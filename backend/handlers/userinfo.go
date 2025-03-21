package handlers

import (
	"heydays/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) GetProfileInfo(c *gin.Context) {
	userId, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "cannot get user_id"})
		c.Abort()
		return
	}

	user := models.User{}
	result := s.db.First(&user, userId)
	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "cannot find user"})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"name":          user.Name,
		"surname":       user.Surname,
		"profile_photo": user.ProfilePhoto,
	})
}

func (s *Server) GetStrangerProfileInfo(c *gin.Context) {
	userId := c.Param("userId")

	user := models.User{}
	result := s.db.First(&user, userId)
	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "cannot find user"})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"name":          user.Name,
		"surname":       user.Surname,
		"profile_photo": user.ProfilePhoto,
	})
}
