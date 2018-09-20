package builder

import (
	"fmt"

	"github.com/Noxdew/Knights-Of-Discord/db"
	"github.com/Noxdew/Knights-Of-Discord/logger"
	"github.com/Noxdew/Knights-Of-Discord/utils"
	"github.com/bwmarrin/discordgo"
)

// BuildServer creates a new server for the game to run on housed in Discord Server g
func BuildServer(g *discordgo.Guild) {
	logger.Log.Info("Building server...")
	db.CreateServer(db.Server{
		ID:    g.ID,
		Power: 0,
	})
	logger.Log.Info("Game started for %s (%s)", g.Name, g.ID)
}

// DestroyServer removes the Discord server from the DB
func DestroyServer(g *discordgo.Guild) {
	logger.Log.Info("Destroying server...")
	db.RemoveServer(g.ID)
	logger.Log.Info("Server destroyed.")
}

// BuildRoles creates Discord roles for the game to use
func BuildRoles(s *discordgo.Session, g *discordgo.Guild, roles []string) {
	logger.Log.Info("Building roles...")
	baseRole := utils.GetRoleByName(g, "@everyone")
	for _, role := range roles {
		r, err := s.GuildRoleCreate(g.ID)
		if err != nil {
			logger.Log.Error(err.Error())
		} else {
			s.GuildRoleEdit(g.ID, r.ID, role, 0, false, baseRole.Permissions, true)
			_, err := db.GetRole(g.ID, role)
			if err != nil && err != db.NotFound {
				logger.Log.Error(err.Error())
			}
			if err == db.NotFound {
				logger.Log.Info("Creating role %s", role)
				db.CreateRole(db.Role{
					ID:       r.ID,
					ServerID: g.ID,
					Type:     role,
				})
			} else {
				logger.Log.Info("Updating role %s", role)
				logger.Log.Debug("%s", r.ID)
				db.UpdateRole(db.Role{
					ID:       r.ID,
					ServerID: g.ID,
					Type:     role,
				})
			}
		}
	}
	logger.Log.Info("Roles built.")
}

// CheckRoles checks the condition of the game roles on an existing server and rebuilds them
func CheckRoles(s *discordgo.Session, g *discordgo.Guild) {
	logger.Log.Info("Checking Roles...")
	var rebuilding []string
	dbRoles, err := db.GetRoles(g.ID)
	if err != nil {
		logger.Log.Error(err.Error())
	}
	for _, role := range *dbRoles {
		r := utils.GetRoleByID(g, role.ID)
		if r == nil {
			rebuilding = append(rebuilding, role.Type)
		}
	}
	if len(rebuilding) > 0 {
		BuildRoles(s, g, rebuilding)
	}
	logger.Log.Info("Roles checked.")
}

// DestroyRoles removes the Discord game roles from server g
func DestroyRoles(s *discordgo.Session, g *discordgo.Guild) {
	logger.Log.Info("Destroying Roles...")
	dbRoles, err := db.GetRoles(g.ID)
	if err != nil {
		logger.Log.Error(err.Error())
	}
	for _, role := range *dbRoles {
		err := s.GuildRoleDelete(g.ID, role.ID)
		if err != nil {
			logger.Log.Error(err.Error())
		}
	}
	db.RemoveRoles(g.ID)
	logger.Log.Info("Roles destroyed.")
}

// BuildChannels checks if game channels exist and create them otherwise
func BuildChannels(s *discordgo.Session, g *discordgo.Guild, c *discordgo.Channel) {
	fmt.Println("Building channels...")
	ch := utils.GetChannelInCategory(g, c, "rules")
	if ch == nil {
		ch, err := s.GuildChannelCreate(g.ID, "rules", "0")
		if err != nil {
			fmt.Println(err.Error())
		}
		s.ChannelEditComplex(ch.ID, &discordgo.ChannelEdit{
			Topic:    "Game Rules and Information",
			Position: 0,
			ParentID: c.ID,
		})
	}
	ch = utils.GetChannelInCategory(g, c, "announcements")
	if ch == nil {
		ch, err := s.GuildChannelCreate(g.ID, "announcements", "0")
		if err != nil {
			fmt.Println(err.Error())
		}
		s.ChannelEditComplex(ch.ID, &discordgo.ChannelEdit{
			Topic:    "General Castle Information",
			Position: 1,
			ParentID: c.ID,
		})
	}
	ch = utils.GetChannelInCategory(g, c, "logs")
	if ch == nil {
		ch, err := s.GuildChannelCreate(g.ID, "logs", "0")
		if err != nil {
			fmt.Println(err.Error())
		}
		s.ChannelEditComplex(ch.ID, &discordgo.ChannelEdit{
			Topic:    "Activity Log For The Castle",
			Position: 2,
			ParentID: c.ID,
		})
	}
	ch = utils.GetChannelInCategory(g, c, "small-tavern")
	if ch == nil {
		ch, err := s.GuildChannelCreate(g.ID, "small-tavern", "0")
		if err != nil {
			fmt.Println(err.Error())
		}
		s.ChannelEditComplex(ch.ID, &discordgo.ChannelEdit{
			Topic:    "General Social Space For Villagers",
			Position: 3,
			ParentID: c.ID,
		})
	}
	ch = utils.GetChannelInCategory(g, c, "medium-tavern")
	if ch == nil {
		ch, err := s.GuildChannelCreate(g.ID, "medium-tavern", "0")
		if err != nil {
			fmt.Println(err.Error())
		}
		s.ChannelEditComplex(ch.ID, &discordgo.ChannelEdit{
			Topic:    "General Social Space For Esquires",
			Position: 4,
			ParentID: c.ID,
		})
	}
	ch = utils.GetChannelInCategory(g, c, "large-tavern")
	if ch == nil {
		ch, err := s.GuildChannelCreate(g.ID, "large-tavern", "0")
		if err != nil {
			fmt.Println(err.Error())
		}
		s.ChannelEditComplex(ch.ID, &discordgo.ChannelEdit{
			Topic:    "General Social Space For Knights",
			Position: 5,
			ParentID: c.ID,
		})
	}
	ch = utils.GetChannelInCategory(g, c, "outer-city")
	if ch == nil {
		ch, err := s.GuildChannelCreate(g.ID, "outer-city", "0")
		if err != nil {
			fmt.Println(err.Error())
		}
		s.ChannelEditComplex(ch.ID, &discordgo.ChannelEdit{
			Topic:    "Activity Center For Villagers",
			Position: 6,
			ParentID: c.ID,
		})
	}
	ch = utils.GetChannelInCategory(g, c, "inner-city")
	if ch == nil {
		ch, err := s.GuildChannelCreate(g.ID, "inner-city", "0")
		if err != nil {
			fmt.Println(err.Error())
		}
		s.ChannelEditComplex(ch.ID, &discordgo.ChannelEdit{
			Topic:    "Activity Center For Esquires",
			Position: 7,
			ParentID: c.ID,
		})
	}
	ch = utils.GetChannelInCategory(g, c, "castle")
	if ch == nil {
		ch, err := s.GuildChannelCreate(g.ID, "castle", "0")
		if err != nil {
			fmt.Println(err.Error())
		}
		s.ChannelEditComplex(ch.ID, &discordgo.ChannelEdit{
			Topic:    "Activity Center For Knights",
			Position: 8,
			ParentID: c.ID,
		})
	}
	fmt.Println("Channels built.")
}
