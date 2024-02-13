package api

import (
	"context"
	"database/sql"
	"net/http"
	"strconv"
	"strings"
	"twitter_golang_backend/db/generated" // sqlcで生成されたパッケージをインポート

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

// CreateTweetRequest はツイート投稿のためのリクエストボディを定義
type CreateTweetRequest struct {
	UserID   int32  `json:"user_id"`
	Message  string `json:"message"`
	ImageURL string `json:"image_url"` // 画像URLを追加
}

// CreateTweetHandler はツイートを投稿するためのハンドラ
func CreateTweetWithImageHandler(db *generated.Queries, rdb *redis.Client, ctx context.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 認証トークンの処理
		token := c.GetHeader("Authorization")
		// トークンがBearerトークンであることを確認
		token = strings.TrimPrefix(token, "Bearer ")

		// RedisからユーザーIDを取得
		userIDStr, err := rdb.Get(ctx, token).Result()
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "認証が必要です"})
			return
		}
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ユーザーIDの取得に失敗しました"})
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
			UserID:   int32(userID), // userIDの取得と変換は省略
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
