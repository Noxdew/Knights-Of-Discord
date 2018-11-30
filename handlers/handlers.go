package handlers

import (
	"strings"

	"github.com/Noxdew/Knights-Of-Discord/builder"
	"github.com/Noxdew/Knights-Of-Discord/command"
	"github.com/Noxdew/Knights-Of-Discord/config"
	"github.com/Noxdew/Knights-Of-Discord/db"
	"github.com/Noxdew/Knights-Of-Discord/logger"
	"github.com/Noxdew/Knights-Of-Discord/structure"

	"github.com/bwmarrin/discordgo"
)

// ReadyHandler is called when `Ready` event is triggered
func ReadyHandler(s *discordgo.Session, r *discordgo.Ready) {
	s.UpdateStatus(0, "Knights of Discord")
	structure.DefaultServer.BuildServer()
	logger.Log.Info("Knights of Discord has successfully started.")
}

// ServerJoinHandler is called when `GuildCreate` event is triggered
func ServerJoinHandler(s *discordgo.Session, g *discordgo.GuildCreate) {
	server, err := db.GetServer(g.Guild.ID)
	if err != nil && err != db.NotFound {
		logger.Log.Error(err.Error())
	} else if err == db.NotFound {
		builder.BuildServer(server, s, g.Guild)
	} else {
		// Server exists
	}
}

// MessageReceiveHandler function called when Message is sent
func MessageReceiveHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore bot messages
	if m.Author.Bot {
		return
	}

	// Check for command
	if strings.HasPrefix(m.Content, config.Get().Prefix) {
		// Get Server object
		c, err := s.Channel(m.ChannelID)
		if err != nil {
			logger.Log.Error(err.Error())
			return
		}
		server, err := db.GetServer(c.GuildID)
		if err != nil {
			logger.Log.Error(err.Error())
			return
		}

		// Execute command
		trigger := strings.Split(m.Content, config.Get().Prefix)[1]
		for _, cmd := range command.MessageCommands {
			if cmd.Trigger() == trigger {
				cmd.Execute(server, s, m)
				return
			}
		}

		// Fallback for wrong command
		// TODO
	}
}

// ReactionAddHandler function called when a Reaction is sent
func ReactionAddHandler(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	// Check author for bot
	u, err := s.User(r.UserID)
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}
	if u.Bot {
		return
	}

	// Get Server object
	c, err := s.Channel(r.ChannelID)
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}
	server, err := db.GetServer(c.GuildID)
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}

	// Call reaction command
	for _, cmd := range command.ReactionCommands {
		if cmd.Trigger() == r.Emoji.ID {
			cmd.Execute(server, s, r)
			return
		}
	}
}
