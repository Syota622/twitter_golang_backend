package api

import (
	"net/http"
	"twitter_golang_backend/db/generated"

	"github.com/gin-gonic/gin"
)

// DeleteUserHandler はユーザーを退会させるためのハンドラ
func DeleteUserHandler(db *generated.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		// ログインユーザーのIDを取得
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "ユーザーが認証されていません"})
			return
		}

		// ユーザーを退会させる
		err := db.DeactivateUser(c, int32(userID.(int)))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ユーザーの退会に成功しました"})
	}
}
