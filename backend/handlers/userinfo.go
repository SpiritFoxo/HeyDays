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

//update profile

func (s *Server) UpdateProfileInfo(c *gin.Context) {
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

	type UpdateProfileInput struct {
		Name         string `json:"name"`
		Surname      string `json:"surname"`
		ProfilePhoto string `json:"profile_photo"`
	}

	var input UpdateProfileInput

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Name != "" {
		user.Name = input.Name
	}

	if input.Surname != "" {
		user.Surname = input.Surname
	}

	if input.ProfilePhoto != "" {
		user.ProfilePhoto = input.ProfilePhoto
	}

	if err := s.db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile updated"})
}
