package handlers

import (
	"strings"

	"github.com/Noxdew/Knights-Of-Discord/builder"
	"github.com/Noxdew/Knights-Of-Discord/config"
	"github.com/Noxdew/Knights-Of-Discord/db"
	"github.com/Noxdew/Knights-Of-Discord/logger"
	"github.com/bwmarrin/discordgo"
)

// ReadyHandler is called when `Ready` event is triggered
func ReadyHandler(s *discordgo.Session, r *discordgo.Ready) {
	s.UpdateStatus(0, "Knights of Discord")
	logger.Log.Info("Knights of Discord has successfully started.")
}

// ServerJoinHandler is called when `GuildCreate` event is triggered
func ServerJoinHandler(s *discordgo.Session, g *discordgo.GuildCreate) {
	server, err := db.GetServer(g.Guild)
	if err != nil && err != db.NotFound {
		logger.Log.Error(err.Error())
	} else if err == db.NotFound {
		builder.BuildServer(server, s, g.Guild)
	} else {
		// Server exists
	}
}

// MessageReceiveHandler function called when message is sent
func MessageReceiveHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
		return
	}
	c, err := s.Channel(m.ChannelID)
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}
	g, err := s.Guild(c.GuildID)
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}
	server, err := db.GetServer(g)
	if err != nil && err != db.NotFound {
		logger.Log.Error(err.Error())
		return
	}
	if strings.HasPrefix(m.Content, config.Get().Prefix) {
		command := strings.Split(m.Content, config.Get().Prefix)[1]
		if command == "closeGame" && m.Author.ID == g.OwnerID {
			s.ChannelMessageSend(m.ChannelID, "Closing Game...")
			builder.DestroyServer(server, s, g)
		}
	}
}
