package database

import (
	"context"
	"fmt"
	"notes/internal/config"
	"notes/internal/errors"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// NewDatabase - функция для создания нового подключения к базе данных
// Принимает контекст, DSN (строку подключения)
// Возвращает указатель на gorm.DB или ошибку, если она произошла
func NewDatabase(cfg *config.Config) (*mongo.Client, error) {
	// Проверяем, что DSN не пустой
	if cfg.DBDSN == "" {
		return nil, fmt.Errorf("%w", errors.ErrEmptyDSN)
	}

	// Создаем контекст с таймаутом для инициализации БД
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.Timeout)*time.Second)
	// Отмена контекста при завершении работы функции
	defer cancel()

	// Применяем контекст к подключению к БД
	db, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.DBDSN))
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errors.ErrDatabaseConnection, err)
	}

	return db, nil
}

// CloseDB также должен работать с переданным db
func CloseDB(db *mongo.Client, cgf *config.Config) error {
	// Создаем новый контекст с таймаутом и предусматриваем его корректное завершение.
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cgf.Timeout)*time.Second)

	// Отмена контекста при завершении работы функции
	defer cancel()
	if db == nil {
		return fmt.Errorf("%w", errors.ErrDatabaseNotInit)
	}

	return db.Disconnect(ctx)
}
