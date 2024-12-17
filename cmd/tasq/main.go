package main

import (
	"log"
	"slim-connector-back/internal"
	"slim-connector-back/internal/handler"
)

func main() {
	initializer, err := internal.NewInitializer()
	if err != nil {
		log.Fatal(err)
		return
	}
	initializer.InitRoute(
		handler.InitRoute(initializer)...,
	)
	if err := initializer.Server.Run(":3000"); err != nil {
		log.Fatal(err)
	}
}
