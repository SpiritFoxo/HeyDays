package models

import (
	"time"

	"gorm.io/gorm"
)

type Chat struct {
	gorm.Model
	Name    string // Empty for direct chats
	IsGroup bool   `gorm:"not null;default:false"`
	OwnerID uint   `gorm:"index"`

	Messages []Message  `gorm:"foreignKey:ChatID"`
	Users    []ChatUser `gorm:"foreignKey:ChatID"`
}

type ChatUser struct {
	gorm.Model
	ChatID   uint      `gorm:"uniqueIndex:idx_chat_user"`
	UserID   uint      `gorm:"uniqueIndex:idx_chat_user"`
	JoinedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	IsAdmin  bool      `gorm:"default:false"`
}

type Message struct {
	gorm.Model
	ChatID   uint   `gorm:"index"`
	SenderID uint   `gorm:"index"`
	Content  string `gorm:"type:text"`
	IsRead   bool   `gorm:"default:false"`
	ReadAt   *time.Time

	Attachments []MessageAttachment `gorm:"foreignKey:MessageID"`
}

type ArchivedMessage struct {
	gorm.Model
	OriginalMessageID uint
	ChatID            uint
	SenderID          uint
	Content           string
	ArchivedAt        time.Time
}

type MessageAttachment struct {
	gorm.Model
	MessageID uint   `gorm:"index"`
	URL       string `gorm:"not null"`
	Type      string `gorm:"not null"` // 'image', 'file', etc.
}
