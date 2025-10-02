// Package jwt предоставляет функционал для работы с JWT токенами
package jwtmanager

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

// Константы для типов токенов
const (
	ACCESS_TOKEN  = "accessToken"  // Тип токена для доступа к API
	REFRESH_TOKEN = "refreshToken" // Тип токена для обновления access токена
)

// JWTConfig представляет конфигурацию JWT
type JWTConfig struct {
	SecretKey              string // Секретный ключ для подписи JWT токенов
	AccessTokenExpiration  int    // Срок действия access токена в часах
	RefreshTokenExpiration int    // Срок действия refresh токена в часах
}

// JWTManager предоставляет функционал для работы с JWT токенами.
// Структура инкапсулирует логику создания, проверки и извлечения токенов.
type JWTManager struct {
	config JWTConfig // Конфигурация JWT
}

// NewJWTManager создает новый сервис для работы с JWT токенами.
func NewJWTManager(config JWTConfig) *JWTManager {
	return &JWTManager{
		config: config,
	}
}

// GenerateTokens генерирует пару токенов (access и refresh) для указанного ID пользователя.
// Возвращает строки токенов и ошибку, если генерация не удалась.
func (s *JWTManager) GenerateTokens(id int) (access string, refresh string, err error) {
	// Генерация access token
	accessTokenString, err := s.generateToken(id, ACCESS_TOKEN, s.config.AccessTokenExpiration)
	if err != nil {
		return "", "", fmt.Errorf("%s: %w", ErrTokenGeneration, err)
	}

	// Генерация refresh token
	refreshTokenString, err := s.generateToken(id, REFRESH_TOKEN, s.config.RefreshTokenExpiration)
	if err != nil {
		return "", "", fmt.Errorf("%s: %w", ErrTokenGeneration, err)
	}

	return accessTokenString, refreshTokenString, nil
}

// ValidateAccessToken проверяет корректность access токена.
// Возвращает ID пользователя из токена и ошибку валидации.
func (s *JWTManager) ValidateAccessToken(tokenString string) (int, error) {
	return s.validateToken(tokenString, ACCESS_TOKEN)
}

// ValidateRefreshToken проверяет корректность refresh токена.
// Возвращает ID пользователя из токена и ошибку валидации.
func (s *JWTManager) ValidateRefreshToken(tokenString string) (int, error) {
	return s.validateToken(tokenString, REFRESH_TOKEN)
}

// generateToken создает и подписывает токен заданного типа.
// Функция используется внутри сервиса для генерации как access, так и refresh токенов.
func (s *JWTManager) generateToken(id int, tokenType string, expirationHours int) (string, error) {
	now := time.Now()
	expiration := now.Add(time.Hour * time.Duration(expirationHours))

	// Создаем набор утверждений (claims) для токена
	claims := jwt.MapClaims{
		"id":   id,                // ID пользователя
		"type": tokenType,         // Тип токена (access или refresh)
		"iat":  now.Unix(),        // Время выпуска токена
		"exp":  expiration.Unix(), // Время истечения срока действия
	}

	// Создаем новый токен с выбранным алгоритмом подписи и утверждениями
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Подписываем токен секретным ключом
	tokenString, err := token.SignedString([]byte(s.config.SecretKey))
	if err != nil {
		return "", fmt.Errorf("%s: %w", ErrInvalidSignature, err)
	}

	return tokenString, nil
}

// validateToken - общая функция для валидации токенов.
// Проверяет подпись, срок действия и тип токена.
func (s *JWTManager) validateToken(tokenString, tokenType string) (int, error) {
	// Функция для проверки ключа подписи
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		// Проверка алгоритма подписи
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("%s: %v", ErrInvalidSignature, token.Header["alg"])
		}
		return []byte(s.config.SecretKey), nil
	}
	// Парсинг токена
	token, err := jwt.Parse(tokenString, keyFunc)
	if err != nil { // Проверяем тип ошибки, чтобы дать более точную информацию
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return 0, fmt.Errorf("%s: %w", ErrTokenExpired, err)
			}
		}
		return 0, fmt.Errorf("%s: %w", ErrInvalidToken, err)
	}

	// Проверка валидности токена и получение данных
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

		// Проверка типа токена, если требуется
		if tokenType != "" {
			if claimType, exists := claims["type"].(string); !exists || claimType != tokenType {
				return 0, fmt.Errorf("%s: ожидается %s, получен %s", ErrInvalidTokenType, tokenType, claims["type"])
			}
		}
		// Получение ID пользователя из токена
		idValue, exists := claims["id"].(float64)
		if !exists {
			return 0, fmt.Errorf("%s", ErrMissingUserID)
		}

		return int(idValue), nil
	}

	return 0, fmt.Errorf("%s", ErrInvalidToken)
}