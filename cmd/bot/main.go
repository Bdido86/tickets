package main

import (
	"fmt"
	"gitlab.ozon.dev/Bdido86/movie-tickets/config"
	"gitlab.ozon.dev/Bdido86/movie-tickets/internal/commander"
	"gitlab.ozon.dev/Bdido86/movie-tickets/internal/handlers"
)

func main() {
	c := config.GetConfig()
	if len(c.Token()) == 0 {
		fmt.Println("Config error: token is empty")
		return
	}

	cmd, err := commander.Init(c)
	if err != nil {
		fmt.Println("Error: " + err.Error())
		return
	}

	handlers.AddHandler(cmd)

	if err := cmd.Run(); err != nil {
		fmt.Println("Error: " + err.Error())
		return
	}
}
