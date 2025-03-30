package handlers

import (
	"encoding/json"
	"fmt"
	"heydays/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/streadway/amqp"
	"gorm.io/gorm"
)

type ChatServer struct {
	db            *gorm.DB
	redisClient   *redis.Client
	rabbitManager *RabbitMQManager
}

func NewChatServer(db *gorm.DB, redisClient *redis.Client, rabbitManager *RabbitMQManager) *ChatServer {
	return &ChatServer{
		db:            db,
		redisClient:   redisClient,
		rabbitManager: rabbitManager,
	}
}

func (s *ChatServer) SendMessage(c *gin.Context) {
	userId, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	type MessageInput struct {
		ChatID  uint   `json:"chat_id" binding:"required"`
		Content string `json:"content" binding:"required"`
	}

	var input MessageInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var chatUser models.ChatUser
	if err := s.db.Where("chat_id = ? AND user_id = ?", input.ChatID, userId).First(&chatUser).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not a member of this chat"})
		return
	}

	message := models.Message{
		ChatID:   input.ChatID,
		SenderID: userId.(uint),
		Content:  input.Content,
		IsRead:   false,
	}

	if err := s.db.Create(&message).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save message"})
		return
	}

	messageBody, _ := json.Marshal(message)

	err := s.rabbitManager.PublishMessage(
		"chat_messages",       // exchange
		"chat_messages_queue", // routing key
		false,                 // mandatory
		false,                 // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        messageBody,
		},
	)

	if err != nil {
		fmt.Printf("Warning: Failed to queue message in RabbitMQ: %v\n", err)
	}

	redisKey := fmt.Sprintf("chat:%d:messages", input.ChatID)
	err = s.redisClient.LPush(c.Request.Context(), redisKey, messageBody).Err()
	if err != nil {
		fmt.Printf("Failed to cache message in Redis: %v\n", err)
	}

	s.redisClient.LTrim(c.Request.Context(), redisKey, 0, 99)

	lastMsgKey := fmt.Sprintf("chat:%d:last_message", input.ChatID)
	s.redisClient.Set(c.Request.Context(), lastMsgKey, messageBody, 24*time.Hour)

	c.JSON(http.StatusOK, gin.H{
		"message": "Message sent successfully",
		"data":    message,
	})
}

func (s *ChatServer) GetChatMessages(c *gin.Context) {
	userId, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	chatID, err := strconv.Atoi(c.Param("chatId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chat ID"})
		return
	}
	var chatUser models.ChatUser
	if err := s.db.Where("chat_id = ? AND user_id = ?", chatID, userId).First(&chatUser).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not a member of this chat"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset := (page - 1) * limit

	var messages []models.Message
	var total int64

	s.db.Model(&models.Message{}).Where("chat_id = ?", chatID).Count(&total)

	result := s.db.Where("chat_id = ?", chatID).
		Preload("Attachments").
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&messages)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch messages"})
		return
	}

	currentTime := time.Now()
	s.db.Model(&models.Message{}).
		Where("chat_id = ? AND sender_id != ? AND is_read = ?", chatID, userId, false).
		Updates(map[string]interface{}{
			"is_read": true,
			"read_at": currentTime,
		})

	c.JSON(http.StatusOK, gin.H{
		"messages": messages,
		"total":    total,
		"page":     page,
		"limit":    limit,
	})
}

func (s *ChatServer) GetUserChats(c *gin.Context) {
	userId, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	type ChatPreview struct {
		ID              uint      `json:"id"`
		Name            string    `json:"name"`
		IsGroup         bool      `json:"is_group"`
		LastMessage     string    `json:"last_message"`
		LastSenderID    uint      `json:"last_sender_id"`
		LastSenderName  string    `json:"last_sender_name"`
		LastMessageTime time.Time `json:"last_message_time"`
		UnreadCount     int       `json:"unread_count"`
		Participants    []struct {
			ID       uint   `json:"id"`
			Name     string `json:"name"`
			Surname  string `json:"surname"`
			PhotoURL string `json:"photo_url"`
		} `json:"participants,omitempty"`
	}

	var userChats []models.ChatUser
	if err := s.db.Where("user_id = ?", userId).
		Preload("Chat").
		Find(&userChats).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch chats"})
		return
	}

	var chatPreviews []ChatPreview
	for _, uc := range userChats {
		chat := uc.Chat

		var lastMessage models.Message
		s.db.Where("chat_id = ?", chat.ID).Order("created_at DESC").Limit(1).First(&lastMessage)

		var lastSenderName string
		if lastMessage.ID != 0 {
			var sender models.User
			s.db.Select("name, surname").First(&sender, lastMessage.SenderID)
			lastSenderName = sender.Name + " " + sender.Surname
		}

		var unreadCount int64
		s.db.Model(&models.Message{}).
			Where("chat_id = ? AND sender_id != ? AND is_read = ?", chat.ID, userId, false).
			Count(&unreadCount)

		preview := ChatPreview{
			ID:              chat.ID,
			IsGroup:         chat.IsGroup,
			LastMessage:     lastMessage.Content,
			LastSenderID:    lastMessage.SenderID,
			LastSenderName:  lastSenderName,
			LastMessageTime: lastMessage.CreatedAt,
			UnreadCount:     int(unreadCount),
		}

		if chat.IsGroup {
			preview.Name = chat.Name
		} else {
			var otherUser models.User
			s.db.Table("chat_users").
				Select("users.id, users.name, users.surname, users.profile_photo").
				Joins("JOIN users ON users.id = chat_users.user_id").
				Where("chat_users.chat_id = ? AND chat_users.user_id != ?", chat.ID, userId).
				First(&otherUser)

			preview.Name = otherUser.Name + " " + otherUser.Surname

			preview.Participants = append(preview.Participants, struct {
				ID       uint   `json:"id"`
				Name     string `json:"name"`
				Surname  string `json:"surname"`
				PhotoURL string `json:"photo_url"`
			}{
				ID:       otherUser.ID,
				Name:     otherUser.Name,
				Surname:  otherUser.Surname,
				PhotoURL: otherUser.ProfilePhoto,
			})
		}

		chatPreviews = append(chatPreviews, preview)
	}

	c.JSON(http.StatusOK, gin.H{
		"chats": chatPreviews,
	})
}

