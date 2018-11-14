package builder

import (
	"github.com/Noxdew/Knights-Of-Discord/db"
	"github.com/Noxdew/Knights-Of-Discord/logger"
	"github.com/Noxdew/Knights-Of-Discord/structure"
	"github.com/bwmarrin/discordgo"
)

// BuildServer initializes a new game on guild `g`
func BuildServer(server *structure.Server, s *discordgo.Session, g *discordgo.Guild) {
	logger.Log.Info("Building server for Guild %s (id: %s)...", g.Name, g.ID)

	// Create Server object
	server.BuildServer(g)

	// Build Discord Structure
	buildRoles(server, s, g)
	buildCategory(server, s, g)
	buildChannels(server, s, g)

	// Upload to DB
	err := db.CreateServer(server)
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}

	logger.Log.Info("Server for Guild %s (id: %s) successfully built.", g.Name, g.ID)
}

// DestroyServer removes game instance from guild `g`
func DestroyServer(server *structure.Server, s *discordgo.Session, g *discordgo.Guild) {
	logger.Log.Info("Destroying server %s (%s)...", g.Name, g.ID)

	// Close the game running on Server
	err := db.UpdateServerPlaying(server, false)
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}

	// Destroy Discord Structure
	destroyChannels(server, s, g)
	destroyCategory(server, s, g)
	destroyRoles(server, s, g)

	// Leave Guild
	err = s.GuildLeave(g.ID)
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}

	// Clear Server from DB
	err = db.DeleteServer(server)
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}

	logger.Log.Info("Server %s (%s) successfully destroyed.", g.Name, g.ID)
}

func buildRoles(server *structure.Server, s *discordgo.Session, g *discordgo.Guild) {
	logger.Log.Info("Building roles for server %s (%s)...", g.Name, g.ID)

	// Get @everyone Role
	var everyone *discordgo.Role
	dRoles, err := s.GuildRoles(g.ID)
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}
	for _, role := range dRoles {
		if role.Name == "@everyone" {
			everyone = role
		}
	}

	roles := server.Roles.GetRoles()
	for i, r := range roles {
		if r.DefName == "everyone" {
			// Save @everyone to object
			roles[i].ID = everyone.ID
		} else {
			// Create Role
			role, err := s.GuildRoleCreate(g.ID)
			if err != nil {
				logger.Log.Error(err.Error())
				return
			}
			_, err = s.GuildRoleEdit(g.ID, role.ID, r.DefName, 0, false, r.Permission, true)
			if err != nil {
				logger.Log.Error(err.Error())
				return
			}
			roles[i].ID = role.ID
		}
	}

	logger.Log.Info("Roles for server %s (%s) successfully built.", g.Name, g.ID)
}

func destroyRoles(server *structure.Server, s *discordgo.Session, g *discordgo.Guild) {
	logger.Log.Info("Destroying Roles from server %s (%s)...", g.Name, g.ID)

	roles := server.Roles.GetRoles()
	for _, role := range roles {
		if role.Level > 0 {
			err := s.GuildRoleDelete(g.ID, role.ID)
			if err != nil {
				logger.Log.Error(err.Error())
			}
		}
	}

	logger.Log.Info("Roles from server %s (%s) successfully destroyed.", g.Name, g.ID)
}

func buildCategory(server *structure.Server, s *discordgo.Session, g *discordgo.Guild) {
	logger.Log.Info("Building category for server %s (%s)", g.Name, g.ID)

	category, err := s.GuildChannelCreate(g.ID, server.Category.DefName, "4")
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}
	server.Category.ID = category.ID

	logger.Log.Info("Category for server %s (%s) successfully built.", g.Name, g.ID)
}

func destroyCategory(server *structure.Server, s *discordgo.Session, g *discordgo.Guild) {
	logger.Log.Info("Destroying Category from server %s (%s)...", g.Name, g.ID)

	_, err := s.ChannelDelete(server.Category.ID)
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}

	logger.Log.Info("Category from server %s (%s) successfully destroyed.", g.Name, g.ID)
}

func buildChannels(server *structure.Server, s *discordgo.Session, g *discordgo.Guild) {
	logger.Log.Info("Building channels for server %s (%s)...", g.Name, g.ID)

	roles := server.Roles.GetRoles()
	channels := server.Channels.GetChannels()
	for i, c := range channels {
		// Create Channel
		channel, err := s.GuildChannelCreate(g.ID, c.DefName, "0")
		if err != nil {
			logger.Log.Error(err.Error())
			return
		}
		_, err = s.ChannelEditComplex(channel.ID, &discordgo.ChannelEdit{
			ParentID: server.Category.ID,
		})
		if err != nil {
			logger.Log.Error(err.Error())
			return
		}
		channels[i].ID = channel.ID
		// Add Permissions

		// Set Bot Permissions
		bot, err := s.User("@me")
		err = s.ChannelPermissionSet(c.ID, bot.ID, "member", server.BotPerm, 0)
		if err != nil {
			logger.Log.Error(err.Error())
			return
		}

		// Set Permissions
		for _, r := range roles {
			if c.Level < 0 {
				// Hub Channel
				if r.Level < 0 {
					err := s.ChannelPermissionSet(c.ID, r.ID, "role", c.Allow, c.Deny)
					if err != nil {
						logger.Log.Error(err.Error())
						return
					}
					break
				}
			} else {
				// Game Channels
				if r.Level < 0 {
					err := s.ChannelPermissionSet(c.ID, r.ID, "role", 0, server.BotPerm)
					if err != nil {
						logger.Log.Error(err.Error())
						return
					}
				} else if r.Level >= c.Level {
					err := s.ChannelPermissionSet(c.ID, r.ID, "role", c.Allow, c.Deny)
					if err != nil {
						logger.Log.Error(err.Error())
						return
					}
				}
			}
		}

		// Send Messages
		for _, message := range c.Messages.GetMessages() {
			embed := discordgo.MessageEmbed{
				Title:       message.Title,
				Description: message.Description,
				Color:       message.Color,
				Footer: &discordgo.MessageEmbedFooter{
					Text: "Knights of Discord",
				},
				Author: &discordgo.MessageEmbedAuthor{
					Name:    bot.Username,
					IconURL: bot.AvatarURL("32"),
				},
			}
			for _, field := range message.Fields {
				embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
					Name:  field.Title,
					Value: field.Value,
				})
			}
			if message.Icon != "" {
				embed.Thumbnail = &discordgo.MessageEmbedThumbnail{
					URL: message.Icon,
				}
			}
			_, err := s.ChannelMessageSendEmbed(c.ID, &embed)
			if err != nil {
				logger.Log.Error(err.Error())
				return
			}
		}
	}

	logger.Log.Info("Channels for server %s (%s) successfully built.", g.Name, g.ID)
}

func destroyChannels(server *structure.Server, s *discordgo.Session, g *discordgo.Guild) {
	logger.Log.Info("Destroying Roles from server %s (%s)...", g.Name, g.ID)

	channels := server.Channels.GetChannels()
	for _, channel := range channels {
		_, err := s.ChannelDelete(channel.ID)
		if err != nil {
			logger.Log.Error(err.Error())
		}
	}

	logger.Log.Info("Roles from server %s (%s) successfully destroyed.", g.Name, g.ID)
}
