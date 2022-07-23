package main

import (
	"fmt"
	"gitlab.ozon.dev/Bdido86/movie-tickets/internal/config"
	"gitlab.ozon.dev/Bdido86/movie-tickets/internal/pkg/bot"
	"gitlab.ozon.dev/Bdido86/movie-tickets/internal/pkg/bot/command"
)

func main() {
	c := config.GetConfig()
	if len(c.Token()) == 0 {
		fmt.Println("Config error: token is empty")
		return
	}

	cmd, err := bot.Init(c)
	if err != nil {
		fmt.Println("Error: " + err.Error())
		return
	}

	command.AddHandler(cmd)

	if err := cmd.Run(); err != nil {
		fmt.Println("Error: " + err.Error())
		return
	}
}
