package service

import (
	"auth/internal/models"
	"context"
)

// Service - интерфейс для управления пользователями
// Он определяет методы согласно паттерна CRUD для создания, поиска, обновления и удаления пользователей
// Все методы принимают контекст для управления временем выполнения и отмены операций
// Это позволяет гибко управлять жизненным циклом операций с пользователями
type Service interface {
	// Создает нового пользователя
	Create(ctx context.Context, user *models.User) (*models.User, error)
	// Находит пользователя по ID
	Read(ctx context.Context, id int) (*models.User, error)
	// Обновляет данные пользователя
	Update(ctx context.Context, user *models.User) error
	// Удаляет пользователя по ID
	Delete(ctx context.Context, id int) error
	// Аутентифицирует пользователя по имени пользователя и паролю
	Authenticate(ctx context.Context, username, password string) (*models.User, error)
	// Находит пользователя по имени пользователя
	ReadByUsername(ctx context.Context, username string) (*models.User, error)
	// Close закрывает соединение с базой данных
	Close() error
}
