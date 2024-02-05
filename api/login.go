package api

import (
	"context"
	"net/http"
	"twitter_golang_backend/db/generated"
	"twitter_golang_backend/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

// ログインリクエストボディ
type LoginUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginHandler は、ユーザーを認証するハンドラー
func LoginHandler(db *generated.Queries, rdb *redis.Client, ctx context.Context) gin.HandlerFunc {
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
		user, err := db.GetUserByEmail(c, req.Email)
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

		// ユーザー認証が成功したら、リクエストデータを含むレスポンスを返す
		c.JSON(http.StatusOK, gin.H{
			"message":      "ログインに成功しました",
			"request_data": req,
		})

		// ユーザー認証が成功した場合
		// セッションID（トークン）を生成
		sessionToken := uuid.NewString()

		// トークンとユーザーIDをRedisに保存
		_, err = rdb.Set(ctx, sessionToken, user.ID, 0).Result()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "セッションの保存に失敗しました"})
			return
		}

		// クライアントにセッションID（トークン）を返す
		c.JSON(http.StatusOK, gin.H{"token": sessionToken})
	}
}
