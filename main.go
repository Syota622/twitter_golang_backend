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

	// SignupHandlerを/signup ルートにマッピング
	router.POST("/signup", api.SignupHandler(queryHandler))

	// GETリクエストに対して "Hello World" を返すルートを追加
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello World")
	})

	// HTTPサーバー起動
	router.Run(":8080") // デフォルトでは localhost:8080 でサーバーを起動
}