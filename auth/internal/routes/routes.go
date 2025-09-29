package routes

import (
	"auth/internal/handler"

	"github.com/gin-gonic/gin"
)

func SetupRouter(h *handler.Handler) *gin.Engine {
	router := gin.Default()
	
	//Группа маршрутов для аутентификации
	auth := router.Group("/auth")
	{
		//Публичные endpoints
		auth.POST("/register", h.RegisterUser)
		auth.POST("/login", h.LoginUser)
		//Защитные endpoints
		auth.GET("/user", h.GetUserInfo)
		auth.PUT("/user", h.UpdateUser)
		auth.DELETE("/user", h.DeleteUser)
	}
	
	return router
}