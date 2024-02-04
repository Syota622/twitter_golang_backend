package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// EnvConfig は環境変数を保持する構造体
type EnvConfig struct {
	FrontendURL string
	DBHost      string
	DBPort      string
	DBUser      string
	DBPassword  string
	DBName      string
}

// GetEnvConfig は環境変数を取得する関数
func GetEnvConfig() *EnvConfig {
	// rootディレクトリの.envファイルを読み込む
	err := godotenv.Load()
	if err != nil {
		log.Fatal(".envファイルは見つかりませんでした")
	}

	// 環境変数を取得
	return &EnvConfig{
		FrontendURL: os.Getenv("FRONTEND_URL"),
		DBHost:      os.Getenv("DB_HOST"),
		DBPort:      os.Getenv("DB_PORT"),
		DBUser:      os.Getenv("DB_USER"),
		DBPassword:  os.Getenv("DB_PASSWORD"),
		DBName:      os.Getenv("DB_NAME"),
	}
}
