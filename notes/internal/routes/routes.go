package routes

import (
	"notes/internal/handler"

	"github.com/gin-gonic/gin"
)

// SetupRouter инициализирует все роуты приложения
// Возвращает настроенный экземпляр gin.Engine
func SetupRouter(noteHandler *handler.Handler) *gin.Engine {
	// Инициализация роутера (по умолчанию)
	router := gin.Default()

	// Группа API для заметок с JWT авторизацией
	noteAPI := router.Group("/notes")
	noteAPI.Use(noteHandler.GetJWTMiddleware()) // Применяем JWT middleware ко всем роутам
	{
		// Создание заметки
		noteAPI.POST("/note", noteHandler.CreateNote)
		// Получение заметки по ID
		noteAPI.GET("/note/:id", noteHandler.GetNoteByID)
		// Редактирование заметки
		noteAPI.PUT("/note/:id", noteHandler.UpdateNote)
		// Удаление заметки
		noteAPI.DELETE("/note/:id", noteHandler.DeleteNote)
		// Получение списка всех заметок
		noteAPI.GET("/notes", noteHandler.GetAllNotes)
	}

	return router
}