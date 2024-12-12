package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"slim-connector-back/config"
)

func main() {
	engine := gin.Default()
	engine.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "hello world",
		})

		config.FetchDB()
	})

	err := engine.Run(":3000")

	if err != nil {
		println(err)
		return
	}
}
