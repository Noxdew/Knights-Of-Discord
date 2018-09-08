package main

import (
	"fmt"
	"github.com/noxdew/knights-of-discord/config"
	"github.com/noxdew/knights-of-discord/bot"
)

func main() {
	err := config.ReadConfig()

	if err != nil {
		fmt.Println(err.Error)
		return
	}

	bot.Start()

	<-make(chan struct{})
	return
}