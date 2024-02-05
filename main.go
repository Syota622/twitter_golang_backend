package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
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

	// データベース接続設定
	psqlInfo := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		envConfig.DBHost,
		envConfig.DBPort,
		envConfig.DBUser,
		envConfig.DBPassword,
		envConfig.DBName,
	)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// データベース接続確認
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("データベースに正常に接続しました！")

	// sqlc用のクエリーハンドラーを生成
	queryHandler := generated.New(db)

	// strconv.Atoiを使用して文字列からintへ変換
	redisDBInt, _ := strconv.Atoi(envConfig.RedisDB)

	// Redisクライアントを初期化
	rdb := redis.NewClient(&redis.Options{
		Addr:     envConfig.RedisAddr,
		Password: envConfig.RedisPassword,
		DB:       redisDBInt,
	})

	// Redisに接続
	ctx := context.Background()
	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Redisに接続できません: %v", err)
	}
	fmt.Println(pong, "Redisに正常に接続しました!")

	// Ginルーターを初期化
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

	// SignupHandlerを/signup ルートにマッピング
	router.POST("/signup", api.SignupHandler(queryHandler))

	// メール確認エンドポイントをルートにマッピング
	router.GET("/confirm", api.ConfirmEmailHandler(queryHandler))

	// LoginHandlerを/login ルートにマッピング
	router.POST("/login", api.LoginHandler(queryHandler, rdb, ctx))

	// GETリクエストに対して "Hello World" を返すルートを追加
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello World")
	})

	// HTTPサーバー起動
	router.Run(":8080") // デフォルトでは localhost:8080 でサーバーを起動
}
