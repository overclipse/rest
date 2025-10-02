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
	//time.Sleep(1 * time.Second) как было 
	db, err := openDBWithRetry(ctx, cfg.DBDSN)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к базе данных")
	}
	//тут мы проверяем реальную работу сокета как что он работает или нет

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

func openDBWithRetry(ctx context.Context, dsn string) (*gorm.DB, error) {
	backoff := 100 * time.Millisecond
	for {
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			return db, nil
		}
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("gorm.Open: не успели до дедлайна: %w", err)
		case <-time.After(backoff):
			if backoff < 2*time.Second {
				backoff *= 2
			}
		}
	}
}