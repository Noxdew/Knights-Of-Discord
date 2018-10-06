package builder

import (
	"github.com/Noxdew/Knights-Of-Discord/config"
	"github.com/Noxdew/Knights-Of-Discord/db"
	"github.com/Noxdew/Knights-Of-Discord/logger"

	"github.com/bwmarrin/discordgo"
)

// BuildServer initializes a new game on guild `g`
func BuildServer(s *discordgo.Session, g *discordgo.Guild) {
	logger.Log.Info("Building server %s (%s)...", g.Name, g.ID)
	err := db.CreateServer(db.Server{
		ID:       g.ID,
		Checked:  true,
		Playing:  true,
		Power:    0,
		Roles:    []db.Role{},
		Category: "",
		Channels: []db.Channel{},
	})
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}
	BuildRoles(s, g)
	BuildChannels(s, g)
	logger.Log.Info("Server %s (%s) successfully built.", g.Name, g.ID)
}

// DestroyServer removes game instance from guild `g`
func DestroyServer(s *discordgo.Session, g *discordgo.Guild) {
	logger.Log.Info("Destroying server %s (%s)...", g.Name, g.ID)
	// Flag the server for destruction
	err := db.UpdateServerStatus(g.ID, false)
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}
	DestroyChannels(s, g)
	DestroyRoles(s, g)
	err = db.RemoveServer(g.ID)
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}
	logger.Log.Info("Server %s (%s) successfully destroyed.", g.Name, g.ID)
}

// BuildRoles initializes new game roles for `g`
func BuildRoles(s *discordgo.Session, g *discordgo.Guild) {
	logger.Log.Info("Building roles for server %s (%s)...", g.Name, g.ID)
	for _, role := range config.Get().Roles {
		BuildRole(s, g, role)
	}
	logger.Log.Info("Roles for server %s (%s) successfully built.", g.Name, g.ID)
}

// DestroyRoles removes the Discord game roles from server g
func DestroyRoles(s *discordgo.Session, g *discordgo.Guild) {
	logger.Log.Info("Destroying Roles from server %s (%s)...", g.Name, g.ID)
	server, err := db.GetServer(g.ID)
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}
	for _, role := range server.Roles {
		// Remove role from server
		err := s.GuildRoleDelete(g.ID, role.ID)
		if err != nil {
			logger.Log.Error(err.Error())
			continue
		}
	}
	logger.Log.Info("Roles from server %s (%s) successfully destroyed.", g.Name, g.ID)
}

// BuildRole creates a single missing Discord game role
func BuildRole(s *discordgo.Session, g *discordgo.Guild, role string) {
	logger.Log.Info("Building role %s", role)
	// Create new Role
	r, err := s.GuildRoleCreate(g.ID)
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}
	// Assign new Role options
	_, err = s.GuildRoleEdit(g.ID, r.ID, role, 0, false, config.Get().RolePerm, true)
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}
	server, err := db.GetServer(g.ID)
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}
	// Upload to DB
	if len(server.Roles) > 0 {
		for _, rS := range server.Roles {
			if rS.Type == role {
				err = db.UpdateRole(db.Role{
					ID:   r.ID,
					Type: role,
				}, g.ID)
				if err != nil {
					logger.Log.Error(err.Error())
				}
				return
			}
		}
	}
	err = db.CreateRole(db.Role{
		ID:   r.ID,
		Type: role,
	}, g.ID)
	if err != nil {
		logger.Log.Error(err.Error())
	}
}

// FixRole rebuilds Discord game role to desired options
func FixRole(s *discordgo.Session, g *discordgo.Guild, r *discordgo.Role) {
	logger.Log.Info("Fixing role %s for server %s (%s)", r.Name, g.Name, g.ID)
	_, err := s.GuildRoleEdit(g.ID, r.ID, r.Name, r.Color, false, config.Get().RolePerm, true)
	if err != nil {
		logger.Log.Error(err.Error())
	}
}

