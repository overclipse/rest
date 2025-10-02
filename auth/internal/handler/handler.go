package handler

import (
	"auth/internal/config"
	"auth/internal/errors"
	"auth/internal/models"
	"auth/internal/service"
	"context"
	"time"

	jwtmanager "jwt_manager"

	"github.com/gin-gonic/gin"
)

// Handler содержит все обработчики для работы с пользователями
type Handler struct {
	service    service.Service        // Сервис для работы с БД
	jwtManager *jwtmanager.JWTManager // JWT менеджер для работы с токенами
	cfg        *config.Config         // Конфигурация сервера
}

// NewHandler создает новый экземпляр обработчика пользователей
func NewHandler(service service.Service, cfg *config.Config) *Handler {

	// Создаем JWT менеджер
	jwtConfig := jwtmanager.JWTConfig{
		SecretKey:              cfg.JWTSecretKey,
		AccessTokenExpiration:  cfg.AccessTokenExpiration,
		RefreshTokenExpiration: cfg.RefreshTokenExpiration,
	}
	jwtManager := jwtmanager.NewJWTManager(jwtConfig)

	return &Handler{
		service:    service,    // Сохраняем сервис в обработчике
		jwtManager: jwtManager, // Сохраняем JWT менеджер в обработчике
		cfg:        cfg,        // Сохраняем конфигурацию в обработчике
	}
}

// RegisterUser обрабатывает запрос на регистрацию нового пользователя
func (h *Handler) RegisterUser(c *gin.Context) {
	var user models.User

	// Парсим JSON из тела запроса
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{
			"error":   errors.MsgInvalidData,
			"details": err.Error(),
		})
		return
	}

	// Проверяем обязательные поля
	if user.Username == "" || user.Password == "" {
		c.JSON(400, gin.H{
			"error": errors.MsgInvalidUserData,
		})
		return
	}

	// Создаем пользователя в базе данных с таймаутом
	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Duration(h.cfg.DBTimeout)*time.Second)
	defer cancel()

	createdUser, err := h.service.Create(ctx, &user)
	if err != nil {
		c.JSON(500, gin.H{
			"error":   errors.MsgUserCreation,
			"details": err.Error(),
		})
		return
	}

	// Убираем пароль из ответа
	createdUser.Password = ""

	c.JSON(201, gin.H{
		"message": errors.MsgUserRegistered,
		"user":    createdUser,
	})
}

// LoginUser обрабатывает запрос на авторизацию пользователя
func (h *Handler) LoginUser(c *gin.Context) {
	var loginRequest struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	// Парсим JSON из тела запроса
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(400, gin.H{
			"error":   errors.MsgInvalidData,
			"details": err.Error(),
		})
		return
	}

	// Аутентифицируем пользователя с таймаутом
	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Duration(h.cfg.DBTimeout)*time.Second)
	// Отменяем контекст после завершения работы функции
	defer cancel()

	user, err := h.service.Authenticate(ctx, loginRequest.Username, loginRequest.Password)
	if err != nil {
		c.JSON(401, gin.H{
			"error": errors.MsgInvalidCredentials,
		})
		return
	}

	// Убираем пароль из ответа
	user.Password = ""

	// Генерируем JWT токены
	accessToken, refreshToken, err := h.jwtManager.GenerateTokens(user.ID)
	if err != nil {
		c.JSON(500, gin.H{
			"error":   errors.MsgTokenGeneration,
			"details": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message":       errors.MsgLoginSuccess,
		"user":          user,
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

// GetUserInfo обрабатывает запрос на получение информации о пользователе
func (h *Handler) GetUserInfo(c *gin.Context) {
	// Получаем ID пользователя из токена
	userID, err := h.GetCurrentUserID(c)
	if err != nil {
		c.JSON(401, gin.H{
			"error": errors.MsgAuthRequired,
		})
		return
	}

	// Аутентифицируем пользователя с таймаутом
	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Duration(h.cfg.DBTimeout)*time.Second)
	// Отменяем контекст после завершения работы функции
	defer cancel()

	user, err := h.service.Read(ctx, userID)
	if err != nil {
		c.JSON(404, gin.H{
			"error": errors.MsgUserNotFound,
		})
		return
	}

	// Убираем пароль из ответа
	user.Password = ""

	c.JSON(200, gin.H{
		"user": user,
	})
}

// GetCurrentUserID получает ID текущего пользователя из контекста (для совместимости)
func (h *Handler) GetCurrentUserID(c *gin.Context) (int, error) {
	return jwtmanager.GetCurrentUserID(c)
}

// UpdateUser обрабатывает запрос на обновление данных пользователя
func (h *Handler) UpdateUser(c *gin.Context) {
	// Получаем ID пользователя из токена
	userID, err := h.GetCurrentUserID(c)
	if err != nil {
		c.JSON(401, gin.H{
			"error": errors.MsgAuthRequired,
		})
		return
	}

	var updateData models.User

	// Парсим JSON из тела запроса
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(400, gin.H{
			"error":   errors.MsgInvalidData,
			"details": err.Error(),
		})
		return
	}

	// Устанавливаем ID для обновления
	updateData.ID = userID

	// Аутентифицируем пользователя с таймаутом
	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Duration(h.cfg.DBTimeout)*time.Second)
	// Отменяем контекст после завершения работы функции
	defer cancel()

	err = h.service.Update(ctx, &updateData)
	if err != nil {
		c.JSON(500, gin.H{
			"error":   errors.MsgDatabaseOperation,
			"details": err.Error(),
		})
		return
	}

	// Получаем обновленного пользователя с таймаутом
	readCtx, readCancel := context.WithTimeout(c.Request.Context(), time.Duration(h.cfg.DBTimeout)*time.Second)
	defer readCancel()

	updatedUser, err := h.service.Read(readCtx, userID)
	if err != nil {
		c.JSON(500, gin.H{
			"error": errors.MsgDatabaseOperation,
		})
		return
	}

	// Убираем пароль из ответа
	updatedUser.Password = ""

	c.JSON(200, gin.H{
		"message": errors.MsgUserUpdated,
		"user":    updatedUser,
	})
}

