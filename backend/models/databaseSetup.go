package models

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Setup() (*gorm.DB, error) {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Cant load .env file", err)
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SSLMODE"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Cant connect to bd:", err)
	}

	if err = db.AutoMigrate(
		&User{},
		&Friendship{},
		&Post{},
		&Image{},
		&Chat{},
		&ChatUser{},
		&Message{},
		&MessageAttachment{},
	); err != nil {
		log.Println("cant migrate db:", err)
	}

	log.Println("Success")
	return db, nil
}
