package main

import (
	"log"
	"slim-connector-back/config"
	"slim-connector-back/internal/repository"
	"slim-connector-back/internal/server"
)

func main() {
	cfg, err := config.LoadConfig("")

	if err != nil {
		log.Fatal(err)
	}

	mongoDB, err := repository.FetchMongoDB(cfg)
	if err != nil {
		log.Fatal(err)
		return
	}

	srv := server.TasQServer(mongoDB)
	if err := srv.Run(":3000"); err != nil {
		log.Fatal(err)
	}
}
