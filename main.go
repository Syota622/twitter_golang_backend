package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

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
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
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

	// HTTPサーバー設定
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!")
	})

	fmt.Println("ポート8080でサーバーを起動しています...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
