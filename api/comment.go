package api

import (
	"database/sql"
	"net/http"
	"strconv"
	"twitter_golang_backend/db/generated" // sqlcで生成されたパッケージをインポート

	"github.com/gin-gonic/gin"
)

// CreateCommentHandler は新しいコメントを作成するハンドラ
func CreateCommentHandler(db *generated.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		// リクエストから必要なパラメータを取得
		tweetID, err := strconv.Atoi(c.PostForm("tweet_id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ツイートIDが無効です"})
			return
		}

		// ツイートの詳細を取得
		tweet, err := db.GetTweetByID(c, int32(tweetID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ツイートの取得に失敗しました"})
			return
		}

		// リクエストからコメントを取得
		comment := c.PostForm("comment")
		if comment == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "コメントが空です"})
			return
		}

		// リクエストからユーザーIDを取得
		userID := c.MustGet("userID").(int)

		// データベースに新しいコメントを作成
		newComment, err := db.CreateComment(c, generated.CreateCommentParams{
			TweetID: tweet.ID,
			Comment: comment,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "コメントの作成に失敗しました"})
			return
		}

		// 新しい通知を作成
		notificationParams := generated.CreateNotificationParams{
			UserID:       tweet.UserID,                                          // 通知を受け取るユーザーID
			NotifiedByID: int32(userID),                                         // 通知を送るユーザーID
			Type:         "comment",                                             // 通知のタイプ
			PostID:       sql.NullInt32{Int32: newComment.TweetID, Valid: true}, // コメントされたツイートID
			CommentID:    sql.NullInt32{Int32: newComment.ID, Valid: true},      // 新しいコメントID
		}
		_, err = db.CreateNotification(c, notificationParams)
		if err != nil {
			// 通知の作成に失敗した場合でも、コメントの作成は続行されます
			c.JSON(http.StatusInternalServerError, gin.H{"error": "通知の作成に失敗しました"})
		}

		// 成功した場合は新しいコメントの詳細を含むレスポンスを返す
		c.JSON(http.StatusOK, gin.H{"comment": newComment})
	}
}

// GetCommentsHandler は指定されたツイートIDに関連するすべてのコメントを取得するハンドラ
func GetCommentsHandler(db *generated.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		// リクエストから必要なパラメータを取得
		tweetID, err := strconv.Atoi(c.Param("tweetId"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ツイートIDが無効です"})
			return
		}

		// データベースからコメントを取得
		comments, err := db.GetCommentsByTweetID(c, int32(tweetID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "コメントの取得に失敗しました"})
			return
		}

		// 成功した場合はコメントのリストを含むレスポンスを返す
		c.JSON(http.StatusOK, gin.H{"comments": comments})
	}
}

// DeleteCommentHandler は指定されたコメントIDを持つコメントを削除するハンドラ
func DeleteCommentHandler(db *generated.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		// リクエストから必要なパラメータを取得
		commentID, err := strconv.Atoi(c.Param("commentId"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "コメントIDが無効です"})
			return
		}

		// データベースからコメントを削除
		err = db.DeleteComment(c, int32(commentID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "コメントの削除に失敗しました"})
			return
		}

		// 成功した場合はステータスコード200を返す
		c.JSON(http.StatusOK, gin.H{"message": "コメントが正常に削除されました"})
	}
}
