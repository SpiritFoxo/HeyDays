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
