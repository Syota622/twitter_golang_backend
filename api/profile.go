package api

import (
	"database/sql"
	"net/http"
	"twitter_golang_backend/db/generated" // このパスはプロジェクトによって異なる

	"github.com/gin-gonic/gin"
)

// UpdateUserProfileRequest はリクエストボディの構造体です
type UpdateUserProfileRequest struct {
	Username           string `json:"username"`
	Email              string `json:"email"`
	Bio                string `json:"bio,omitempty"`
	ProfileImageURL    string `json:"profile_image_url,omitempty"`
	BackgroundImageURL string `json:"background_image_url,omitempty"`
}

// UpdateUserProfileHandler はプロフィール更新のハンドラーです
func UpdateUserProfileHandler(queries *generated.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "認証が必要です"})
			return
		}

		var req UpdateUserProfileRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// userIDの型アサーション（ここではuserIDがint型であることを前提としています）
		userIDInt, ok := userID.(int)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "サーバー内部エラー"})
			return
		}

		// リクエストから受け取った string 型のフィールドを sql.NullString 型に変換
		bio := sql.NullString{String: req.Bio, Valid: req.Bio != ""}
		profileImageUrl := sql.NullString{String: req.ProfileImageURL, Valid: req.ProfileImageURL != ""}
		backgroundImageUrl := sql.NullString{String: req.BackgroundImageURL, Valid: req.BackgroundImageURL != ""}

		err := queries.UpdateUserProfile(c, generated.UpdateUserProfileParams{
			ID:                 int32(userIDInt), // userIDをint32型に変換
			Username:           req.Username,
			Email:              req.Email,
			Bio:                bio,
			ProfileImageUrl:    profileImageUrl,
			BackgroundImageUrl: backgroundImageUrl,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "プロフィールの更新に失敗しました"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "プロフィールが正常に更新されました"})
	}
}
