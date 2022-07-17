package main

import (
	"fmt"
	"ozon/go-hw-bot/config"
	"ozon/go-hw-bot/internal/commander"
)

func main() {
	fmt.Println("Start bot")

	c := config.New()
	if len(c.Token()) == 0 {
		fmt.Println("Config error: token is empty")
		return
	}

	cmd, err := commander.Init(c)
	if err != nil {
		fmt.Println("Error: " + err.Error())
		return
	}

	if err := cmd.Run(); err != nil {
		fmt.Println("Error: " + err.Error())
		return
	}
}
