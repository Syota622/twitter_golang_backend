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
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},                     // 許可するHTTPメソッド
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
	router.POST("/signup", api.SignupHandler(queryHandler))            // ユーザー登録
	router.GET("/confirm", api.ConfirmEmailHandler(queryHandler))      // メールアドレスの確認
	router.POST("/login", api.LoginHandler(queryHandler, rdb, ctx))    // ログイン
	router.GET("/tweets/:id", api.GetTweetDetailHandler(queryHandler)) // 特定(1つ)のツイートIDに対する詳細情報を取得
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello World")
	})
	router.Static("/uploads", "./uploads") // 画像ファイルのアップロード先のディレクトリを指定

	// 認証が必要なAPIグループ
	authGroup := router.Group("/")
	authGroup.Use(AuthMiddleware(rdb)) // 認証のミドルウェアを設定

	// ルートの設定
	authGroup.POST("/tweet", api.CreateTweetWithImageHandler(queryHandler))               // ツイート登録
	authGroup.GET("/tweets", api.GetAllTweetsHandler(queryHandler))                       // 前ユーザーのツイートリストを取得
	authGroup.DELETE("/tweet/:tweetId", api.DeleteTweetHandler(queryHandler))             // ツイート削除
	authGroup.GET("/user/:userId", api.GetUserProfileHandler(queryHandler))               // ユーザープロフィールを取得
	authGroup.PUT("/user/profile", api.UpdateUserProfileHandler(queryHandler))            // プロフィールを更新
	authGroup.GET("/users/:userId/tweets", api.GetUserTweetsHandler(queryHandler))        // 特定のユーザーのツイートリストを取得
	authGroup.POST("/comment", api.CreateCommentHandler(queryHandler))                    // コメント登録
	authGroup.GET("/comments/:tweetId", api.GetCommentsHandler(queryHandler))             // 特定のツイートIDに対するコメントを取得
	authGroup.DELETE("/comment/:commentId", api.DeleteCommentHandler(queryHandler))       // コメント削除
	authGroup.POST("/retweet/:tweetId", api.CreateRetweetHandler(queryHandler))           // リツイート登録
	authGroup.POST("/like/:tweetId", api.CreateLikeHandler(queryHandler))                 // いいね登録
	authGroup.POST("/follow/:followId", api.CreateFollowHandler(queryHandler))            // フォロー登録
	authGroup.GET("/follow/:followId/is-following", api.IsFollowingHandler(queryHandler)) // フォロー状態を確認
	authGroup.DELETE("/unfollow/:unfollowId", api.UnfollowHandler(queryHandler))          // フォロー解除
}
