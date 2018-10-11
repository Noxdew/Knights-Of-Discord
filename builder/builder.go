package builder

import (
	"github.com/Noxdew/Knights-Of-Discord/config"
	"github.com/Noxdew/Knights-Of-Discord/db"
	"github.com/Noxdew/Knights-Of-Discord/logger"
	"github.com/Noxdew/Knights-Of-Discord/utils"
	"github.com/bwmarrin/discordgo"
)

// BuildServer initializes a new game on guild `g`
func BuildServer(s *discordgo.Session, g *discordgo.Guild) {
	logger.Log.Info("Building server %s (%s)...", g.Name, g.ID)
	err := db.CreateServer(db.Server{
		ID:       g.ID,
		Checked:  true,
		Playing:  true,
		Roles:    []db.Role{},
		Category: "",
		Channels: []db.Channel{},
	})
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}
	BuildRoles(s, g)
	BuildCategory(s, g)
	BuildChannels(s, g)
	logger.Log.Info("Server %s (%s) successfully built.", g.Name, g.ID)
}

// DestroyServer removes game instance from guild `g`
func DestroyServer(s *discordgo.Session, g *discordgo.Guild) {
	logger.Log.Info("Destroying server %s (%s)...", g.Name, g.ID)
	// Flag the server for destruction
	err := db.UpdateServerPlaying(g.ID, false)
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}
	DestroyChannels(s, g)
	DestroyCategory(s, g)
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

// DestroyRoles removes the Discord game roles from server `g`
func DestroyRoles(s *discordgo.Session, g *discordgo.Guild) {
	logger.Log.Info("Destroying Roles from server %s (%s)...", g.Name, g.ID)
	server, err := db.GetServer(g.ID)
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}
	for _, role := range server.Roles {
		err := s.GuildRoleDelete(g.ID, role.ID)
		if err != nil {
			logger.Log.Error(err.Error())
		}
	}
	logger.Log.Info("Roles from server %s (%s) successfully destroyed.", g.Name, g.ID)
}

// BuildRole Creates role `r` for Guild `g`
func BuildRole(s *discordgo.Session, g *discordgo.Guild, r string) {
	logger.Log.Info("Building role %s for server %s", r, g.Name)
	// Create new Role
	role, err := s.GuildRoleCreate(g.ID)
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}
	// Assign new Role options
	_, err = s.GuildRoleEdit(g.ID, role.ID, r, 0, false, config.Get().RolePerm, true)
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}
	// Upload to DB
	server, err := db.GetServer(g.ID)
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}
	_, err = utils.GetServerRoleByName(r, server)
	if err != nil {
		err = db.CreateRole(db.Role{
			ID:      role.ID,
			DefName: r,
		}, g.ID)
		if err != nil {
			logger.Log.Error(err.Error())
		}
	} else {
		err = db.UpdateRole(db.Role{
			ID:      role.ID,
			DefName: r,
		}, g.ID)
		if err != nil {
			logger.Log.Error(err.Error())
		}
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

// BuildCategory initializes new game category for `g`
func BuildCategory(s *discordgo.Session, g *discordgo.Guild) {
	// Create Category
	logger.Log.Info("Building category for server %s (%s)", g.Name, g.ID)
	c, err := s.GuildChannelCreate(g.ID, "Knights of Discord", "4")
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}
	// Rearange game channels
	server, err := db.GetServer(g.ID)
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}
	if len(server.Channels) > 0 {
		for _, channel := range server.Channels {
			_, err = s.ChannelEditComplex(channel.ID, &discordgo.ChannelEdit{
				ParentID: c.ID,
			})
			if err != nil {
				logger.Log.Error(err.Error())
				return
			}
		}
	}
	// Upload category to DB
	err = db.CreateCategory(c.ID, g.ID)
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}
	logger.Log.Info("Category for server %s (%s) successfully built.", g.Name, g.ID)
}

// DestroyCategory removes the Discord game category from server g
func DestroyCategory(s *discordgo.Session, g *discordgo.Guild) {
	logger.Log.Info("Destroying Category from server %s (%s)...", g.Name, g.ID)
	server, err := db.GetServer(g.ID)
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}
	_, err = s.ChannelDelete(server.Category)
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}
	logger.Log.Info("Category from server %s (%s) successfully destroyed.", g.Name, g.ID)
}

// BuildChannels initializes new game channels for `g`
func BuildChannels(s *discordgo.Session, g *discordgo.Guild) {
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
		_, err := s.ChannelDelete(channel.ID)
		if err != nil {
			logger.Log.Error(err.Error())
		}
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
		ParentID: server.Category,
	})
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}
	// Upload to DB
	_, err = utils.GetServerChannelByName(channel.Name, server)
	if err != nil {
		err = db.CreateChannel(db.Channel{
			ID:      c.ID,
			DefName: channel.Name,
			Type:    channel.Type,
			Perms:   []db.Perm{},
		}, g.ID)
		if err != nil {
			logger.Log.Error(err.Error())
		}
	} else {
		err = db.UpdateChannel(db.Channel{
			ID:      c.ID,
			DefName: channel.Name,
		}, g.ID)
		if err != nil {
			logger.Log.Error(err.Error())
		}
	}
	// Add Permissions
	// TODO
}

// FixChannel rebuilds Discord game channel to desired options
func FixChannel(s *discordgo.Session, g *discordgo.Guild, c *discordgo.Channel) {
	logger.Log.Info("Fixing channel %s for server %s (%s)", c.Name, g.Name, g.ID)
	server, err := db.GetServer(g.ID)
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}
	_, err = s.ChannelEditComplex(c.ID, &discordgo.ChannelEdit{
		ParentID: server.Category,
	})
	if err != nil {
		logger.Log.Error(err.Error())
	}
}
