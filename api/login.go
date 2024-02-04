package api

import (
	"net/http"
	"twitter_golang_backend/db/generated"
	"twitter_golang_backend/utils"

	"github.com/gin-gonic/gin"
)

// ログインリクエストボディ
type LoginUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginHandler は、ユーザーを認証するハンドラー
func LoginHandler(db *generated.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {

		// ログインリクエストボディをパース
		var req LoginUserRequest

		// HTTPリクエストボディをパース
		// リクエストボディがJSON形式でない場合はエラーを返す
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// ユーザーを検索
		user, err := db.GetUserByEmail(c, req.Username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ユーザーが見つかりません"})
			return
		}

		// パスワードを検証
		match := utils.CheckPasswordHash(req.Password, user.HashedPassword)
		if !match {
			c.JSON(http.StatusBadRequest, gin.H{"error": "パスワードが間違っています"})
			return
		}

		// ユーザー認証が成功
		c.JSON(http.StatusOK, user)
	}
}
