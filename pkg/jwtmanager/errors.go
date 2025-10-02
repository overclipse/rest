package jwtmanager

import "errors"

// JWT-специфичные ошибки
var (
	ErrInvalidToken      = errors.New("неверный токен")
	ErrTokenExpired      = errors.New("токен истек")
	ErrInvalidTokenType  = errors.New("неверный тип токена")
	ErrTokenGeneration   = errors.New("ошибка генерации токена")
	ErrMissingUserID     = errors.New("ID пользователя отсутствует в токене")
	ErrInvalidSignature  = errors.New("неверная подпись токена")
	ErrMissingMetadata   = errors.New("метаданные отсутствуют в контексте")
	ErrMissingAuthHeader = errors.New("отсутствует заголовок Authorization")
	ErrInvalidAuthFormat = errors.New("неверный формат токена")
)

// Сообщения для JWT ошибок
const (
	MsgTokenExpired         = "токен истек"
	MsgInvalidToken         = "неверный токен"
	MsgInvalidTokenType     = "неверный тип токена"
	MsgTokenGenerationError = "ошибка генерации токена"
	MsgMissingUserID        = "ID пользователя отсутствует в токене"
	MsgInvalidSignature     = "неверная подпись токена"
	MsgMissingMetadata      = "метаданные отсутствуют в контексте"
	MsgMissingAuthHeader    = "отсутствует заголовок Authorization"
	MsgInvalidAuthFormat    = "неверный формат токена"
	MsgTokenRequired        = "токен отсутствует или неверный формат"
)
