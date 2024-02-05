package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"twitter_golang_backend/db/generated"
	"twitter_golang_backend/utils"

	"github.com/gin-gonic/gin"

	"crypto/rand"
	"encoding/hex"
)

// サインアップリクエストボディ
type CreateUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

// トークン生成
func generateConfirmationToken() (string, error) {
	tokenBytes := make([]byte, 16)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(tokenBytes), nil
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
		hashedPassword, err := utils.HashPassword(req.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
			return
		}

		// トークン生成
		confirmationToken, err := generateConfirmationToken()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "トークン生成に失敗しました"})
			return
		}

		// ユーザーを作成
		user, err := db.CreateUser(c, generated.CreateUserParams{
			Username:          req.Username,
			HashedPassword:    hashedPassword,
			Email:             req.Email,
			ConfirmationToken: sql.NullString{String: confirmationToken, Valid: true},
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// メールアドレス確認用のリンクを作成(router.GET("/confirm", api.ConfirmEmailHandler(queryHandler)))
		confirmationLink := fmt.Sprintf("http://localhost:8080/confirm?token=%s", confirmationToken)
		subject := "メールアドレスの確認をお願いします"
		body := fmt.Sprintf("以下のリンクをクリックしてメールアドレスの確認を完了してください: %s", confirmationLink)

		// メール送信
		err = utils.SendEmail(req.Email, subject, body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "確認メールの送信に失敗しました"})
			return
		}

		c.JSON(http.StatusOK, user)
	}
}
