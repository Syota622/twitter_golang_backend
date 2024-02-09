package api

import (
	"net/http"
	"twitter_golang_backend/db/generated" // sqlcで生成されたパッケージをインポート

	"github.com/gin-gonic/gin"
)

// CreateTweetRequest はツイート投稿のためのリクエストボディを定義します。
type CreateTweetRequest struct {
	UserID  int32  `json:"user_id"`
	Message string `json:"message"`
}

// CreateTweetHandler はツイートをデータベースに保存するハンドラです。
func CreateTweetHandler(db *generated.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {

		// ツイート投稿リクエストボディをパース
		var req CreateTweetRequest

		// HTTPリクエストボディをパース
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// ツイートデータをデータベースに挿入する
		arg := generated.CreateTweetParams{
			UserID:  req.UserID,
			Message: req.Message,
		}
		tweet, err := db.CreateTweet(c, arg)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ツイートの投稿に失敗しました"})
			return
		}

		// 成功した場合はツイートデータを含むレスポンスを返す
		c.JSON(http.StatusOK, gin.H{"tweet": tweet})
	}
}
