package config

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

// AuthMiddleware 認証のミドルウェア
func AuthMiddleware(rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 認証トークンの取得("Bearer xxxxxxxxx...")
		token := c.GetHeader("Authorization")

		// トークンの有無を確認
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "認証が必要です"})
			c.Abort()
			return
		}

		// トークンがBearerトークンであることを確認
		token = strings.TrimPrefix(token, "Bearer ")

		// RedisからユーザーIDを取得
		userIDStr, err := rdb.Get(context.Background(), token).Result()
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "認証が必要です"})
			c.Abort()
			return
		}

		// ユーザーIDの変換
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ユーザーIDの取得に失敗しました"})
			c.Abort()
			return
		}

		// 後続のハンドラでユーザーIDを使用できるように設定
		c.Set("userID", userID)

		// ミドルウェアチェーンを継続
		c.Next()
	}
}
