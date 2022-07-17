package main

import (
	"fmt"
	"ozon/go-hw-bot/config"
)

func main() {
	fmt.Println("Start bot")

	c := config.New()
	if len(c.Token()) == 0 {
		fmt.Println("Config error: token is empty")
		return
	}

	fmt.Print(c)
}
