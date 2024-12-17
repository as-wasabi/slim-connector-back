package server

import (
	"github.com/gin-gonic/gin"
)

type Server struct {
	router *gin.Engine
}

func TasQServer(engine *gin.Engine) *Server {
	return &Server{
		router: engine,
	}
}

func (s *Server) Run(addr string) error {
	return s.router.Run(addr)
}