func (s *ChatServer) CreateDirectChat(c *gin.Context) {
	userId, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	type CreateChatInput struct {
		UserID uint `json:"user_id" binding:"required"`
	}

	var input CreateChatInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var otherUser models.User
	if err := s.db.First(&otherUser, input.UserID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var existingChatID uint
	err := s.db.Raw(`
		SELECT c1.chat_id 
		FROM chat_users c1 
		JOIN chat_users c2 ON c1.chat_id = c2.chat_id 
		JOIN chats ON chats.id = c1.chat_id 
		WHERE c1.user_id = ? AND c2.user_id = ? AND chats.is_group = false`,
		userId, input.UserID).Scan(&existingChatID).Error

	if err == nil && existingChatID > 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "Chat already exists",
			"chat_id": existingChatID,
		})
		return
	}
	newChat := models.Chat{
		IsGroup: false,
		OwnerID: userId.(uint),
	}

	if err := s.db.Create(&newChat).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create chat"})
		return
	}

	chatUsers := []models.ChatUser{
		{
			ChatID:   newChat.ID,
			UserID:   userId.(uint),
			JoinedAt: time.Now(),
			IsAdmin:  true,
		},
		{
			ChatID:   newChat.ID,
			UserID:   input.UserID,
			JoinedAt: time.Now(),
			IsAdmin:  false,
		},
	}

	if err := s.db.Create(&chatUsers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add users to chat"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Chat created successfully",
		"chat_id": newChat.ID,
	})
}

func (s *ChatServer) CreateGroupChat(c *gin.Context) {
	userId, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	type CreateGroupChatInput struct {
		Name    string `json:"name" binding:"required"`
		UserIDs []uint `json:"user_ids" binding:"required"`
	}

	var input CreateGroupChatInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newChat := models.Chat{
		Name:    input.Name,
		IsGroup: true,
		OwnerID: userId.(uint),
	}

	if err := s.db.Create(&newChat).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create group chat"})
		return
	}

	ownerChatUser := models.ChatUser{
		ChatID:   newChat.ID,
		UserID:   userId.(uint),
		JoinedAt: time.Now(),
		IsAdmin:  true,
	}

	if err := s.db.Create(&ownerChatUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add owner to chat"})
		return
	}

	for _, userID := range input.UserIDs {
		var user models.User
		if err := s.db.First(&user, userID).Error; err != nil {
			continue
		}

		chatUser := models.ChatUser{
			ChatID:   newChat.ID,
			UserID:   userID,
			JoinedAt: time.Now(),
			IsAdmin:  false,
		}

		s.db.Create(&chatUser)
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Group chat created successfully",
		"chat_id": newChat.ID,
	})
}

func (s *ChatServer) GetChatInfo(c *gin.Context) {
	userId, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	chatID, err := strconv.Atoi(c.Param("chatId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chat ID"})
		return
	}

	var chatUser models.ChatUser
	if err := s.db.Where("chat_id = ? AND user_id = ?", chatID, userId).First(&chatUser).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not a member of this chat"})
		return
	}

	var chat models.Chat
	if err := s.db.First(&chat, chatID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Chat not found"})
		return
	}

	type Participant struct {
		ID       uint      `json:"id"`
		Name     string    `json:"name"`
		Surname  string    `json:"surname"`
		PhotoURL string    `json:"photo_url"`
		IsAdmin  bool      `json:"is_admin"`
		JoinedAt time.Time `json:"joined_at"`
	}

	var participants []Participant
	if err := s.db.Table("chat_users").
		Select("users.id, users.name, users.surname, users.profile_photo, chat_users.is_admin, chat_users.joined_at").
		Joins("JOIN users ON users.id = chat_users.user_id").
		Where("chat_users.chat_id = ?", chatID).
		Scan(&participants).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch participants"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":           chat.ID,
		"name":         chat.Name,
		"is_group":     chat.IsGroup,
		"owner_id":     chat.OwnerID,
		"created_at":   chat.CreatedAt,
		"participants": participants,
	})
}

func (s *ChatServer) ArchiveOldMessages() {
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		threeDaysAgo := time.Now().AddDate(0, 0, -3)

		var chats []models.Chat
		s.db.Find(&chats)

		for _, chat := range chats {
			var oldMessages []models.Message
			s.db.Where("chat_id = ? AND created_at < ?", chat.ID, threeDaysAgo).
				Find(&oldMessages)
			for _, msg := range oldMessages {
				archivedMessage := models.ArchivedMessage{
					OriginalMessageID: msg.ID,
					ChatID:            msg.ChatID,
					Content:           msg.Content,
					SenderID:          msg.SenderID,
					ArchivedAt:        time.Now(),
				}
				s.db.Create(&archivedMessage)
				s.db.Delete(&msg)
			}
		}
	}
}
