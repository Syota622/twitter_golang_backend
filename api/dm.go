package api

import (
	"database/sql"
	"net/http"
	"strconv"
	"twitter_golang_backend/db/generated" // sqlcで生成されたパッケージをインポート

	"github.com/gin-gonic/gin"
)

// CreateGroupHandler は新しいグループを作成するためのハンドラ
func CreateGroupHandler(db *generated.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Name string `json:"name"`
		}
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// グループを作成
		group, err := db.CreateGroup(c, req.Name)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, group)
	}
}

// CreateGroupMessageHandler は新しいグループメッセージを作成するためのハンドラ
func CreateGroupMessageHandler(db *generated.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			GroupID string `json:"group_id"`
			UserID  string `json:"user_id"`
			Message string `json:"message"`
		}
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		groupID, err := strconv.ParseInt(req.GroupID, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group_id"})
			return
		}

		userID, err := strconv.ParseInt(req.UserID, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id"})
			return
		}

		params := generated.CreateGroupMessageParams{
			GroupID: sql.NullInt32{Int32: int32(groupID), Valid: groupID != 0},
			UserID:  int32(userID),
			Message: req.Message,
		}
		message, err := db.CreateGroupMessage(c, params)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, message)
	}
}

// GetAllGroupsHandler は全てのグループを取得するためのハンドラ
func GetAllGroupsHandler(db *generated.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 全てのグループを取得
		groups, err := db.GetAllGroups(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, groups)
	}
}

// GetGroupMessagesHandler は特定のグループのメッセージ一覧を取得するためのハンドラ
func GetGroupMessagesHandler(db *generated.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 特定のグループIDに対するメッセージを取得
		groupID, _ := strconv.ParseInt(c.Param("groupId"), 10, 32)
		messages, err := db.GetGroupMessages(c, sql.NullInt32{Int32: int32(groupID), Valid: groupID != 0})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, messages)
	}
}
