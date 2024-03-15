package api

import (
	"net/http"
	"strconv"
	"twitter_golang_backend/db/generated" // このパスはプロジェクトによって異なります

	"github.com/gin-gonic/gin"
)

// 特定のユーザーのツイート一覧を取得するハンドラー
func GetUserTweetsHandler(queries *generated.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		// URLパラメータからユーザーIDを取得
		userIdParam := c.Param("userId")
		// 文字列のユーザーIDをint32に変換
		userIDInt, err := strconv.ParseInt(userIdParam, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "無効なユーザーID"})
			return
		}

		tweets, err := queries.GetUserTweets(c, int32(userIDInt))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ツイートの取得に失敗しました"})
			return
		}

		// tweetsキーにtweetsの配列を入れる
		c.JSON(http.StatusOK, gin.H{"tweets": tweets})
	}
}
