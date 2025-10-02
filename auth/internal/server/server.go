package server

import (
	"auth/internal/config"
	"auth/internal/handler"
	"auth/internal/routes"
	"auth/internal/service"
	"auth/internal/errors"
	"fmt"

	"github.com/gin-gonic/gin"
)

type Server struct {
	//Конфиг подглючения
	cfg *config.Config
	//Роутер сервера
	router *gin.Engine
}

// NewServer - конструктор сервера
func NewServer(cfg *config.Config) (*Server, error) {
	// Проверяем, что конфигурация не пустая
	if cfg == nil {
		return nil, fmt.Errorf("конфигурация сервера не может быть nil")
	}
	service, err := service.NewService(cfg)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errors.ErrServiceCreation, err)
	}
	//Создаем новый экземпляр обработчика
	handler := handler.NewHandler(service, cfg)
	if handler == nil {
		return nil, fmt.Errorf("не удалось создать обработчик сервера")
	}
	router := routes.SetupRouter(handler)
	fmt.Println("Обработчик сервера успешно создан")
	// Создаем новый экземпляр сервера
	return &Server{
		router: router,
		cfg: cfg,
	}, nil
}

// Stop - остановка сервера
func (s *Server) Stop() error {
	fmt.Println("Сервер остановлен")
	return nil
}

// Serve - основной метод сервера
func (s *Server) Serve() error {
	// Запускаем сервер
	address := fmt.Sprintf("%s:%s", s.cfg.Host, s.cfg.Port)
	fmt.Printf("Сервер готов к обработке запросов на %s...\n", address)
	return s.router.Run(address)
}