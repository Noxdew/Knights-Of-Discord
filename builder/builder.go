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
	db.CreateServer(db.Server{
		ID:      g.ID,
		Power:   0,
		Checked: true,
		Roles:   []db.Role{},
	})
	BuildRoles(s, g)
	logger.Log.Info("Server %s (%s) successfully built.", g.Name, g.ID)
}

// DestroyServer removes game instance from guild `g`
func DestroyServer(s *discordgo.Session, g *discordgo.Guild) {
	logger.Log.Info("Destroying server %s (%s)...", g.Name, g.ID)
	DestroyRoles(s, g)
	db.RemoveServer(g.ID)
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
	}
	// Assign new Role options
	_, err = s.GuildRoleEdit(g.ID, r.ID, role, 0, false, config.Get().RolePerm, true)
	if err != nil {
		logger.Log.Error(err.Error())
	}
	server, err := db.GetServer(g.ID)
	if err != nil {
		logger.Log.Error(err.Error())
	}
	// Upload to DB
	if len(server.Roles) > 0 {
		for _, rS := range server.Roles {
			if rS.Type == role {
				err = db.UpdateRole(db.Role{
					ID:   r.ID,
					Type: role,
				}, g.ID)
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
