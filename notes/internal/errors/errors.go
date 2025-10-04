package errors

import "errors"

// Общие ошибки сервиса заметок
var (
	// Ошибки заметок
	ErrNoteNotFound      = errors.New("заметка не найдена")
	ErrInvalidNoteID     = errors.New("некорректный ID заметки")
	ErrInvalidNoteData   = errors.New("неверные данные заметки")
	ErrNoteAlreadyExists = errors.New("заметка уже существует")
	ErrNoteCreation      = errors.New("ошибка создания заметки")
	ErrNoteUpdate        = errors.New("ошибка обновления заметки")
	ErrNoteDeletion      = errors.New("ошибка удаления заметки")

	// Ошибки авторизации
	ErrMissingAuthHeader = errors.New("отсутствует заголовок Authorization")
	ErrInvalidAuthFormat = errors.New("неверный формат токена")
	ErrTokenRequired     = errors.New("токен отсутствует или неверный формат")
	ErrInvalidToken      = errors.New("неверный или истекший токен")
	ErrRefreshToken      = errors.New("неверный или истекший refresh токен")
	ErrAuthRequired      = errors.New("необходима авторизация")
	ErrMissingUserID     = errors.New("ID пользователя не найден в токене")

	// Ошибки токенов
	ErrTokenGeneration = errors.New("ошибка генерации токенов")

	// Ошибки базы данных
	ErrDatabaseConnection = errors.New("ошибка подключения к базе данных")
	ErrDatabaseOperation  = errors.New("ошибка операции с базой данных")
	ErrDatabaseClose      = errors.New("ошибка закрытия соединения с базой данных")
	ErrDatabaseNotInit    = errors.New("база данных не инициализирована")
	ErrCacheConnection    = errors.New("ошибка подключения к кэшу")
	ErrCacheClose         = errors.New("ошибка закрытия соединения с кэшем")
	ErrCacheSet           = errors.New("ошибка записи в кэш")
	ErrCacheGet           = errors.New("ошибка чтения из кэша")
	ErrCacheSerialization = errors.New("ошибка сериализации данных для кэша")
	ErrIterationNotes     = errors.New("ошибка итерации по заметкам")
	ErrDecodeNote         = errors.New("ошибка декодирования заметки")

	// Ошибки конфигурации
	ErrMissingEnvVar = errors.New("переменная окружения не установлена")
	ErrEmptyDSN      = errors.New("строка подключения к базе данных не указана")

	// Ошибки сервиса
	ErrServiceCreation = errors.New("ошибка создания сервиса")
	ErrInvalidData     = errors.New("неверный формат данных")
)

// Сообщения для ошибок
const (
	// Сообщения для заметок
	MsgNoteNotFound      = "Заметка не найдена"
	MsgInvalidNoteID     = "Некорректный ID заметки"
	MsgInvalidNoteData   = "Неверные данные заметки"
	MsgNoteAlreadyExists = "Заметка уже существует"
	MsgNoteCreation      = "Ошибка создания заметки"
	MsgNoteUpdate        = "Ошибка обновления заметки"
	MsgNoteDeletion      = "Ошибка удаления заметки"

	// Сообщения для авторизации
	MsgMissingAuthHeader = "Отсутствует заголовок Authorization"
	MsgInvalidAuthFormat = "Неверный формат токена"
	MsgTokenRequired     = "Токен отсутствует или неверный формат"
	MsgInvalidToken      = "Неверный или истекший токен"
	MsgRefreshToken      = "Неверный или истекший refresh токен"
	MsgAuthRequired      = "Необходима авторизация"
	MsgMissingUserID     = "ID пользователя не найден в токене"

	// Сообщения для токенов
	MsgTokenGeneration = "Ошибка генерации токенов"

	// Сообщения для базы данных
	MsgDatabaseConnection = "Ошибка подключения к базе данных"
	MsgDatabaseOperation  = "Ошибка операции с базой данных"
	MsgDatabaseClose      = "Ошибка закрытия соединения с базой данных"
	MsgDatabaseNotInit    = "База данных не инициализирована"
	MsgIterationNotes     = "Ошибка итерации по заметкам"
	MsgDecodeNote         = "Ошибка декодирования заметки"

	// Сообщения для конфигурации
	MsgMissingEnvVar = "Переменная окружения не установлена"
	MsgEmptyDSN      = "Строка подключения к базе данных не указана"

	// Сообщения для сервиса
	MsgServiceCreation = "Ошибка создания сервиса"
	MsgInvalidData     = "Неверный формат данных"

	// Успешные сообщения
	MsgNoteCreated = "Заметка успешно создана"
	MsgNoteUpdated = "Заметка успешно обновлена"
	MsgNoteDeleted = "Заметка успешно удалена"
	MsgNoteFound   = "Заметка найдена"
	MsgNotesFound  = "Заметки получены"
)
