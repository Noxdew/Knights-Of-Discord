package utils

import (
	"github.com/Noxdew/Knights-Of-Discord/config"
	"github.com/Noxdew/Knights-Of-Discord/db"
	"github.com/Noxdew/Knights-Of-Discord/logger"
	"github.com/bwmarrin/discordgo"
)

// GetServerRoleByName returns Role with default name `r` from Server `s`
func GetServerRoleByName(r string, s *db.Server) (db.Role, error) {
	for _, role := range s.Roles {
		if role.DefName == r {
			return role, nil
		}
	}
	return db.Role{}, db.NotFound
}

// GetServerRoleByID returns Role with ID `r` from Server `s`
func GetServerRoleByID(r string, s *db.Server) (db.Role, error) {
	if len(s.Roles) < 1 {
		return db.Role{}, db.NotFound
	}
	for _, role := range s.Roles {
		if role.ID == r {
			return role, nil
		}
	}
	return db.Role{}, db.NotFound
}

// CheckRole evaluates Role `r` against desired game role options
func CheckRole(r *discordgo.Role) bool {
	if !r.Mentionable || r.Permissions != config.Get().RolePerm {
		return false
	}
	return true
}

// GetServerChannelByName returns Channel with default name `c` from Server `s`
func GetServerChannelByName(c string, s *db.Server) (db.Channel, error) {
	if len(s.Channels) < 1 {
		return db.Channel{}, db.NotFound
	}
	for _, channel := range s.Channels {
		if channel.DefName == c {
			return channel, nil
		}
	}
	return db.Channel{}, db.NotFound
}

// GetServerChannelByID returns Channel with ID `c` from Server `s`
func GetServerChannelByID(c string, s *db.Server) (db.Channel, error) {
	if len(s.Channels) < 1 {
		return db.Channel{}, db.NotFound
	}
	for _, channel := range s.Channels {
		if channel.ID == c {
			return channel, nil
		}
	}
	return db.Channel{}, db.NotFound
}

// CheckChannel evaluates Channel `c` against desired game channel options
func CheckChannel(c *discordgo.Channel, g *discordgo.Guild) bool {
	server, err := db.GetServer(g.ID)
	if err != nil {
		logger.Log.Error(err.Error())
		return true
	}
	if c.ParentID != server.Category {
		return false
	}
	return true
}

// GetConfigChannelByName returns ChannelConfig with default name `c` from the config
func GetConfigChannelByName(c string) (config.ChannelConfig, error) {
	conf := config.Get()
	for _, channel := range conf.Channels {
		if channel.Name == c {
			return channel, nil
		}
	}
	return config.ChannelConfig{}, db.NotFound
}
