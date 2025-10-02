package jwtmanager

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// JWTInterceptor создает middleware для проверки JWT токена
// Этот middleware можно использовать в любом сервисе с Gin роутером
func (j *JWTManager) JWTInterceptor() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Извлекаем токен из заголовка
		tokenString, err := j.extractTokenFromHeader(c)
		if err != nil {
			c.JSON(401, gin.H{
				"error": MsgTokenRequired,
			})
			c.Abort()
			return
		}

		// Валидируем токен
		userID, err := j.ValidateAccessToken(tokenString)
		if err != nil {
			c.JSON(401, gin.H{
				"error": MsgInvalidToken,
			})
			c.Abort()
			return
		}

		// Сохраняем ID пользователя в контексте
		c.Set("user_id", userID)
		c.Next()
	}
}

// extractTokenFromHeader извлекает JWT токен из HTTP заголовка Authorization
func (j *JWTManager) extractTokenFromHeader(c *gin.Context) (string, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return "", ErrMissingAuthHeader
	}

	// Проверяем на формат Bearer
	const bearerPrefix = "Bearer "
	if len(authHeader) <= len(bearerPrefix) || !strings.HasPrefix(authHeader, bearerPrefix) {
		return "", ErrInvalidAuthFormat
	}

	// Возвращаем токен без префикса "Bearer "
	return authHeader[len(bearerPrefix):], nil
}

// GetCurrentUserID получает ID текущего пользователя из контекста Gin
// Этот метод можно использовать в handler'ах для получения ID аутентифицированного пользователя
func GetCurrentUserID(c *gin.Context) (int, error) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, ErrMissingUserID
	}

	id, ok := userID.(int)
	if !ok {
		return 0, ErrMissingUserID
	}

	return id, nil
}
