package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"twitter_golang_backend/api"
	"twitter_golang_backend/db/generated" // generatedパッケージをインポート

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
)

const (
	host     = "db"
	port     = 5432
	user     = "postgres"
	password = "Passw0rd"
	dbname   = "db"
)

func main() {
	// データベース接続設定
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
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

	// Ginルーターを初期化
	router := gin.Default()

	// CORSの設定
	corsConfig := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},                                       // ReactアプリのURLを許可
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

	// GETリクエストに対して "Hello World" を返すルートを追加
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello World")
	})

	// LoginHandlerを/login ルートにマッピング
	router.POST("/login", api.LoginHandler(queryHandler))

	// HTTPサーバー起動
	router.Run(":8080") // デフォルトでは localhost:8080 でサーバーを起動
}
