package bot

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/Noxdew/Knights-Of-Discord/config"
	"github.com/Noxdew/Knights-Of-Discord/handlers"
	"github.com/Noxdew/Knights-Of-Discord/logger"
	"github.com/bwmarrin/discordgo"
)

var s *discordgo.Session

// Start start the bot, connects it to the Discord servers and starts the game
func Start() {
	// Login the bot client and save its user ID
	s, err := discordgo.New("Bot " + config.Get().Token)
	if err != nil {
		logger.Log.Panic(err)
	}

	// Add event handlers
	s.AddHandler(handlers.ReadyHandler)
	s.AddHandler(handlers.ServerJoinHandler)
	s.AddHandler(handlers.ServerLeaveHandler)
	s.AddHandler(handlers.RoleEditHandler)
	s.AddHandler(handlers.RoleDeleteHandler)
	s.AddHandler(handlers.ChannelEditHandler)
	s.AddHandler(handlers.ChannelDeleteHandler)
	s.AddHandler(handlers.MessageReceiveHandler)

	// Start the bot's session
	err = s.Open()
	if err != nil {
		logger.Log.Panic(err)
	}

	// Wait here until CTRL-C or other term signal is received.
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}
