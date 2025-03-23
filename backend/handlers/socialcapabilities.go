package handlers

import (
	"heydays/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) GetFriendList(c *gin.Context) {
	userId, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "cannot get user_id"})
		return
	}

	var friends []models.User
	if err := s.db.Joins(
		"JOIN friendships ON (friendships.user_id = users.id AND friendships.friend_id = ?) OR (friendships.friend_id = users.id AND friendships.user_id = ?)",
		userId, userId).
		Where("friendships.status = ?", "accepted").
		Find(&friends).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get friends"})
		return
	}

	friendList := make([]gin.H, len(friends))
	for i, friend := range friends {
		friendList[i] = gin.H{
			"id":            friend.ID,
			"name":          friend.Name,
			"surname":       friend.Surname,
			"profile_photo": friend.ProfilePhoto,
		}
	}

	c.JSON(http.StatusOK, gin.H{"friends": friendList})
}

func (s *Server) GetPendingFriendInvies(c *gin.Context) {
	userId, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "cannot get user_id"})
		return
	}

	var friends []models.User
	if err := s.db.Joins(
		"JOIN friendships ON (friendships.user_id = users.id AND friendships.friend_id = ?) OR (friendships.friend_id = users.id AND friendships.user_id = ?)",
		userId, userId).
		Where("friendships.status = ?", "pending").
		Find(&friends).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no friends for you"})
		return
	}

	friendList := make([]gin.H, len(friends))
	for i, friend := range friends {
		friendList[i] = gin.H{
			"id":            friend.ID,
			"name":          friend.Name,
			"surname":       friend.Surname,
			"profile_photo": friend.ProfilePhoto,
		}
	}

	c.JSON(http.StatusOK, gin.H{"friends": friendList})
}

func (s *Server) SendFriendRequest(c *gin.Context) {
	userId, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "cannot get user_id"})
		c.Abort()
		return
	}

	type request struct {
		FriendId uint `json:"friend_id" binding:"required"`
	}

	var req request
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newFriendship := models.Friendship{UserID: userId.(uint), FriendID: req.FriendId, Status: "pending"}
	if err := s.db.Create(&newFriendship).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot send friend request"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "succes"})

}

func (s *Server) AcceptFriendRequest(c *gin.Context) {
	userId, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "cannot get user_id"})
		c.Abort()
		return
	}

	type request struct {
		FriendId uint `json:"friend_id" binding:"required"`
	}

	var req request
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := s.db.Model(&models.Friendship{}).
		Where("user_id = ? AND friend_id = ? AND status = ?", req.FriendId, userId, "pending").
		Update("status", "accepted").Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "friend request not found or cannot be accepted"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "succes"})
}

func (s *Server) DeclineFriendRequest(c *gin.Context) {

	userId, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "cannot get user_id"})
		c.Abort()
		return
	}

	type request struct {
		FriendId uint `json:"friend_id" binding:"required"`
	}

	var req request
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var friendship models.Friendship
	if err := s.db.Where("user_id = ? AND friend_id = ? AND status = ?", req.FriendId, userId, "pending").First(&friendship).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "friend request not found"})
		return
	}

	if err := s.db.Delete(&friendship).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to decline friend request"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success"})
}
