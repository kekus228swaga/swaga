package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"github.com/kekus228swaga/orderflow/auth-service/internal/domain/user"
	"github.com/kekus228swaga/orderflow/auth-service/internal/service"
)

type AuthHandler struct {
	service   *service.AuthService
	jwtSecret string
}

func NewAuthHandler(service *service.AuthService, jwtSecret string) *AuthHandler {
	return &AuthHandler{service: service, jwtSecret: jwtSecret}
}

// Register обрабатывает POST /auth/register
func (h *AuthHandler) Register(c *gin.Context) {
	var req user.RegisterRequest
	// Валидация JSON (используем теги binding из model.go)
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// Вызываем сервис
	createdUser, err := h.service.Register(c.Request.Context(), req)
	if err != nil {
		// Простая обработка ошибок (в проде лучше проверить тип ошибки)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":    createdUser.ID,
		"email": createdUser.Email,
	})
}

// Login обрабатывает POST /auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req user.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// Проверка логина/пароля
	loggedUser, err := h.service.Login(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// Генерация JWT токена
	token, err := h.generateToken(loggedUser.ID, loggedUser.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "token generation failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":    loggedUser.ID,
			"email": loggedUser.Email,
		},
	})
}

// generateToken - вспомогательная функция для создания JWT
func (h *AuthHandler) generateToken(userID int64, email string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Токен живёт 24 часа
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(h.jwtSecret))
}
