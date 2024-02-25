package api

import (
	"database/sql"
	"net/http"
	"strconv"
	"twitter_golang_backend/db/generated" // sqlcで生成されたパッケージをインポート

	"github.com/gin-gonic/gin"
)

// CreateTweetRequest はツイート投稿のためのリクエストボディを定義
type CreateTweetRequest struct {
	UserID   int32  `json:"user_id"`
	Message  string `json:"message"`
	ImageURL string `json:"image_url"` // 画像URLを追加
}

// CreateTweetWithImageHandler はツイートを投稿するためのハンドラ
func CreateTweetWithImageHandler(db *generated.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		// ユーザーIDの取得
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "認証が必要です"})
			return
		}

		// フォームデータの処理
		message := c.PostForm("message")
		file, err := c.FormFile("image")
		var imageUrl sql.NullString

		if err != nil && err != http.ErrMissingFile {
			c.JSON(http.StatusBadRequest, gin.H{"error": "画像ファイルの取得に失敗しました"})
			return
		}

		// ファイルが存在する場合、サーバーに保存
		var filePath string
		if file != nil {
			filePath = "./uploads/" + file.Filename // 保存先のパスを指定
			// ファイルを保存
			if err := c.SaveUploadedFile(file, filePath); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "ファイルの保存に失敗しました"})
				return
			}
			// クライアントからアクセス可能なURLを組み立て
			imageUrl = sql.NullString{String: "http://localhost:8080/uploads/" + file.Filename, Valid: true}
		}

		// データベースにツイートを保存
		arg := generated.CreateTweetParams{
			UserID:   int32(userID.(int)), // ユーザーIDの取得
			Message:  message,
			ImageUrl: imageUrl,
		}
		tweet, err := db.CreateTweet(c, arg)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ツイートの投稿に失敗しました"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"tweet": tweet})
	}
}

// GetAllTweetsHandler はデータベースからツイートのリストを取得するハンドラ
func GetAllTweetsHandler(db *generated.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		// リクエストから limit と offset の値を取得
		limitParam := c.DefaultQuery("limit", "10")
		offsetParam := c.DefaultQuery("offset", "0")

		limit, err := strconv.Atoi(limitParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit value"})
			return
		}

		offset, err := strconv.Atoi(offsetParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset value"})
			return
		}

		// データベースからツイートを取得する
		params := generated.GetAllTweetsParams{
			Limit:  int32(limit),
			Offset: int32(offset),
		}
		tweets, err := db.GetAllTweets(c, params)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ツイートの取得に失敗しました"})
			return
		}

		// 成功した場合はツイートのリストを含むレスポンスを返す
		c.JSON(http.StatusOK, gin.H{"tweets": tweets})
	}
}
