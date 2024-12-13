package routes

import (
	"github.com/gin-gonic/gin"
	"slim-connector-back/internal/handler"
)

func InitRoutes(router *gin.Engine, userHandler *handler.UserHandler) {
	api := router.Group("/tasq")
	UserRoutes(api, userHandler)

}
