package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"slim-connector-back/config"
	"slim-connector-back/internal/repository"
)

func main() {
	engine := gin.Default()
	engine.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "hello world",
		})

		cfg, err := config.LoadConfig("")

		if err != nil {
			log.Fatalf("Cant Load config: %v", err)
		}

		mongoDB, err := repository.FetchMongoDB(cfg)
		if err != nil {
			return
		}

		defer func(mongoDB *repository.MongoDB) {
			err := mongoDB.Close()
			if err != nil {

			}
		}(mongoDB)
	})

	err := engine.Run(":3000")

	if err != nil {
		println(err)
		return
	}
}
