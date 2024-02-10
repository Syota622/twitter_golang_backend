package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"path/filepath"
	"twitter_golang_backend/db/generated" // sqlcで生成されたパッケージをインポート

	"github.com/gin-gonic/gin"
)

// CreateTweetRequest はツイート投稿のためのリクエストボディを定義します。
type CreateTweetRequest struct {
	UserID   int32  `json:"user_id"`
	Message  string `json:"message"`
	ImageURL string `json:"image_url"` // 画像URLを追加
}

// UploadImageHandler は画像をアップロードするためのハンドラ
func UploadImageHandler(c *gin.Context) {
	file, err := c.FormFile("image") // file: アップロードされたファイル
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 画像ファイルを保存するパスを生成します。
	filename := filepath.Base(file.Filename)       // filename: "image.png"
	savePath := filepath.Join("uploads", filename) // savePath: "uploads/image.png"

	// ファイルをサーバーに保存します。
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 保存したファイルへのURLをレスポンスとして返します。
	// 実際のURLは、アプリケーションの設定に応じて変更してください。
	fileURL := fmt.Sprintf("http://localhost:8080/%s", savePath)
	c.JSON(http.StatusOK, gin.H{"url": fileURL})
}

// CreateTweetWithImageHandler はツイートをデータベースに保存するハンドラ
// TweetHandler はツイートと画像URLをデータベースに保存するハンドラ
func CreateTweetWithImageHandler(db *generated.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {

		// ツイート投稿リクエストボディをパース
		var req CreateTweetRequest

		// HTTPリクエストボディをパース
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// ツイートをデータベースに保存
		arg := generated.CreateTweetParams{
			UserID:   req.UserID,
			Message:  req.Message,
			ImageUrl: sql.NullString{String: req.ImageURL, Valid: req.ImageURL != ""},
		}
		tweet, err := db.CreateTweet(c, arg)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ツイートの投稿に失敗しました"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"tweet": tweet})
	}
}

// GetTweetsHandler はデータベースからツイートのリストを取得するハンドラ
func GetTweetsHandler(db *generated.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		tweets, err := db.GetTweets(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ツイートの取得に失敗しました"})
			return
		}

		// 成功した場合はツイートのリストを含むレスポンスを返す
		c.JSON(http.StatusOK, gin.H{"tweets": tweets})
	}
}
