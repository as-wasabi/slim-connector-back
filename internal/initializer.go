package internal

import (
	"github.com/gin-gonic/gin"
	"slim-connector-back/config"
	"slim-connector-back/internal/repository"
	"slim-connector-back/internal/server"
)

type Initializer struct {
	*repository.MongoDB
	*server.Server
	*gin.Engine
}

type Route interface {
	InitRoute(group *gin.RouterGroup)
}

func NewInitializer() (*Initializer, error) {
	cfg, err := config.LoadConfig("")
	if err != nil {
		return nil, err
	}
	mongoDB, err := initMongo(cfg)
	if err != nil {
		return nil, err
	}
	engine := gin.Default()
	srv := server.TasQServer(engine)
	return &Initializer{
		mongoDB,
		srv,
		engine,
	}, nil
}
func (receiver Initializer) InitRoute(routes ...Route) {
	api := receiver.Engine.Group("/tasq")

	for _, route := range routes {
		route.InitRoute(api)
	}
}

func (receiver Initializer) InitAIRoute(routes ...Route) {
	api := receiver.Engine.Group("/openai")

	for _, route := range routes {
		route.InitRoute(api)
	}
}

func initMongo(cfg *config.Config) (*repository.MongoDB, error) {
	mongoDB, err := repository.FetchMongoDB(cfg)
	if err != nil {
		return nil, err
	}
	return mongoDB, nil
}
