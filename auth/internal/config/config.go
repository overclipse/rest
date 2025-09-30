package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config - структура для хранения конфигурации приложения
// Содержит параметры, необходимые для запуска сервера
// Параметры могут быть загружены из файла, переменных окружения или других источников
type Config struct {
	Port                   string // Порт, на котором будет запущен сервер
	Host                   string // Хост, на котором будет запущен сервер
	Timeout                int    // Таймаут для операций с сервером в секундах
	DBDSN                  string // Строка подключения к базе данных, например, "postgres://user:password@localhost:5432/dbname"
	DBTimeout              int    // Таймаут для операций с базой данных в секундах
	JWTSecretKey           string // Секретный ключ для JWT токенов
	AccessTokenExpiration  int    // Срок действия access токена в часах
	RefreshTokenExpiration int    // Срок действия refresh токена в часах
}

func getEnv(key string) (string, error) {
	value := os.Getenv(key)
	if value == "" {
		return "", fmt.Errorf("%s: %s", "не установленна переменая окружения", key)
	}
	return value, nil
}

// NewConfig - конструктор для создания новой конфигурации
// Возвращает указатель на Config с параметрами по умолчанию
func NewConfig() *Config {
	// Попытка получить порт из переменной окружения
	port, err := getEnv("PORT")
	if err != nil {
		fmt.Println("Не удалось получить PORT из переменной окружения, используется порт по умолчанию")
	}
	// Попытка получить хост из переменной окружения
	host, err := getEnv("HOST")
	if err != nil {
		fmt.Println("Не удалось получить HOST из переменной окружения, используется хост по умолчанию")
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
	dbHost, err := getEnv("POSTGRES_HOST")
	if err != nil {
		fmt.Println("Не удалось получить POSTGRES_HOST из переменной окружения")
	}
	dbPort, err := getEnv("POSTGRES_PORT")
	if err != nil {
		fmt.Println("Не удалось получить POSTGRES_PORT из переменной окружения")
	}
	dbUser, err := getEnv("POSTGRES_USER")
	if err != nil {
		fmt.Println("Не удалось получить POSTGRES_USER из переменной окружения")
	}
	dbPassword, err := getEnv("POSTGRES_PASSWORD")
	if err != nil {
		fmt.Println("Не удалось получить POSTGRES_PASSWORD из переменной окружения")
	}
	dbName, err := getEnv("POSTGRES_DB")
	if err != nil {
		fmt.Println("Не удалось получить POSTGRES_DB из переменной окружения")
	}
	dbSSL, err := getEnv("POSTGRES_USE_SSL")
	if err != nil {
		fmt.Println("Не удалось получить POSTGRES_USE_SSL из переменной окружения")
	}
	// Формирование строки подключения к базе данных
	dbDSN := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		dbUser, dbPassword, dbHost, dbPort, dbName, dbSSL)
	// JWT настройки
	jwtSecretKey, err := getEnv("JWT_SECRET_KEY")
	if err != nil {
		fmt.Println("Не удалось получить JWT_SECRET_KEY из переменной окружения, используется значение по умолчанию")
	}

	accessTokenExpiration := 24 // по умолчанию 24 часа
	if envValue, err := getEnv("JWT_ACCESS_TOKEN_EXPIRATION"); err == nil {
		if parsed, parseErr := strconv.Atoi(envValue); parseErr == nil {
			accessTokenExpiration = parsed
		}
	}

	refreshTokenExpiration := 168 // по умолчанию 7 дней (168 часов)
	if envValue, err := getEnv("JWT_REFRESH_TOKEN_EXPIRATION"); err == nil {
		if parsed, parseErr := strconv.Atoi(envValue); parseErr == nil {
			refreshTokenExpiration = parsed
		}
	}

	return &Config{
		Port:                   port,
		Host:                   host,
		DBDSN:                  dbDSN, // Строка подключения к базе данных
		JWTSecretKey:           jwtSecretKey,
		AccessTokenExpiration:  accessTokenExpiration,
		RefreshTokenExpiration: refreshTokenExpiration,
		Timeout:                timeout,   // Таймаут для операций с сервером
		DBTimeout:              dbTimeout, // Таймаут для операций с базой данных
	}
}
