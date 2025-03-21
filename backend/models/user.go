package models

import (
	"html"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email        string `gorm:"size:50;uniqueIndex;not null"`
	Name         string `gorm:"size:32;not null"`
	Surname      string `gorm:"size:32;not null"`
	ProfilePhoto string `gorm:"default: '"`
	Password     string `gorm:"not null"`

	Friends        []*User    `gorm:"many2many:friendships;joinForeignKey:UserID;joinReferences:FriendID"`
	Posts          []Post     `gorm:"foreignKey:UserID"`
	SentMessages   []Message  `gorm:"foreignKey:SenderID"`
	OwnedChats     []Chat     `gorm:"foreignKey:OwnerID"`
	Participations []ChatUser `gorm:"foreignKey:UserID"`
}

type Friendship struct {
	gorm.Model
	UserID   uint   `gorm:"uniqueIndex:idx_friendship"`
	FriendID uint   `gorm:"uniqueIndex:idx_friendship"`
	Status   string `gorm:"default:'pending'"` // 'pending', 'accepted', 'rejected'
}

type Post struct {
	gorm.Model
	UserID  uint   `gorm:"index"`
	Content string `gorm:"type:text"`

	Images []Image `gorm:"foreignKey:PostID"`
}

type Image struct {
	gorm.Model
	PostID uint   `gorm:"index"`
	URL    string `gorm:"not null"`
}

func (user *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	user.Email = html.EscapeString(strings.TrimSpace(user.Email))

	user.Name = html.EscapeString(strings.TrimSpace(user.Name))

	user.Surname = html.EscapeString(strings.TrimSpace(user.Surname))

	return nil
}

func (user *User) VerifyPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
}
