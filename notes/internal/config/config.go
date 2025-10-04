package config

import (
	"fmt"
	"notes/internal/errors"
	"os"
	"strconv"
)

// Config - структура для хранения конфигурации приложения
// Содержит параметры, необходимые для запуска сервера
// Параметры могут быть загружены из файла, переменных окружения или других источников
type Config struct {
	Port          string // Порт, на котором будет запущен сервер
	Host          string // Хост, на котором будет запущен сервер
	DB_NAME       string // Имя базы данных, например "notes_db"
	DB_COLLECTION string // Коллекция в базе данных, например "notes"
	DBDSN         string // Строка подключения к базе данных, например, "postgres://user:password@localhost:5432/dbname"
	DBSSL         string // Настройка SSL для подключения к MongoDB
	JWTSecretKey  string // Секретный ключ для JWT токенов
	Timeout       int    // Таймаут для операций с сервером в секундах
	DBTimeout     int    // Таймаут для операций с базой данных в секундах
	RedisHost     string // Хост Redis сервера
	RedisPort     string // Порт Redis сервера
	RedisPassword string // Пароль для подключения к Redis
}

// NewConfig - конструктор для создания новой конфигурации
// Возвращает указатель на Config
func NewConfig() *Config {
	// Попытка получить порт из переменной окружения
	port, err := getEnv("PORT")
	if err != nil {
		fmt.Println("Не удалось получить PORT из переменной окружения")
	}
	// Попытка получить хост из переменной окружения
	host, err := getEnv("HOST")
	if err != nil {
		fmt.Println("Не удалось получить HOST из переменной окружения")
	}

	dbUsername, err := getEnv("MONGO_INITDB_ROOT_USERNAME")
	if err != nil {
		fmt.Println("Не удалось получить MONGO_INITDB_ROOT_USERNAME из переменной окружения")
	}
	dbPassword, err := getEnv("MONGO_INITDB_ROOT_PASSWORD")
	if err != nil {
		fmt.Println("Не удалось получить MONGO_INITDB_ROOT_PASSWORD из переменной окружения")
	}
	dbPort, err := getEnv("MONGO_INITDB_PORT")
	if err != nil {
		fmt.Println("Не удалось получить MONGO_INITDB_PORT из переменной окружения")
	}
	dbHost, err := getEnv("MONGO_INITDB_HOST")
	if err != nil {
		fmt.Println("Не удалось получить MONGO_INITDB_HOST из переменной окружения")
	}
	dbName, err := getEnv("MONGO_INITDB_DATABASE")
	if err != nil {
		fmt.Println("Не удалось получить MONGO_INITDB_DATABASE из переменной окружения")
	}

	// Формирование строки подключения к базе данных
	dbDSN := fmt.Sprintf(
		"mongodb://%s:%s@%s:%s/%s?authSource=admin",
		dbUsername,
		dbPassword,
		dbHost,
		dbPort,
		dbName,
	)
	// Параметры SSL для подключения к базе данных
	dbSSL, err := getEnv("MONGO_USE_SSL")
	if err != nil {
		fmt.Println("Не удалось получить MONGO_USE_SSL из переменной окружения")
	}
	// Если SSL не используется, добавляем параметр в строку подключения
	if dbSSL == "disable" {
		dbDSN += "&ssl=false"
	} else {
		dbDSN += "&ssl=true"
	}

	// JWT настройки
	jwtSecretKey, err := getEnv("JWT_SECRET_KEY")
	if err != nil {
		fmt.Println("Не удалось получить JWT_SECRET_KEY из переменной окружения")
	}

	// Попытка получить таймаут из переменной окружения
	timeout := 10
	if envValue, err := getEnv("SERVER_TIMEOUT"); err == nil {
		if parsed, parseErr := strconv.Atoi(envValue); parseErr == nil {
			timeout = parsed
		}
	} else {
		fmt.Println("Не удалось получить SERVER_TIMEOUT из переменной окружения, используется 10 секунд")
	}

	// Попытка получить таймаут для базы данных из переменной окружения
	dbTimeout := 5 // по умолчанию 5 секунд
	if envValue, err := getEnv("DB_TIMEOUT"); err == nil {
		if parsed, parseErr := strconv.Atoi(envValue); parseErr == nil {
			dbTimeout = parsed
		}
	} else {
		fmt.Println("Не удалось получить DB_TIMEOUT из переменной окружения, используется 5 секунд")
	}

	// Попытка получить хост Redis из переменной окружения
	redisHost, err := getEnv("REDIS_HOST")
	if err != nil {
		fmt.Println("Не удалось получить REDIS_HOST из переменной окружения")
	}

	// Попытка получить порт Redis из переменной окружения
	redisPort, err := getEnv("REDIS_PORT")
	if err != nil {
		fmt.Println("Не удалось получить REDIS_PORT из переменной окружения")
	}
	// Попытка получить пароль Redis из переменной окружения
	redisPassword, err := getEnv("REDIS_PASSWORD")
	if err != nil {
		fmt.Println("Не удалось получить REDIS_PASSWORD из переменной окружения")
	}

	dbCollection, err := getEnv("DB_COLLECTION")
	if err != nil {
		fmt.Println("Не удалось получить DB_COLLECTION из переменной окружения")
	}

	return &Config{
		Port:          port,
		Host:          host,
		DBDSN:         dbDSN, // Строка подключения к базе данных
		DBSSL:         dbSSL, // Параметры SSL для подключения к базе данных
		JWTSecretKey:  jwtSecretKey,
		Timeout:       timeout,       // Таймаут для операций с сервером
		DBTimeout:     dbTimeout,     // Таймаут для операций с базой данных
		RedisHost:     redisHost,     // Хост Redis сервера
		RedisPort:     redisPort,     // Порт Redis сервера
		RedisPassword: redisPassword, // Пароль для подключения к Redis
		DB_NAME:       dbName,        // Имя базы данных
		DB_COLLECTION: dbCollection,  // Коллекция в базе данных
	}
}

// getEnv получает значение переменной окружения
// Принимает ключ переменной в качестве аргумента
// Возвращает значение переменной или ошибку, если переменная не установлена
func getEnv(key string) (string, error) {
	value := os.Getenv(key)
	if value == "" {
		return "", fmt.Errorf("%w: %s", errors.ErrMissingEnvVar, key)
	}
	return value, nil
}
