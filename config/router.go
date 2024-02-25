package config

import (
	"context"
	"database/sql"
	"net/http"
	"twitter_golang_backend/api"
	"twitter_golang_backend/db/generated"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/rs/cors"
)

// SetupRouter Ginルーターを初期化
func SetupRouter(envConfig EnvConfig) *gin.Engine {
	router := gin.Default()
	// CORSの設定
	corsConfig := cors.New(cors.Options{
		AllowedOrigins:   []string{envConfig.FrontendURL},                                         // ReactアプリのURLを許可
		AllowCredentials: true,                                                                    // クッキーを許可
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},                                      // 許可するHTTPメソッド
		AllowedHeaders:   []string{"Content-Type", "Accept", "Authorization", "X-Requested-With"}, // 許可するHTTPヘッダー
	})
	// GolangのrouterにCORSミドルウェアを使用
	// CORSミドルウェアを使用することで、異なるオリジンからのリクエストを許可する
	router.Use(func(c *gin.Context) {
		handler := corsConfig.Handler(
			http.HandlerFunc(
				func(w http.ResponseWriter, r *http.Request) {
					c.Next()
				}))
		handler.ServeHTTP(c.Writer, c.Request) // CORSミドルウェアを使用
	})
	return router
}

// SetupRoutes は、ルートを設定する
func SetupRoutes(router *gin.Engine, db *sql.DB, rdb *redis.Client) {
	queryHandler := generated.New(db)
	ctx := context.Background()

	// ルートの設定
	router.POST("/signup", api.SignupHandler(queryHandler))
	router.GET("/confirm", api.ConfirmEmailHandler(queryHandler))
	router.POST("/login", api.LoginHandler(queryHandler, rdb, ctx))
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello World")
	})
	router.Static("/uploads", "./uploads") // 画像ファイルのアップロード先のディレクトリを指定

	// 認証が必要なAPIグループ
	authGroup := router.Group("/")
	authGroup.Use(AuthMiddleware(rdb)) // 認証のミドルウェアを設定

	// ルートの設定
	authGroup.POST("/tweet", api.CreateTweetWithImageHandler(queryHandler))
	authGroup.GET("/tweets", api.GetAllTweetsHandler(queryHandler))
}
