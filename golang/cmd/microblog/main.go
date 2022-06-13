package main

import (
	"log"
	"time"

	"github.com/chipocrudos/microblog/config"
	"github.com/chipocrudos/microblog/internal/server"
	"github.com/chipocrudos/microblog/internal/server/data"
)

func main() {

	if config.Config.Debug {
		log.Println("Debug mode")
	}

	for {
		log.Println("Waiting to database")
		if data.MongoCN.PingDB() {
			log.Println("Database connection ready")
			break
		}
		time.Sleep(2 * time.Second)

	}

	serv, err := server.New(config.Config.PORT)
	if err != nil {
		log.Fatal(err)
	}

	serv.Start()
}
