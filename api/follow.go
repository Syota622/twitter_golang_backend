package api

import (
	"net/http"
	"strconv"
	"twitter_golang_backend/db/generated"

	"github.com/gin-gonic/gin"
)

// CreateFollowHandler はフォローを作成するためのハンドラ
func CreateFollowHandler(db *generated.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 認証されたユーザーのIDを取得: フォローするユーザーのID
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "認証が必要です"})
			return
		}

		// パラメーターから取得: フォローされるユーザーのIDを取得
		followID, err := strconv.Atoi(c.Param("followId"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "無効なユーザーID"})
			return
		}

		// データベースにフォローを保存
		arg := generated.CreateFollowParams{
			UserID:   int32(userID.(int)), // ユーザーIDの取得
			FollowID: int32(followID),     // フォローIDの取得
		}
		follow, err := db.CreateFollow(c, arg)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "フォローに失敗しました"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"follow": follow})
	}
}
