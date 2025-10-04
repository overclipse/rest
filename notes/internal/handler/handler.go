package handler

import (
	"github.com/gin-gonic/gin"
	jwtmanager "jwt_manager"
	"notes/internal/config"
)

// Handler содержит все обработчики для работы с заметками
type Handler struct {
	cfg        *config.Config         // Конфигурация приложения
	jwtManager *jwtmanager.JWTManager // JWT менеджер для работы с токенами
}

// NewHandler создает новый экземпляр обработчика заметок
func NewHandler(cfg *config.Config) *Handler {
	// Создаем JWT менеджер
	jwtConfig := jwtmanager.JWTConfig{
		SecretKey:              cfg.JWTSecretKey,
		AccessTokenExpiration:  24,  // 24 часа по умолчанию
		RefreshTokenExpiration: 168, // 7 дней по умолчанию
	}
	jwtManager := jwtmanager.NewJWTManager(jwtConfig)

	return &Handler{
		cfg:        cfg,
		jwtManager: jwtManager,
	}
}

// CreateNote создает новую заметку
// POST /api/v1/note
func (h *Handler) CreateNote(c *gin.Context) {
	c.JSON(201, gin.H{
		"message": "Создание заметки не реализовано",
	})
}

// GetNoteByID получает заметку по ID
// GET /api/v1/note/:id
func (h *Handler) GetNoteByID(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Получение заметки по ID не реализовано",
	})
}

// UpdateNote обновляет существующую заметку
// PUT /api/v1/note/:id
func (h *Handler) UpdateNote(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Обновление заметки не реализовано",
	})
}

// DeleteNote удаляет заметку по ID
// DELETE /api/v1/note/:id
func (h *Handler) DeleteNote(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Удаление заметки не реализовано",
	})
}

// GetAllNotes получает список всех заметок текущего пользователя
// GET /api/v1/notes
func (h *Handler) GetAllNotes(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Получение всех заметок не реализовано",
	})
}

// GetJWTMiddleware возвращает JWT middleware для использования в роутах
func (h *Handler) GetJWTMiddleware() gin.HandlerFunc {
	return h.jwtManager.JWTInterceptor()
}