package handlers

import (
	"heydays/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

//send frend request

func (s *Server) SendFriendRequest(c *gin.Context) {
	userId, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "cannot get user_id"})
		c.Abort()
		return
	}

	type request struct {
		friendId uint `json:"friend_id"`
	}

	var req request
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newFriendship := models.Friendship{UserID: userId.(uint), FriendID: req.friendId}
	if err := s.db.Create(&newFriendship); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot send friend request"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"messgae": "succes"})

}

//accept/reject friend request

func (s *Server) AcceptFriendRequest(c *gin.Context) {
	userId, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "cannot get user_id"})
		c.Abort()
		return
	}

	type request struct {
		friendId uint `json:"friend_id"`
	}

	var req request
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var newFriendship models.Friendship
	if err := s.db.Where("UserId = ? AND FriendId = ? AND Status = 'pending'", req.friendId, userId).First(&newFriendship).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "this relation does not exist"})
		return
	}
	newFriendship.Status = "accepted"
	if err := s.db.Save(&newFriendship).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "this relation does not exist"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "succes"})
}

// unfinished
func (s *Server) DeclineFriendRequest(c *gin.Context) {

	userId, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "cannot get user_id"})
		c.Abort()
		return
	}

	type request struct {
		friendId uint `json:"friend_id"`
	}

	var req request
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := s.db.Where("UserId = ? AND FriendId = ? AND Status = 'pending'", req.friendId, userId).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "this relation does not exist"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "succes"})
}
