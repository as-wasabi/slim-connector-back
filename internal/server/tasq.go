package server

import (
	"github.com/gin-gonic/gin"
	"slim-connector-back/internal/handler"
	"slim-connector-back/internal/repository"
	"slim-connector-back/internal/routes"
)

type Server struct {
	router *gin.Engine
}

func TasQServer(mongoDB *repository.MongoDB) *Server {
	router := gin.Default()
	userHandler := handler.NewUserHandler(mongoDB)

	routes.InitRoutes(router, userHandler)

	return &Server{
		router: router,
	}
}

func (s *Server) Run(addr string) error {
	return s.router.Run(addr)
}
