package api

import (
	"net/http"
	"strconv"
	"twitter_golang_backend/db/generated"

	"github.com/gin-gonic/gin"
)

// CreateFollowHandler はフォローを作成するためのハンドラ
func CreateFollowHandler(db *generated.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 認証されたユーザーのIDを取得: フォローするユーザーのID
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "認証が必要です"})
			return
		}

		// パラメーターから取得: フォローされるユーザーのIDを取得
		followID, err := strconv.Atoi(c.Param("followId"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "無効なユーザーID"})
			return
		}

		// データベースにフォローを保存
		arg := generated.CreateFollowParams{
			UserID:   int32(userID.(int)), // ユーザーIDの取得
			FollowID: int32(followID),     // フォローIDの取得
		}
		follow, err := db.CreateFollow(c, arg)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "フォローに失敗しました"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"follow": follow})
	}
}

// IsFollowingHandler はログインユーザーが指定したユーザーをフォローしているかどうかを確認するためのハンドラ
func IsFollowingHandler(db *generated.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 認証されたユーザーのIDを取得
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "認証が必要です"})
			return
		}

		// パラメーターから取得: フォロー状態を確認するユーザーのIDを取得
		followIdParam, err := strconv.Atoi(c.Param("followId"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "無効なfollowIdパラメータ"})
			return
		}

		// フォローの状態を取得
		arg := generated.IsFollowingParams{
			UserID:   int32(userID.(int)),  // ユーザーIDの取得
			FollowID: int32(followIdParam), // フォローIDの取得
		}
		isFollowing, err := db.IsFollowing(c, arg)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "フォローの状態の取得に失敗しました"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"isFollowing": isFollowing})
	}
}

// UnfollowHandler はログインユーザーが指定したユーザーのフォローを解除するためのハンドラ
func UnfollowHandler(db *generated.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 認証されたユーザーのIDを取得
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "認証が必要です"})
			return
		}

		// パラメーターから取得: フォローを解除するユーザーのIDを取得
		unfollowIdParam, err := strconv.Atoi(c.Param("unfollowId"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "無効なunfollowIdパラメータ"})
			return
		}

		// データベースからフォローを解除
		arg := generated.UnfollowParams{
			UserID:   int32(userID.(int)),    // ユーザーIDの取得
			FollowID: int32(unfollowIdParam), // フォローを解除するユーザーのIDの取得
		}
		_, err = db.Unfollow(c, arg)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "フォローの解除に失敗しました"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "フォローを解除しました"})
	}
}
