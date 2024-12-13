package routes

import (
	"github.com/gin-gonic/gin"
	"slim-connector-back/internal/handler"
)

func UserRoutes(router *gin.RouterGroup, userHandler *handler.UserHandler) {
	router.POST("/users", userHandler.CreateUser) // Register
	router.GET("/users", userHandler.GetUsers)
}
