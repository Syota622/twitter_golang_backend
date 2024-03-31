package api

import (
	"net/http"
	"strconv"
	"twitter_golang_backend/db/generated"

	"github.com/gin-gonic/gin"
)

// CreateRetweetHandler はリツイートを作成するためのハンドラ
func CreateRetweetHandler(db *generated.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 認証されたユーザーのIDを取得
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "認証が必要です"})
			return
		}

		// リツイートするツイートIDの取得
		tweetID, err := strconv.Atoi(c.Param("tweetId"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "無効なツイートID"})
			return
		}

		// データベースにリツイートを保存
		arg := generated.CreateRetweetParams{
			UserID:  int32(userID.(int)), // ユーザーIDの取得
			TweetID: int32(tweetID),      // ツイートIDの取得
		}
		retweet, err := db.CreateRetweet(c, arg)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "リツイートの作成に失敗しました"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"retweet": retweet})
	}
}
