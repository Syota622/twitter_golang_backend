package api

import (
	"net/http"
	"strconv"
	"twitter_golang_backend/db/generated" // sqlcで生成されたパッケージをインポート

	"github.com/gin-gonic/gin"
)

// CreateLikeHandler は「いいね」を作成するためのハンドラ
func CreateLikeHandler(db *generated.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 認証されたユーザーのIDを取得
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "認証が必要です"})
			return
		}

		// いいねするツイートIDの取得（URLパラメータから取得）
		tweetID, err := strconv.Atoi(c.Param("tweetId"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "無効なツイートID"})
			return
		}

		// データベースにいいねを保存
		arg := generated.CreateLikeParams{
			UserID:  int32(userID.(int)), // ユーザーIDの取得
			TweetID: int32(tweetID),      // ツイートIDの取得
		}
		like, err := db.CreateLike(c, arg)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "いいねに失敗しました"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"like": like})
	}
}
