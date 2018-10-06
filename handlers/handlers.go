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
	} else if err == db.NotFound {
		builder.BuildServer(s, g.Guild)
	}
}

// ServerLeaveHandler is called when `GuildDelete` event is triggered
func ServerLeaveHandler(s *discordgo.Session, g *discordgo.GuildDelete) {
	err := db.RemoveServer(g.Guild.ID)
	if err != nil {
		logger.Log.Error(err.Error())
	}
}

// RoleEditHandler is called when `GuildRoleUpdate` event is triggered
func RoleEditHandler(s *discordgo.Session, r *discordgo.GuildRoleUpdate) {
	server, err := db.GetServer(r.GuildID)
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}
	if !server.Playing {
		return
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

// RoleDeleteHandler is called when `GuildRoleDelete` event is triggered
func RoleDeleteHandler(s *discordgo.Session, r *discordgo.GuildRoleDelete) {
	server, err := db.GetServer(r.GuildID)
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}
	if !server.Playing {
		return
	}
	for _, role := range server.Roles {
		g, _ := s.Guild(r.GuildID)
		if role.ID == r.RoleID {
			builder.BuildRole(s, g, role.Type)
		}
	}
}

// ChannelEditHandler is called when `ChannelUpdate` event is triggered
// TODO

// ChannelDeleteHandler is called when `ChannelDelete` event is triggered
func ChannelDeleteHandler(s *discordgo.Session, c *discordgo.ChannelDelete) {
	logger.Log.Debug("BOOP")
	server, err := db.GetServer(c.Channel.GuildID)
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}
	if !server.Playing {
		return
	}
	g, _ := s.Guild(c.Channel.GuildID)
	for _, channel := range server.Channels {
		if channel.ID == c.Channel.ID {
			for _, ch := range config.Get().Channels {
				if ch.Name == channel.Name {
					builder.BuildChannel(s, g, ch)
					break
				}
			}
		}
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
