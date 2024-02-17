package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"twitter_golang_backend/api"
	"twitter_golang_backend/config"
	"twitter_golang_backend/db/generated" // generatedパッケージをインポート

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
)

func main() {

	// 環境変数の取得
	envConfig := config.GetEnvConfig()

	// データベース接続
	db := connectDatabase(*envConfig)
	defer db.Close()

	// Redis接続
	rdb := connectRedis(*envConfig)

	// Ginルーターの初期化
	router := setupRouter(*envConfig)

	// ルートの設定
	setupRoutes(router, db, rdb)

	// HTTPサーバー起動
	router.Run(":8080") // デフォルトでは localhost:8080 でサーバーを起動
}

// connectDatabase は、データベースに接続する
func connectDatabase(envConfig config.EnvConfig) *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		envConfig.DBHost, envConfig.DBPort, envConfig.DBUser, envConfig.DBPassword, envConfig.DBName)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("データベース接続エラー: %v", err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatalf("データベース接続エラー: %v", err)
	}
	fmt.Println("データベースに接続しました！")
	return db
}

// connectRedis は、Redisに接続する
func connectRedis(envConfig config.EnvConfig) *redis.Client {
	redisDBInt, _ := strconv.Atoi(envConfig.RedisDB)
	rdb := redis.NewClient(&redis.Options{
		Addr:     envConfig.RedisAddr,
		Password: envConfig.RedisPassword,
		DB:       redisDBInt,
	})
	ctx := context.Background()
	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Redis接続エラー: %v", err)
	}
	fmt.Println(pong, "Redisに接続しました！")
	return rdb
}

// setupRouter は、Ginルーターを初期化する
func setupRouter(envConfig config.EnvConfig) *gin.Engine {
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

// setupRoutes は、ルートを設定する
func setupRoutes(router *gin.Engine, db *sql.DB, rdb *redis.Client) {
	queryHandler := generated.New(db)
	ctx := context.Background()

	// ルートの設定
	router.POST("/signup", api.SignupHandler(queryHandler))
	router.GET("/confirm", api.ConfirmEmailHandler(queryHandler))
	router.POST("/login", api.LoginHandler(queryHandler, rdb, ctx))
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello World")
	})

	// 認証が必要なAPIグループ
	authGroup := router.Group("/")
	authGroup.Use(authMiddleware(rdb, ctx)) // 認証のミドルウェアを設定

	// ルートの設定
	authGroup.Static("/uploads", "./uploads") // 画像ファイルのアップロード先のディレクトリを指定
	authGroup.POST("/tweet", api.CreateTweetWithImageHandler(queryHandler, ctx))
	authGroup.GET("/tweets", api.GetAllTweetsHandler(queryHandler))
}

// authMiddleware は認証のミドルウェアです
func authMiddleware(rdb *redis.Client, ctx context.Context) gin.HandlerFunc {
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
		userIDStr, err := rdb.Get(ctx, token).Result()
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
