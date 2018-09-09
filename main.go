package main

import (
	"fmt"

	"github.com/noxdew/knights-of-discord/bot"
	"github.com/noxdew/knights-of-discord/config"
)

func main() {
	// Read the config file, exits on error
	err := config.ReadConfig()
	if err != nil {
		fmt.Println(err.Error)
		return
	}

	// Call the main bot
	bot.Start()
}
