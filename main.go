package main

import (
	"github.com/Noxdew/Knights-Of-Discord/bot"
	"github.com/Noxdew/Knights-Of-Discord/config"
	"github.com/Noxdew/Knights-Of-Discord/logger"
)

func main() {
	// Setup the logger
	logger.Init()

	// Read the config file
	config.Load()

	// Start the game
	bot.Start()
}
