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
			UserID  int32 `json:"user_id"`
			TweetID int32 `json:"tweet_id"`
		}
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		params := generated.CreateBookmarkParams{
			UserID:  req.UserID,
			TweetID: req.TweetID,
		}

		// ブックマークを作成
		err := db.CreateBookmark(c, params)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "success"})
	}
}

// ListBookmarksHandler はブックマーク一覧を取得するためのハンドラ
func ListBookmarksHandler(db *generated.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := strconv.ParseInt(c.Param("userId"), 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
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