// BuildChannels initializes new game channels for `g`
func BuildChannels(s *discordgo.Session, g *discordgo.Guild) {
	// Create Category
	logger.Log.Info("Building category for server %s (%s)", g.Name, g.ID)
	c, err := s.GuildChannelCreate(g.ID, "Knights of Discord", "4")
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}
	// Upload category to DB
	err = db.CreateCategory(c.ID, g.ID)
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}
	logger.Log.Info("Category for server %s (%s) successfully built.", g.Name, g.ID)
	// Create Channels
	logger.Log.Info("Building channels for server %s (%s)...", g.Name, g.ID)
	for _, channel := range config.Get().Channels {
		BuildChannel(s, g, channel)
	}
	logger.Log.Info("Channels for server %s (%s) successfully built.", g.Name, g.ID)
}

// DestroyChannels removes the Discord game channels from server g
func DestroyChannels(s *discordgo.Session, g *discordgo.Guild) {
	logger.Log.Info("Destroying Channels from server %s (%s)...", g.Name, g.ID)
	server, err := db.GetServer(g.ID)
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}
	for _, channel := range server.Channels {
		// Remove role from server
		_, err := s.ChannelDelete(channel.ID)
		if err != nil {
			logger.Log.Error(err.Error())
			continue
		}
	}
	_, err = s.ChannelDelete(server.Category)
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}
	logger.Log.Info("Channels from server %s (%s) successfully destroyed.", g.Name, g.ID)
}

// BuildChannel creates a single missing Discord game channel
func BuildChannel(s *discordgo.Session, g *discordgo.Guild, channel config.ChannelConfig) {
	logger.Log.Info("Building channel %s", channel.Name)
	// Create new Channel
	c, err := s.GuildChannelCreate(g.ID, channel.Name, "0")
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}
	// Assign new Channel options
	server, err := db.GetServer(g.ID)
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}
	_, err = s.ChannelEditComplex(c.ID, &discordgo.ChannelEdit{
		NSFW:     false,
		ParentID: server.Category,
	})
	// Add Channel permissions
	err = s.ChannelPermissionSet(c.ID, "487744442531315712", "member", 871890257, 0)
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}
	if channel.Type == "hub" {
		for _, role := range g.Roles {
			err = s.ChannelPermissionSet(c.ID, role.ID, "role", 328768, 871561489)
			if err != nil {
				logger.Log.Error(err.Error())
				return
			}
		}
	} else if channel.Type == "action" {
		for _, role := range g.Roles {
			err = s.ChannelPermissionSet(c.ID, role.ID, "role", 0, 871890257)
			if err != nil {
				logger.Log.Error(err.Error())
				return
			}
		}
		hash := make(map[string]bool)
		for _, role := range channel.Role {
			hash[role] = true
		}
		for _, role := range server.Roles {
			if _, ok := hash[role.Type]; ok {
				err = s.ChannelPermissionSet(c.ID, role.ID, "role", 328768, 871561489)
				if err != nil {
					logger.Log.Error(err.Error())
					return
				}
			}
		}
	} else if channel.Type == "social" {
		for _, role := range g.Roles {
			err = s.ChannelPermissionSet(c.ID, role.ID, "role", 0, 871890257)
			if err != nil {
				logger.Log.Error(err.Error())
				return
			}
		}
		hash := make(map[string]bool)
		for _, role := range channel.Role {
			hash[role] = true
		}
		for _, role := range server.Roles {
			if _, ok := hash[role.Type]; ok {
				err = s.ChannelPermissionSet(c.ID, role.ID, "role", 379968, 871510289)
				if err != nil {
					logger.Log.Error(err.Error())
					return
				}
			}
		}
	}
	// Upload to DB
	if len(server.Channels) > 0 {
		for _, cS := range server.Channels {
			if cS.Name == channel.Name {
				err = db.UpdateChannel(db.Channel{
					ID:   c.ID,
					Name: channel.Name,
					Type: channel.Type,
				}, g.ID)
				if err != nil {
					logger.Log.Error(err.Error())
				}
				return
			}
		}
	}
	err = db.CreateChannel(db.Channel{
		ID:   c.ID,
		Name: channel.Name,
		Type: channel.Type,
	}, g.ID)
	if err != nil {
		logger.Log.Error(err.Error())
	}
}
