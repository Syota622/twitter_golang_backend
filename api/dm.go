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
func CreateMessageHandler(db *generated.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Message string `json:"message"`
		}
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// URLパラメータからgroup_idを取得
		groupID, err := strconv.ParseInt(c.Param("groupId"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "グループIDが無効です"})
			return
		}

		// コンテキストからuserIDを取得
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ユーザーIDが無効です"})
			return
		}

		// グループメッセージを作成
		params := generated.CreateMessageParams{
			GroupID: sql.NullInt32{Int32: int32(groupID), Valid: groupID != 0},
			UserID:  int32(userID.(int)), // userIDをint型にキャスト
			Message: req.Message,
		}
		message, err := db.CreateMessage(c, params)
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

// GetGroupMessagesHandler は特定のグループIDに対するメッセージを取得するためのハンドラ
func GetMessagesHandler(db *generated.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		// URLパラメータからgroup_idを取得
		groupID, err := strconv.ParseInt(c.Param("groupId"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "グループIDが無効です"})
			return
		}

		// 特定のグループIDに対するメッセージを取得
		params := sql.NullInt32{Int32: int32(groupID), Valid: groupID != 0}
		messages, err := db.GetMessages(c, params)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, messages)
	}
}
