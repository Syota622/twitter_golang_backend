package api

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"twitter_golang_backend/db/generated" // sqlcで生成されたパッケージをインポート

	"github.com/gin-gonic/gin"
)

// CreateLikeHandler は「いいね」を作成するためのハンドラ
func CreateLikeHandler(db *generated.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 認証されたユーザーのIDを取得
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "認証が必要です"})
			return
		}

		// いいねするツイートIDの取得（URLパラメータから取得）
		tweetID, err := strconv.Atoi(c.Param("tweetId"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "無効なツイートID"})
			return
		}

		// データベースにいいねを保存
		arg := generated.CreateLikeParams{
			UserID:  int32(userID.(int)), // ユーザーIDの取得
			TweetID: int32(tweetID),      // ツイートIDの取得
		}
		like, err := db.CreateLike(c, arg)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "いいねに失敗しました"})
			return
		}

		// ツイートを取得
		tweet, err := db.GetTweetByID(c, int32(tweetID))
		log.Println("tweet: ", tweet)
		if err != nil {
			// ツイートが見つからない場合はエラーを返す
			c.JSON(http.StatusNotFound, gin.H{"error": "ツイートが見つかりません"})
			return
		}

		// データベースに通知を保存
		notificationParams := generated.CreateNotificationParams{
			UserID:       tweet.UserID,                                      // 通知を受け取るユーザーID
			NotifiedByID: int32(userID.(int)),                               // 通知を送るユーザーID
			Type:         "like",                                            // 通知のタイプ
			PostID:       sql.NullInt32{Int32: int32(tweetID), Valid: true}, // いいねされたツイートID
		}
		notification, err := db.CreateNotification(c, notificationParams)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "通知の作成に失敗しました"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"like": like, "notification": notification})
	}
}
