package database

import (
	"auth/internal/config"
	"context"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewDatabase - функция для создания нового подключения к базе данных
// Принимает контекст, DSN (строку подключения)
// Возвращает указатель на gorm.DB или ошибку, если она произошла
func NewDatabase(cfg *config.Config, models ...any) (*gorm.DB, error) {

	// Создаем контекст с таймаутом для инициализации БД
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.Timeout)*time.Second)
	// Отмена контекста при завершении работы функции
	defer cancel()

	// Добавляем задержку в 1 секунду, // чтобы дать время на инициализацию других компонентов, если это необходимо
	// Это может быть полезно, если база данных запускается в контейнере или сервисе, который требует времени на инициализацию
	// Например, если база данных запускается в Docker-контейнере, то может потребоваться время на его запуск и готовность к соединению
	time.Sleep(1 * time.Second)

	// Применяем контекст к подключению к БД
	db, err := gorm.Open(postgres.Open(cfg.DBDSN), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("%s: %v", "ошибка подключения к базе данных", err)
	}

	errMigration := runMigrations(db, models...) // Выполняем миграции
	if errMigration != nil {
		return nil, fmt.Errorf("%s: %v", "ошибка миграции базы данных", errMigration)
	}

	return db.WithContext(ctx), nil
}

// Автоматические миграции для моделей
func runMigrations(db *gorm.DB, models ...any) error {
	// Выполняем миграции для всех переданных моделей
	for _, model := range models {
		if err := db.AutoMigrate(model); err != nil {
			return fmt.Errorf("ошибка миграции модели %T: %w", model, err)
		}
	}
	// Если все миграции прошли успешно, возвращаем nil
	fmt.Println("Все миграции успешно выполнены")
	return nil
}