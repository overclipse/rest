package main

import (
	"auth/internal/config"
	"auth/internal/server"
	"fmt"
)

func main() {
	// Инициализируем конфигурации сервера
	cfg := config.NewConfig()

	fmt.Printf("=== Server Configuration ===\n")
	fmt.Printf("Host: %s\n", cfg.Host)
	fmt.Printf("Port: %s\n", cfg.Port)
	fmt.Printf("Database URL: %s\n", cfg.DBDSN)
	fmt.Printf("Access Token Expiration: %d hours\n", cfg.AccessTokenExpiration)
	fmt.Printf("Refresh Token Expiration: %d hours\n", cfg.RefreshTokenExpiration)
	fmt.Printf("Server Timeout: %d seconds\n", cfg.Timeout)
	fmt.Printf("Database Timeout: %d seconds\n", cfg.DBTimeout)
	fmt.Printf("JWT Secret Key: %s\n", cfg.JWTSecretKey)
	fmt.Printf("=============================\n")
	
	server, err := server.NewServer(cfg)
	if err != nil {
		fmt.Printf("Ошибка при создании сервера %v\n", err)
		return
	}
	fmt.Printf("Сервер успешно создан\n")
	if err := server.Serve(); err != nil {
		fmt.Printf("Ошибка запуска сервера %v\n", err)
		return
	}
	fmt.Printf("Сервер запущен успешно\n")
}