package api

import (
	"database/sql"
	"net/http"
	"twitter_golang_backend/db/generated"

	"github.com/gin-gonic/gin"
)

// ConfirmEmailHandler は、メールアドレスを確認するハンドラー
func ConfirmEmailHandler(db *generated.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		// クエリパラメータからトークンを取得
		token := c.Query("token")

		if token == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "トークンが必要です"})
			return
		}

		// サインアップ登録したユーザーのメールアドレスを確認して、is_confirmedをtrueにする
		err := db.ConfirmUserEmail(c, sql.NullString{String: token, Valid: true})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "メールの確認に失敗しました"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "メールアドレスが正常に確認されました"})
	}
}
