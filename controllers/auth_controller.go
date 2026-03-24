package controllers

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"project-e-commerce/config"
	"project-e-commerce/models"
	"project-e-commerce/utils"
)

var jwtSecret = []byte(config.GetEnv("JWT_SECRET"))

type RegisterRequest struct {
	Username string `json:"username" example:"john"`
	Email    string `json:"email" example:"john@email.com"`
	Password string `json:"password" example:"123456"`
}

type LoginRequest struct {
	Username string `json:"username" example:"user"`
	Password string `json:"password" example:"1234"`
}

// Register godoc
// @Summary Register user
// @Description Create new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param body body RegisterRequest true "Register data"
// @Success 200 {object} map[string]interface{}
// @Router /register [post]
func Register(c *gin.Context) {
	var input RegisterRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.Error(utils.BadRequest("invalid input", err))
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), 10)
	if err != nil {
		c.Error(utils.Internal("failed to hash password", err))
		return
	}

	user := models.User{
		ID:       uuid.New(),
		Username: input.Username,
		Email:    input.Email,
		Password: string(hashed),
		Role:     "user",
	}

	if err := config.DB.Create(&user).Error; err != nil {
		c.Error(utils.BadRequest("username or email already exists", err))
		return
	}

	utils.Success(c, "User registered", nil)
}

// Login godoc
// @Summary Login user
// @Description Authenticate user and return JWT
// @Tags auth
// @Accept json
// @Produce json
// @Param body body LoginRequest true "Login credentials"
// @Success 200 {object} map[string]string
// @Router /login [post]
func Login(c *gin.Context) {
	var input LoginRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.Error(utils.BadRequest("invalid input", err))
		return
	}

	var user models.User
	if err := config.DB.Where("username = ?", input.Username).
		First(&user).Error; err != nil {

		c.Error(utils.BadRequest("invalid credentials", err))
		return
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(input.Password),
	); err != nil {

		c.Error(utils.BadRequest("invalid credentials", err))
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID.String(),
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		c.Error(utils.Internal("failed to generate token", err))
		return
	}

	utils.Success(c, "Login successful", gin.H{
		"token": tokenString,
	})
}

// Profile godoc
// @Summary Get user profile
// @Tags auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /admin/profile [get]
func Profile(c *gin.Context) {

	userID, exists := c.Get("user_id")
	if !exists {
		c.Error(utils.Internal("user id not found in context", nil))
		return
	}

	utils.Success(c, "You are authenticated", userID)
}
