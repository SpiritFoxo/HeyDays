package handlers

import (
	"heydays/models"
	"heydays/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Server struct {
	db *gorm.DB
}

func NewServer(db *gorm.DB) *Server {
	return &Server{db: db}
}

func (s *Server) Register(c *gin.Context) {

	type RegisterInput struct {
		Email    string `json:"email" binding:"required"`
		Name     string `json:"name" binding:"required"`
		Surname  string `json:"surname" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	var input RegisterInput

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := models.User{Email: input.Email, Name: input.Name, Surname: input.Surname, Password: input.Password}
	user.HashPassword()

	if err := s.db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created"})
}

func (s *Server) LoginCheck(email, password string) (string, *models.User, error) {
	var err error

	user := models.User{}

	if err = s.db.Model(models.User{}).Where("email=?", email).Take(&user).Error; err != nil {
		return "", nil, err
	}

	err = user.VerifyPassword(password)

	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", nil, err
	}

	token, err := utils.GenerateToken(user)

	if err != nil {
		return "", nil, err
	}

	return token, &user, nil
}

func (s *Server) Login(c *gin.Context) {

	type LoginInput struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	var input LoginInput

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, user, err := s.LoginCheck(input.Email, input.Password)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "The email or password is not correct"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token, "user_id": user.ID})
}
