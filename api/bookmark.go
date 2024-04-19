package api

import (
	"net/http"
	"strconv"
	"twitter_golang_backend/db/generated" // sqlcで生成されたパッケージをインポート

	"github.com/gin-gonic/gin"
)

// CreateBookmarkHandler は新しいブックマークを作成するためのハンドラ
func CreateBookmarkHandler(db *generated.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			TweetID int32 `json:"tweet_id"`
		}
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// ログインユーザーのIDを取得
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "ユーザーが認証されていません"})
			return
		}

		params := generated.CreateBookmarkParams{
			UserID:  int32(userID.(int)), // ログインユーザーのIDを使用
			TweetID: req.TweetID,
		}

		// ブックマークを作成
		err := db.CreateBookmark(c, params)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ブックマークの作成に成功しました"})
	}
}

// ListBookmarksHandler はブックマーク一覧を取得するためのハンドラ
func ListBookmarksHandler(db *generated.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := strconv.ParseInt(c.Param("userId"), 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ユーザーIDが無効です"})
			return
		}

		// ブックマーク一覧を取得
		bookmarks, err := db.ListBookmarks(c, int32(userID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, bookmarks)
	}
}

// DeleteBookmarkHandler はブックマークを削除するためのハンドラ
func DeleteBookmarkHandler(db *generated.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			TweetID int32 `json:"tweet_id"`
		}
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// ログインユーザーのIDを取得
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "ユーザーが認証されていません"})
			return
		}

		params := generated.DeleteBookmarkParams{
			UserID:  int32(userID.(int)), // ログインユーザーのIDを使用
			TweetID: req.TweetID,
		}

		// ブックマークを削除
		err := db.DeleteBookmark(c, params)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ブックマークの削除に成功しました"})
	}
}
