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
		c.Abort()
		return
	}

	var friends []models.Friendship
	if err := s.db.Where("(user_id = ? OR friend_id = ?) AND status = 'accepted'", userId, userId).Find(&friends).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no friends for you"})
		return
	}

	friendsIds := make([]uint, 0)
	for _, friend := range friends {
		if friend.UserID == userId {
			friendsIds = append(friendsIds, friend.FriendID)
		} else {
			friendsIds = append(friendsIds, friend.UserID)
		}
	}

	c.JSON(http.StatusOK, gin.H{"friends": friendsIds})
}

func (s *Server) GetPendingFriendInvies(c *gin.Context) {
	userId, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "cannot get user_id"})
		c.Abort()
		return
	}

	var friends []models.Friendship
	if err := s.db.Where("(user_id = ? OR friend_id = ?) AND status = ?", userId, userId, "pending").Find(&friends).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no friends for you"})
		return
	}

	friendsIds := make([]uint, 0)
	for _, friend := range friends {
		if friend.UserID == userId {
			friendsIds = append(friendsIds, friend.FriendID)
		} else {
			friendsIds = append(friendsIds, friend.UserID)
		}
	}

	c.JSON(http.StatusOK, gin.H{"friends": friendsIds})
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

	var newFriendship models.Friendship
	if err := s.db.Where("user_id = ? AND friend_id = ? AND Status = ?", req.FriendId, userId, "pending").First(&newFriendship).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "friend request not found"})
		return
	}
	newFriendship.Status = "accepted"
	if err := s.db.Save(&newFriendship).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update status"})
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
