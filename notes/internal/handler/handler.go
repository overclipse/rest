package handler

import (
	"context"
	jwtmanager "jwt_manager"
	"net/http"
	"notes/internal/config"
	"notes/internal/errors"
	"notes/internal/models"
	"notes/internal/service"

	"github.com/gin-gonic/gin"
)

// Handler содержит все обработчики для работы с заметками
type Handler struct {
	cfg        *config.Config         // Конфигурация приложения
	jwtManager *jwtmanager.JWTManager // JWT менеджер для работы с токенами
	service    service.Service        // Сервис для работы с заметками
}

// NewHandler создает новый экземпляр обработчика заметок
func NewHandler(cfg *config.Config, service service.Service) *Handler {
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
		service:    service,
	}
}

// GetJWTMiddleware возвращает JWT middleware для использования в роутах
func (h *Handler) GetJWTMiddleware() gin.HandlerFunc {
	return h.jwtManager.JWTInterceptor()
}

// extractAuthorID извлекает ID автора из контекста (установленного JWT middleware)
func (h *Handler) extractAuthorID(c *gin.Context) (int, error) {
	// Используем готовую функцию из JWT manager
	return jwtmanager.GetCurrentUserID(c)
}

// CreateNote создает новую заметку
// POST /note
func (h *Handler) CreateNote(c *gin.Context) {
	// Извлекаем ID автора из JWT токена
	authorID, err := h.extractAuthorID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   errors.MsgMissingUserID,
			"details": err.Error(),
		})
		return
	}
	// Проверяем, что тело запроса содержит корректные данные
	var note models.Note
	if err := c.ShouldBindJSON(&note); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   errors.MsgInvalidData,
			"details": err.Error(),
		})
		return
	}

	// Устанавливаем ID автора из токена
	note.AuthorID = authorID
	// Создаем контекст для работы с сервисом
	ctx := context.Background()
	// Вызываем сервис для создания заметки
	createdNote, err := h.service.Create(ctx, note)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   errors.MsgNoteCreation,
			"details": err.Error(),
		})
		return
	}
	// Возвращаем успешный ответ с созданной заметкой
	c.JSON(http.StatusCreated, gin.H{
		"message": errors.MsgNoteCreated,
		"note":    createdNote,
	})
}

// GetNoteByID получает заметку по ID
// GET /note/:id
func (h *Handler) GetNoteByID(c *gin.Context) {
	// Извлекаем ID автора из JWT токена
	authorID, err := h.extractAuthorID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   errors.MsgMissingUserID,
			"details": err.Error(),
		})
		return
	}
	// Извлекаем ID заметки из параметров запроса
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errors.MsgInvalidNoteID,
		})
		return
	}
	// Создаем контекст для работы с сервисом
	ctx := context.Background()
	// Вызываем сервис для получения заметки по ID
	note, err := h.service.GetByID(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   errors.MsgNoteNotFound,
			"details": err.Error(),
		})
		return
	}

	// Проверяем, что пользователь является владельцем заметки
	if note.AuthorID != authorID {
		c.JSON(http.StatusForbidden, gin.H{})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": errors.MsgNoteFound,
		"note":    note,
	})
}

// UpdateNote обновляет существующую заметку
// PUT /note/:id
func (h *Handler) UpdateNote(c *gin.Context) {
	// Извлекаем ID автора из JWT токена
	authorID, err := h.extractAuthorID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   errors.MsgMissingUserID,
			"details": err.Error(),
		})
		return
	}
	// Извлекаем ID заметки из параметров запроса
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errors.MsgInvalidNoteID,
		})
		return
	}

	// Сначала проверяем, что заметка существует и принадлежит пользователю
	ctx := context.Background()
	existingNote, err := h.service.GetByID(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   errors.MsgNoteNotFound,
			"details": err.Error(),
		})
		return
	}

	// Проверяем владельца
	if existingNote.AuthorID != authorID {
		c.JSON(http.StatusForbidden, gin.H{})
		return
	}
	// Проверяем, что тело запроса содержит корректные данные
	var note models.Note
	if err := c.ShouldBindJSON(&note); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   errors.MsgInvalidData,
			"details": err.Error(),
		})
		return
	}
	// Устанавливаем ID заметки и ID автора
	note.ID = id
	note.AuthorID = authorID // Убеждаемся, что автор не изменился
	// Вызываем сервис для обновления заметки
	updatedNote, err := h.service.Update(ctx, note)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   errors.MsgNoteUpdate,
			"details": err.Error(),
		})
		return
	}
	// Возвращаем успешный ответ с обновленной заметкой
	c.JSON(http.StatusOK, gin.H{
		"message": errors.MsgNoteUpdated,
		"note":    updatedNote,
	})
}

// DeleteNote удаляет заметку по ID
// DELETE /note/:id
func (h *Handler) DeleteNote(c *gin.Context) {
	// Извлекаем ID автора из JWT токена
	authorID, err := h.extractAuthorID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   errors.MsgMissingUserID,
			"details": err.Error(),
		})
		return
	}
	// Извлекаем ID заметки из параметров запроса
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errors.MsgInvalidNoteID,
		})
		return
	}

	// Сначала проверяем, что заметка существует и принадлежит пользователю
	ctx := context.Background()
	// Получаем заметку по ID
	existingNote, err := h.service.GetByID(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   errors.MsgNoteNotFound,
			"details": err.Error(),
		})
		return
	}

	// Проверяем владельца
	if existingNote.AuthorID != authorID {
		c.JSON(http.StatusForbidden, gin.H{})
		return
	}
	// Вызываем сервис для удаления заметки
	err = h.service.Delete(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   errors.MsgNoteDeletion,
			"details": err.Error(),
		})
		return
	}
	// Возвращаем успешный ответ
	c.JSON(http.StatusOK, gin.H{
		"message": errors.MsgNoteDeleted,
	})
}

// GetAllNotes получает список всех заметок текущего пользователя
// GET /notes
func (h *Handler) GetAllNotes(c *gin.Context) {
	// Извлекаем ID автора из JWT токена
	authorID, err := h.extractAuthorID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   errors.MsgMissingUserID,
			"details": err.Error(),
		})
		return
	}
	// Создаем контекст для работы с сервисом
	ctx := context.Background()
	// Вызываем сервис для получения всех заметок текущего пользователя
	notes, err := h.service.GetAll(ctx, authorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   errors.MsgDatabaseOperation,
			"details": err.Error(),
		})
		return
	}
	// Проверяем, что заметки найдены
	c.JSON(http.StatusOK, gin.H{
		"message":   errors.MsgNotesFound,
		"notes":     notes,
		"count":     len(notes),
		"author_id": authorID,
	})
}