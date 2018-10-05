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
	logger.Log.Info("Knights of Discord has successfully started.")
}

// ServerJoinHandler is called when `GuildCreate` event is triggered
func ServerJoinHandler(s *discordgo.Session, g *discordgo.GuildCreate) {
	// Check for existing game
	_, err := db.GetServer(g.Guild.ID)
	if err != nil && err != db.NotFound {
		logger.Log.Error(err.Error())
	}
	if err == db.NotFound {
		builder.BuildServer(s, g.Guild)
	}
}

// RoleEditHandler function called when `GuildRoleUpdate` is triggered
func RoleEditHandler(s *discordgo.Session, r *discordgo.GuildRoleUpdate) {
	server, err := db.GetServer(r.GuildID)
	if err != nil {
		logger.Log.Error(err.Error())
	}
	for _, role := range server.Roles {
		g, _ := s.Guild(r.GuildID)
		if role.ID == r.Role.ID {
			if !r.Role.Mentionable || r.Role.Permissions != config.Get().RolePerm {
				builder.FixRole(s, g, r.Role)
			}
		}
	}
}

// RoleDeleteHandler function called when `GuildRoleDelete` is triggered
func RoleDeleteHandler(s *discordgo.Session, r *discordgo.GuildRoleDelete) {
	server, err := db.GetServer(r.GuildID)
	if err != nil {
		logger.Log.Error(err.Error())
	}
	for _, role := range server.Roles {
		g, _ := s.Guild(r.GuildID)
		if role.ID == r.RoleID {
			builder.BuildRole(s, g, role.Type)
		}
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
			builder.DestroyServer(s, g)
			s.ChannelMessageSend(m.ChannelID, "Game has been closed.")
			err := s.GuildLeave(g.ID)
			if err != nil {
				logger.Log.Error(err.Error())
			}
		}
	}
}
