package main

import (
	"heydays/config"
	"heydays/handlers"
	"heydays/middleware"
	"heydays/models"
	"heydays/ws"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func DbInit() *gorm.DB {
	db, err := models.Setup()
	if err != nil {
		log.Println("Connection error")
	}
	return db
}

func SetupRouter() *gin.Engine {
	r := gin.Default()

	db := DbInit()

	redisClient := config.SetupRedis()

	rabbitManager := config.SetupRabbitMQManager()

	err := rabbitManager.Connect()
	if err != nil {
		log.Printf("Warning: Initial RabbitMQ connection failed: %v", err)
	}

	chatServer := handlers.NewChatServer(db, redisClient, rabbitManager)

	server := handlers.NewServer(db)

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router := r.Group("/auth")
	router.POST("/register", server.Register)
	router.POST("/login", server.Login)

	//openapi := r.Group("/openapi")
	//openapi.GET("/profile/:userId", server.GetStrangerProfileInfo)

	user := r.Group("/user")
	user.Use(middleware.JWTMiddleware())
	user.GET("/profile", server.GetProfileInfo)
	user.GET("/profile/:userId", server.GetStrangerProfileInfo)
	user.POST("/update-profile", server.UpdateProfileInfo)

	friends := r.Group("/friends")
	friends.Use(middleware.JWTMiddleware())
	friends.POST("/send-friend-request", server.SendFriendRequest)
	friends.PATCH("/accept-friend-request", server.AcceptFriendRequest)
	friends.DELETE("/decline-friend-request", server.DeclineFriendRequest)
	friends.GET("/get-friends", server.GetFriendList)
	friends.GET("/get-pending-friend-requests", server.GetPendingFriendInvies)

	chat := r.Group("/chat")
	chat.Use(middleware.JWTMiddleware())
	chat.POST("/send", chatServer.SendMessage)
	chat.GET("/:chatId/messages", chatServer.GetChatMessages)
	chat.GET("/list", chatServer.GetUserChats)
	chat.POST("/direct", chatServer.CreateDirectChat)
	chat.POST("/group", chatServer.CreateGroupChat)
	chat.GET("/:chatId", chatServer.GetChatInfo)

	r.GET("/ws", ws.HandleConnections)
	go ws.BroadcastMessages()

	go chatServer.ArchiveOldMessages()

	return r

}

func main() {
	r := SetupRouter()

	r.Run(":8080")
}
