package api

import (
	"database/sql"
	"net/http"
	"strconv"
	"twitter_golang_backend/db/generated" // sqlcで生成されたパッケージをインポート

	"github.com/gin-gonic/gin"
)

type GetAllTweetsResponse struct {
	generated.GetAllTweetsRow
	RetweetCount int64 `json:"retweet_count"`
}

// CreateTweetWithImageHandler はツイートを投稿するためのハンドラ
func CreateTweetWithImageHandler(db *generated.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		// ユーザーIDの取得
		userID, exists := c.Get("userID")
		println("userID: ", userID)
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

		// レスポンス用のスライスを作成
		response := make([]GetAllTweetsResponse, len(tweets))
		for i, tweet := range tweets {
			response[i] = GetAllTweetsResponse{
				GetAllTweetsRow: tweet,
				RetweetCount:    tweet.RetweetCount,
			}
		}

		// 成功した場合はツイートのリストを含むレスポンスを返す
		c.JSON(http.StatusOK, gin.H{"tweets": response})
	}
}

// GetTweetDetailHandler は特定のツイートIDに対する詳細情報を取得するハンドラ
func GetTweetDetailHandler(db *generated.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		// リクエストからツイートIDを取得
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ツイートIDが無効です"})
			return
		}

		// データベースからツイートの詳細を取得
		tweet, err := db.GetTweetByID(c, int32(id))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ツイートの取得に失敗しました"})
			return
		}

		// 成功した場合はツイートの詳細を含むレスポンスを返す
		c.JSON(http.StatusOK, gin.H{"tweet": tweet})
	}
}

// DeleteTweetHandler はツイートを削除するハンドラ
func DeleteTweetHandler(queryHandler *generated.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		tweetId, err := strconv.ParseInt(c.Param("tweetId"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ツイートIDが無効です"})
			return
		}

		err = queryHandler.DeleteTweet(c, int32(tweetId))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ツイートの削除に失敗しました"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "ツイートの削除に成功しました"})
	}
}
