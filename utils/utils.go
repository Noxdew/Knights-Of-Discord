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

// GetDiscordRoleByName returns a Role with name `r` from guild `g`
func GetDiscordRoleByName(r string, g *discordgo.Guild) (*discordgo.Role, error) {
	for _, role := range g.Roles {
		if role.Name == r {
			return role, nil
		}
	}
	return nil, db.NotFound
}

// GetChannelRoleByName returns a Role with name `r` from Server Channel `c`
func GetChannelRoleByName(r string, c db.Channel) (string, error) {
	for _, role := range c.Roles {
		if role == r {
			return role, nil
		}
	}
	return "", db.NotFound
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

// GetChannelPermByID returns Permission with ID `p` in chanel `c`
func GetChannelPermByID(p string, c db.Channel) (db.Perm, error) {
	for _, perm := range c.Perms {
		if perm.ID == p {
			return perm, nil
		}
	}
	return db.Perm{}, db.NotFound
}

// CheckPermNeeded evaluates if a permission `p` is required for the game
func CheckPermNeeded(p *discordgo.PermissionOverwrite, channel db.Channel) bool {
	if p.ID == "487744442531315712" {
		return true
	}
	_, err := GetChannelPermByID(p.ID, channel)
	if err != nil {
		return false
	}
	return true
}

// CheckPerm evaluates PermissionOverwrite `p` against desired game perm options
func CheckPerm(p *discordgo.PermissionOverwrite, c db.Channel) bool {
	perm, err := GetChannelPermByID(p.ID, c)
	if err != nil {
		logger.Log.Error(err.Error())
		return true
	}
	if perm.Allow != p.Allow {
		return false
	}
	return true
}

// CheckPermExists evaluates if Server Perm `p` exists in Discord Channel `c`
func CheckPermExists(p db.Perm, c *discordgo.Channel) bool {
	for _, perm := range c.PermissionOverwrites {
		if perm.ID == p.ID {
			return true
		}
	}
	return false
}
