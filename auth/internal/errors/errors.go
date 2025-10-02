package errors

import "errors"

// Общие ошибки сервиса аутентификации
var (
	// Ошибки пользователей
	ErrUserNotFound       = errors.New("пользователь не найден")
	ErrInvalidCredentials = errors.New("неверные учетные данные")
	ErrUserAlreadyExists  = errors.New("пользователь уже существует")
	ErrInvalidUserData    = errors.New("неверные данные пользователя")
	ErrUserIdNotFound     = errors.New("ID пользователя не найден в контексте")

	// Ошибки авторизации
	ErrMissingAuthHeader = errors.New("отсутствует заголовок Authorization")
	ErrInvalidAuthFormat = errors.New("неверный формат токена")
	ErrTokenRequired     = errors.New("токен отсутствует или неверный формат")
	ErrInvalidToken      = errors.New("неверный или истекший токен")
	ErrRefreshToken      = errors.New("неверный или истекший refresh токен")
	ErrAuthRequired      = errors.New("необходима авторизация")

	// Ошибки токенов
	ErrTokenGeneration = errors.New("ошибка генерации токенов")

	// Ошибки базы данных
	ErrDatabaseConnection = errors.New("ошибка подключения к базе данных")
	ErrDatabaseMigration  = errors.New("ошибка выполнения миграций")
	ErrDatabaseClose      = errors.New("ошибка закрытия соединения с базой данных")
	ErrDatabaseOperation  = errors.New("ошибка операции с базой данных")

	// Ошибки конфигурации
	ErrMissingEnvVar   = errors.New("переменная окружения не установлена")
	ErrEmptyDSN        = errors.New("строка подключения к базе данных не указана")
	ErrDatabaseNotInit = errors.New("база данных не инициализирована")

	// Ошибки сервиса
	ErrServiceCreation = errors.New("ошибка создания сервиса")
	ErrInvalidData     = errors.New("неверный формат данных")
	ErrUserCreation    = errors.New("ошибка создания пользователя")
)

// Сообщения для ошибок
const (
	// Сообщения для пользователей
	MsgUserNotFound       = "Пользователь не найден"
	MsgInvalidCredentials = "Неверные учетные данные"
	MsgUserAlreadyExists  = "Пользователь уже существует"
	MsgInvalidUserData    = "Неверные данные пользователя"
	MsgUserIdNotFound     = "ID пользователя не найден в контексте"

	// Сообщения для авторизации
	MsgMissingAuthHeader = "Отсутствует заголовок Authorization"
	MsgInvalidAuthFormat = "Неверный формат токена"
	MsgTokenRequired     = "Токен отсутствует или неверный формат"
	MsgInvalidToken      = "Неверный или истекший токен"
	MsgRefreshToken      = "Неверный или истекший refresh токен"
	MsgAuthRequired      = "Необходима авторизация"

	// Сообщения для токенов
	MsgTokenGeneration = "Ошибка генерации токенов"

	// Сообщения для базы данных
	MsgDatabaseConnection = "Ошибка подключения к базе данных"
	MsgDatabaseMigration  = "Ошибка выполнения миграций"
	MsgDatabaseClose      = "Ошибка закрытия соединения с базой данных"
	MsgDatabaseOperation  = "Ошибка операции с базой данных"

	// Сообщения для конфигурации
	MsgMissingEnvVar   = "Переменная окружения не установлена"
	MsgEmptyDSN        = "Строка подключения к базе данных не указана"
	MsgDatabaseNotInit = "База данных не инициализирована"

	// Сообщения для сервиса
	MsgServiceCreation = "Ошибка создания сервиса"
	MsgInvalidData     = "Неверный формат данных"
	MsgUserCreation    = "Ошибка создания пользователя"

	// Успешные сообщения
	MsgUserRegistered  = "Пользователь успешно зарегистрирован"
	MsgLoginSuccess    = "Успешная авторизация"
	MsgTokensRefreshed = "Токены успешно обновлены"
	MsgUserUpdated     = "Данные пользователя успешно обновлены"
	MsgUserDeleted     = "Пользователь успешно удален"
)
