package api

import (
	"database/sql"
	"log"
	"net/http"
	"twitter_golang_backend/db/generated"

	"github.com/gin-gonic/gin"
)

type CreateNotificationRequest struct {
	UserID       int32  `json:"user_id"`
	NotifiedByID int32  `json:"notified_by_id"`
	Type         string `json:"type"`
	PostID       *int32 `json:"post_id"`
	CommentID    *int32 `json:"comment_id"`
}

// CreateNotificationHandler は新しい通知を作成するためのハンドラ
func CreateNotificationHandler(db *generated.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {

		// リクエストデータのバインド
		var req CreateNotificationRequest
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "無効なリクエストデータ"})
			return
		}

		// 通知の作成
		params := generated.CreateNotificationParams{
			UserID:       req.UserID,
			NotifiedByID: req.NotifiedByID,
			Type:         req.Type,
		}

		// リクエストのPostIDとCommentIDがnilでない場合、それぞれの値を設定
		if req.PostID != nil {
			params.PostID = sql.NullInt32{Int32: *req.PostID, Valid: true}
		}
		if req.CommentID != nil {
			params.CommentID = sql.NullInt32{Int32: *req.CommentID, Valid: true}
		}

		notification, err := db.CreateNotification(c.Request.Context(), params)
		if err != nil {
			// エラーの詳細をログに出力
			log.Printf("failed to create notification: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "通知の作成に失敗しました"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"notification": notification})
	}
}

// GetNotificationsHandler はユーザーの通知のリストを取得するハンドラ
func GetNotificationsHandler(db *generated.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		// リクエストからユーザーIDを取得
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "認証が必要です"})
			return
		}

		// データベースから通知を取得
		notifications, err := db.GetNotificationsByUserID(c, int32(userID.(int)))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "通知の取得に失敗しました"})
			return
		}

		// 成功した場合は通知のリストを含むレスポンスを返す
		c.JSON(http.StatusOK, gin.H{"notifications": notifications})
	}
}
