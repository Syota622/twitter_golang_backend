package config

import (
	"database/sql"
	"fmt"
	"log"
)

// ConnectDatabase データベースに接続
func ConnectDatabase(envConfig EnvConfig) *sql.DB {
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
