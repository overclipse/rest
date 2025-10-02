package service

import (
	"auth/internal/models"
	"context"

	"gorm.io/gorm"
)

// DBService - реализация интерфейса Service для работы с базой данных
// Этот сервис взаимодействует с базой данных
// Сервис реализует методы интерфейса Service
type DBService struct {
	db *gorm.DB // Указатель на БД
}

// Проверка, что DBService реализует интерфейс Service
// Это позволяет гарантировать, что DBService соответствует контракту интерфейса Service
var _ Service = (*DBService)(nil)

// NewService - конструктор для создания нового экземпляра Service
// Он принимает конфигурацию и возвращает указатель на текущею реализацию
// сервиса или ошибку, если она произошла

func NewService(db *gorm.DB) Service {
	return &DBService{db: db}
}

// Create создает нового пользователя в базе данных
// Он принимает контекст и указатель на модель User
// Возвращает созданного пользователя или ошибку, если она произошла
func (p *DBService) Create(ctx context.Context, user *models.User) (*models.User, error) {
	if user == nil {
		return nil, gorm.ErrInvalidData
	}

	// Хешируем пароль перед сохранением в базу данных
	hashedPassword, err := user.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}
	user.Password = hashedPassword

	// Используем контекст для управления временем выполнения операции
	if err := p.db.WithContext(ctx).Create(user).Error; err != nil {
		return nil, err
	}

	// Возвращаем созданного пользователя
	return user, nil
}

// Delete Удаляет пользователя из базы данных по ID
func (p *DBService) Delete(ctx context.Context, id int) error {
	if id <= 0 {
		return gorm.ErrRecordNotFound
	}

	// Используем контекст для управления временем выполнения операции
	if err := p.db.WithContext(ctx).Delete(&models.User{ID: id}).Error; err != nil {
		return err
	}

	return nil
}

// Read - находит пользователя по ID в базе данных
func (p *DBService) Read(ctx context.Context, id int) (*models.User, error) {
	if id <= 0 {
		return nil, gorm.ErrRecordNotFound
	}

	var user models.User
	// Используем контекст для управления временем выполнения операции
	if err := p.db.WithContext(ctx).First(&user, id).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// Update - обновляет данные пользователя в базе данных
func (p *DBService) Update(ctx context.Context, user *models.User) error {
	if user == nil || user.ID <= 0 {
		return gorm.ErrInvalidData
	}

	// Если пароль не пустой, значит он был изменен и нужно его хешировать
	if user.Password != "" {
		hashedPassword, err := user.HashPassword(user.Password)
		if err != nil {
			return err
		}
		user.Password = hashedPassword
	}

	// Используем контекст для управления временем выполнения операции
	if err := p.db.WithContext(ctx).Save(user).Error; err != nil {
		return err
	}

	return nil
}

// ReadByUsername - находит пользователя по имени пользователя в базе данных
func (p *DBService) ReadByUsername(ctx context.Context, username string) (*models.User, error) {
	if username == "" {
		return nil, gorm.ErrInvalidData
	}

	var user models.User
	// Используем контекст для управления временем выполнения операции
	if err := p.db.WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// Authenticate - аутентифицирует пользователя по имени пользователя и паролю
func (p *DBService) Authenticate(ctx context.Context, username, password string) (*models.User, error) {
	if username == "" || password == "" {
		return nil, gorm.ErrInvalidData
	}

	var user models.User
	// Ищем пользователя по имени пользователя
	if err := p.db.WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}

	// Проверяем соответствие пароля хешу
	if !user.CheckPassword(password, user.Password) {
		return nil, gorm.ErrRecordNotFound // Возвращаем ошибку "не найден" для безопасности
	}

	return &user, nil
}

// Close закрывает соединение с базой данных
func (p *DBService) Close() error {
	db, err := p.db.DB()
	if err != nil {
		return err
	}
	return db.Close()
}
