package handler

import (
	"notes/internal/config"

	"github.com/gin-gonic/gin"
)

// Handler содержит все обработчики для работы с пользователями
type Handler struct {
	cfg *config.Config // Конфигурация сервера
}

// NewHandler создает новый экземпляр обработчика пользователей
func NewHandler(cfg *config.Config) *Handler {
	return &Handler{
		cfg: cfg, // Сохраняем конфигурацию в обработчике
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