package api

import (
	"net/http"
	"twitter_golang_backend/db/generated"
	"twitter_golang_backend/util"

	"github.com/gin-gonic/gin"
)

// サインアップリクエストボディ
type CreateUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

// SignupHandler は、ユーザーを作成するハンドラー
func SignupHandler(db *generated.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {

		// サインアップリクエストボディをパース
		var req CreateUserRequest

		// HTTPリクエストボディをパース
		// リクエストボディがJSON形式でない場合はエラーを返す
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// パスワードをハッシュ化
		hashedPassword, err := util.HashPassword(req.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
			return
		}

		// ユーザーを作成
		user, err := db.CreateUser(c, generated.CreateUserParams{
			Username:       req.Username,
			HashedPassword: hashedPassword,
			Email:          req.Email,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, user)
	}
}