// DeleteUser обрабатывает запрос на удаление пользователя
func (h *Handler) DeleteUser(c *gin.Context) {
	// Получаем ID пользователя из токена
	userID, err := h.GetCurrentUserID(c)
	if err != nil {
		c.JSON(401, gin.H{
			"error": errors.MsgAuthRequired,
		})
		return
	}

	// Проверяем, что пользователь существует с таймаутом
	checkCtx, checkCancel := context.WithTimeout(c.Request.Context(), time.Duration(h.cfg.DBTimeout)*time.Second)
	defer checkCancel()

	_, err = h.service.Read(checkCtx, userID)
	if err != nil {
		c.JSON(404, gin.H{
			"error": errors.MsgUserNotFound,
		})
		return
	}

	// Удаляем пользователя из базы данных с таймаутом
	deleteCtx, deleteCancel := context.WithTimeout(c.Request.Context(), time.Duration(h.cfg.DBTimeout)*time.Second)
	defer deleteCancel()

	err = h.service.Delete(deleteCtx, userID)
	if err != nil {
		c.JSON(500, gin.H{
			"error":   errors.MsgDatabaseOperation,
			"details": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": errors.MsgUserDeleted,
	})
}

// RefreshToken обрабатывает запрос на обновление access токена с помощью refresh токена
func (h *Handler) RefreshToken(c *gin.Context) {
	var refreshRequest struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	// Парсим JSON из тела запроса
	if err := c.ShouldBindJSON(&refreshRequest); err != nil {
		c.JSON(400, gin.H{
			"error":   errors.MsgInvalidData,
			"details": err.Error(),
		})
		return
	}

	// Валидируем refresh токен
	userID, err := h.jwtManager.ValidateRefreshToken(refreshRequest.RefreshToken)
	if err != nil {
		c.JSON(401, gin.H{
			"error": errors.MsgRefreshToken,
		})
		return
	}

	// Проверяем, что пользователь существует
	user, err := h.service.Read(c.Request.Context(), userID)
	if err != nil {
		c.JSON(404, gin.H{
			"error": errors.MsgUserNotFound,
		})
		return
	}

	// Генерируем новые токены
	accessToken, refreshToken, err := h.jwtManager.GenerateTokens(user.ID)
	if err != nil {
		c.JSON(500, gin.H{
			"error":   errors.MsgTokenGeneration,
			"details": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message":       errors.MsgTokensRefreshed,
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

// ExtractTokenFromHeader извлекает JWT токен из HTTP заголовка Authorization (публичный метод)
func (h *Handler) ExtractTokenFromHeader(c *gin.Context) (string, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return "", errors.ErrMissingAuthHeader
	}

	// Проверяем на формат Bearer
	const bearerPrefix = "Bearer "
	if len(authHeader) <= len(bearerPrefix) || authHeader[:len(bearerPrefix)] != bearerPrefix {
		return "", errors.ErrInvalidAuthFormat
	}

	// Возвращаем токен без префикса "Bearer "
	return authHeader[len(bearerPrefix):], nil
}

// ValidateAccessToken валидирует access токен (публичный метод)
func (h *Handler) ValidateAccessToken(tokenString string) (int, error) {
	return h.jwtManager.ValidateAccessToken(tokenString)
}

// RequireAuth возвращает middleware для проверки JWT токена
func (h *Handler) RequireAuth() gin.HandlerFunc {
	return h.jwtManager.JWTInterceptor()
}
