package main

import (
	"heydays/handlers"
	"heydays/middleware"
	"heydays/models"
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

	openapi := r.Group("/openapi")
	openapi.GET("/profile/:userId", server.GetStrangerProfileInfo)

	auth := r.Group("/user")
	auth.Use(middleware.JWTMiddleware())
	auth.GET("/profile", server.GetProfileInfo)

	friends := r.Group("/friends")
	friends.Use(middleware.JWTMiddleware())
	friends.POST("/send-friend-request", server.SendFriendRequest)
	friends.PATCH("/accept-friend-request", server.AcceptFriendRequest)
	friends.DELETE("/decline-friend-request", server.DeclineFriendRequest)
	friends.GET("/get-friends", server.GetFriendList)
	friends.GET("/get-pending-friend-invies", server.GetPendingFriendInvies)

	return r

}

func main() {
	r := SetupRouter()

	r.Run(":8080")
}
