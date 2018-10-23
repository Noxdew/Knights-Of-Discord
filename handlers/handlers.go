package handlers

import (
	"strings"

	"github.com/Noxdew/Knights-Of-Discord/builder"
	"github.com/Noxdew/Knights-Of-Discord/config"
	"github.com/Noxdew/Knights-Of-Discord/db"
	"github.com/Noxdew/Knights-Of-Discord/logger"
	"github.com/Noxdew/Knights-Of-Discord/utils"
	"github.com/bwmarrin/discordgo"
)

// ReadyHandler is called when `Ready` event is triggered
func ReadyHandler(s *discordgo.Session, r *discordgo.Ready) {
	logger.Log.Info("Knights of Discord has successfully started.")
}

// ServerJoinHandler is called when `GuildCreate` event is triggered
func ServerJoinHandler(s *discordgo.Session, g *discordgo.GuildCreate) {
	_, err := db.GetServer(g.Guild.ID)
	if err != nil && err != db.NotFound {
		logger.Log.Error(err.Error())
	} else if err == db.NotFound {
		builder.BuildServer(s, g.Guild)
	} else {
		db.UpdateServerChecked(g.Guild.ID, false)
		db.UpdateServerPlaying(g.Guild.ID, true)
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
	g, err := s.Guild(r.GuildID)
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}
	_, err = utils.GetServerRoleByID(r.Role.ID, server)
	if err != nil {
		return
	}
	if !utils.CheckRole(r.Role) {
		builder.FixRole(s, g, r.Role)
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
	g, err := s.Guild(r.GuildID)
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}
	role, err := utils.GetServerRoleByID(r.RoleID, server)
	if err != nil {
		return
	}
	builder.BuildRole(s, g, role.DefName)
}

// ChannelEditHandler is called when `ChannelUpdate` event is triggered
func ChannelEditHandler(s *discordgo.Session, c *discordgo.ChannelUpdate) {
	server, err := db.GetServer(c.GuildID)
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}
	if !server.Playing {
		return
	}
	g, err := s.Guild(c.GuildID)
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}
	channel, err := utils.GetServerChannelByID(c.Channel.ID, server)
	if err != nil {
		return
	}
	if !utils.CheckChannel(c.Channel, g) {
		builder.FixChannel(s, g, c.Channel)
	}
	for _, perm := range c.PermissionOverwrites {
		needed := utils.CheckPermNeeded(perm, channel)
		if !needed {
			builder.RemovePerm(s, c.Channel.ID, perm)
		} else {
			if !utils.CheckPerm(perm, channel) {
				p, err := utils.GetChannelPermByID(perm.ID, channel)
				if err != nil {
					logger.Log.Error(err.Error())
					continue
				}
				builder.FixPerm(s, c.Channel.ID, p)
			}
		}
	}
	for _, perm := range channel.Perms {
		if !utils.CheckPermExists(perm, c.Channel) {
			s.ChannelPermissionSet(c.Channel.ID, perm.ID, perm.Type, perm.Allow, perm.Deny)
		}
	}
}

// ChannelDeleteHandler is called when `ChannelDelete` event is triggered
func ChannelDeleteHandler(s *discordgo.Session, c *discordgo.ChannelDelete) {
	server, err := db.GetServer(c.Channel.GuildID)
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}
	if !server.Playing {
		return
	}
	g, err := s.Guild(c.Channel.GuildID)
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}
	if c.Channel.ID == server.Category {
		builder.BuildCategory(s, g)
		return
	}
	channel, err := utils.GetServerChannelByID(c.Channel.ID, server)
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}
	ch, err := utils.GetConfigChannelByName(channel.DefName)
	builder.BuildChannel(s, g, ch)
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
