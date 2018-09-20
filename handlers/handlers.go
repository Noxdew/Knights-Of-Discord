package handlers

import (
	"strings"

	"github.com/Noxdew/Knights-Of-Discord/builder"
	"github.com/Noxdew/Knights-Of-Discord/config"
	"github.com/Noxdew/Knights-Of-Discord/db"
	"github.com/Noxdew/Knights-Of-Discord/logger"

	"github.com/bwmarrin/discordgo"
)

var roles = []string{
	"KoD-King",
	"KoD-Knight",
	"KoD-Esquire",
	"KoD-Villager",
}

// ServerJoinHandler function called when a server is joined
func ServerJoinHandler(s *discordgo.Session, g *discordgo.GuildCreate) {
	// Check for existing game
	_, err := db.GetServer(g.Guild.ID)
	if err != nil && err != db.NotFound {
		logger.Log.Error(err.Error())
	}
	if err == db.NotFound {
		builder.BuildServer(g.Guild)
		builder.BuildRoles(s, g.Guild, roles)
	} else {
		builder.CheckRoles(s, g.Guild)
		logger.Log.Info("Game already running on %s (%s)", g.Guild.Name, g.Guild.ID)
	}
}

// MessageReceiveHandler function called when message is sent
func MessageReceiveHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	c, _ := s.Channel(m.ChannelID)
	g, _ := s.Guild(c.GuildID)
	if m.Author.Bot {
		return
	}
	if strings.HasPrefix(m.Content, config.Get().Prefix) {
		command := strings.Split(m.Content, config.Get().Prefix)[1]
		if command == "closeGame" && m.Author.ID == g.OwnerID {
			builder.DestroyRoles(s, g)
			builder.DestroyServer(g)
			s.ChannelMessageSend(m.ChannelID, "Game has been closed.")
			err := s.GuildLeave(g.ID)
			if err != nil {
				logger.Log.Error(err.Error())
			}
		}
	}
}
