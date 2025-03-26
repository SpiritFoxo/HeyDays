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
	db          *gorm.DB
	redisClient *redis.Client
	rabbitConn  *amqp.Connection
	rabbitChan  *amqp.Channel
}

func NewChatServer(db *gorm.DB, redisClient *redis.Client, rabbitConn *amqp.Connection) *ChatServer {
	// Create RabbitMQ channel
	ch, err := rabbitConn.Channel()
	if err != nil {
		panic(fmt.Sprintf("Failed to open a channel: %v", err))
	}

	// Declare exchanges and queues
	err = ch.ExchangeDeclare(
		"chat_messages", // name
		"direct",        // type
		true,            // durable
		false,           // auto-deleted
		false,           // internal
		false,           // no-wait
		nil,             // arguments
	)
	if err != nil {
		panic(fmt.Sprintf("Failed to declare exchange: %v", err))
	}

	_, err = ch.QueueDeclare(
		"chat_messages_queue", // name
		true,                  // durable
		false,                 // delete when unused
		false,                 // exclusive
		false,                 // no-wait
		nil,                   // arguments
	)
	if err != nil {
		panic(fmt.Sprintf("Failed to declare queue: %v", err))
	}

	return &ChatServer{
		db:          db,
		redisClient: redisClient,
		rabbitConn:  rabbitConn,
		rabbitChan:  ch,
	}
}

func (s *ChatServer) SendMessage(c *gin.Context) {
	// Get user ID from context
	userId, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Define message input struct
	type MessageInput struct {
		ChatID  uint   `json:"chat_id" binding:"required"`
		Content string `json:"content" binding:"required"`
	}

	var input MessageInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create message model
	message := models.Message{
		ChatID:   input.ChatID,
		SenderID: userId.(uint),
		Content:  input.Content,
		IsRead:   false,
	}

	// Publish to RabbitMQ
	messageBody, _ := json.Marshal(message)
	err := s.rabbitChan.Publish(
		"chat_messages",
		"chat_messages_queue",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        messageBody,
		},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to queue message"})
		return
	}

	redisKey := fmt.Sprintf("chat:%d:messages", input.ChatID)
	err = s.redisClient.LPush(c.Request.Context(), redisKey, messageBody).Err()
	if err != nil {
		fmt.Printf("Failed to cache message in Redis: %v\n", err)
	}

	s.redisClient.LTrim(c.Request.Context(), redisKey, 0, 99)

	c.JSON(http.StatusOK, gin.H{"message": "Message sent successfully"})
}

func (s *ChatServer) GetChatMessages(c *gin.Context) {
	chatID, err := strconv.Atoi(c.Param("chatId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chat ID"})
		return
	}

	redisKey := fmt.Sprintf("chat:%d:messages", chatID)
	cachedMessages, err := s.redisClient.LRange(c.Request.Context(), redisKey, 0, 99).Result()

	var messages []models.Message
	if err == nil && len(cachedMessages) > 0 {
		for _, msg := range cachedMessages {
			var message models.Message
			json.Unmarshal([]byte(msg), &message)
			messages = append(messages, message)
		}
	} else {
		threeDaysAgo := time.Now().AddDate(0, 0, -3)
		s.db.Where("chat_id = ? AND created_at >= ?", chatID, threeDaysAgo).
			Order("created_at DESC").
			Limit(100).
			Find(&messages)
	}

	c.JSON(http.StatusOK, gin.H{"messages": messages})
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
				}
				s.db.Create(&archivedMessage)
				s.db.Delete(&msg)
			}
		}
	}
}
