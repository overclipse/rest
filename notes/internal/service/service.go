package service

import (
	"context"
	"notes/internal/models"
)

// Service - интерфейс для работы с базой данных
// Он определяет методы, которые должны быть реализованы для работы с данными
type Service interface {
	Close() error                                                       // Закрывает соединение с базой данных
	Create(ctx context.Context, note models.Note) (*models.Note, error) // Создает новую заметку в базе данных
	GetByID(ctx context.Context, id string) (*models.Note, error)       // Получает заметку по идентификатору
	GetAll(ctx context.Context, authorId int) ([]models.Note, error)    // Получает все заметки из базы данных
	Update(ctx context.Context, note models.Note) (*models.Note, error) // Обновляет существующую заметку
	Delete(ctx context.Context, id string) error                        // Удаляет заметку по идентификатору
}