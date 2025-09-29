package handler

import (
	"auth/internal/config"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	cfg *config.Config
}

// Создание нового обработчика пользователей
func NewHandler (cfg *config.Config) *Handler {
	return &Handler{
		cfg: cfg, //Сохраняем конфигурацию в обработчике
	}
}


// RegisterUser обрабатывает запрос на регистрацию нового пользователя
func (h *Handler) RegisterUser(c *gin.Context) {
	c.JSON(201, gin.H{
		"message": "Пользователь успешно зарегистрирован",
	})
}

// LoginUser обрабатывает запрос на авторизацию пользователя
func (h *Handler) LoginUser(c *gin.Context) {
	c.JSON(200, gin.H{
		"message":       "Пользователь успешно авторизован",
		"access_token":  "example_access_token",
		"refresh_token": "example_refresh_token",
	})
}

// GetUserInfo обрабатывает запрос на получение информации о пользователе
func (h *Handler) GetUserInfo(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Информация о пользователе успешно получена",
	})
}

// UpdateUser обрабатывает запрос на обновление данных пользователя
func (h *Handler) UpdateUser(c *gin.Context) {

	c.JSON(200, gin.H{
		"message": "Данные пользователя успешно обновлены",
	})
}

// DeleteUser обрабатывает запрос на удаление пользователя
func (h *Handler) DeleteUser(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Пользователь успешно удален",
	})
}