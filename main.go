package main

import (
	"twitter_golang_backend/config"

	_ "github.com/lib/pq"
)

func main() {

	// 環境変数の取得
	envConfig := config.GetEnvConfig()

	// データベース接続
	db := config.ConnectDatabase(*envConfig)
	defer db.Close()

	// Redis接続
	rdb := config.ConnectRedis(*envConfig)

	// Ginルーターの初期化
	router := config.SetupRouter(*envConfig)

	// ルートの設定
	config.SetupRoutes(router, db, rdb)

	// HTTPサーバー起動
	router.Run(":8080") // デフォルトでは localhost:8080 でサーバーを起動
}
