package main

import (
	"fmt"
	"ozon/go-hw-bot/config"
	"ozon/go-hw-bot/internal/commander"
	"ozon/go-hw-bot/internal/handlers"
)

func main() {
	fmt.Println("Start bot")

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
