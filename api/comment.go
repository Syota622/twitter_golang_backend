package api

import (
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
		comment := c.PostForm("comment")
		if comment == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "コメントが空です"})
			return
		}

		// データベースに新しいコメントを作成
		newComment, err := db.CreateComment(c, generated.CreateCommentParams{
			TweetID: int32(tweetID),
			Comment: comment,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "コメントの作成に失敗しました"})
			return
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
