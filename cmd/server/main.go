package main

import (
	"gitlab.ozon.dev/Bdido86/movie-tickets/internal/config"
	"log"
)

func main() {
	c := config.GetConfig()
	if len(c.Port()) == 0 {
		log.Fatal("Config error: port is empty")
	}

	address := ":" + c.Port()
	runGRPCServer(address)

}
