package builder

import (
	"github.com/Noxdew/Knights-Of-Discord/db"
	"github.com/Noxdew/Knights-Of-Discord/logger"
	"github.com/Noxdew/Knights-Of-Discord/structure"
	"github.com/bwmarrin/discordgo"
)

// BuildServer initializes a new game on guild `g`
func BuildServer(server *structure.Server, s *discordgo.Session, g *discordgo.Guild) {
	logger.Log.Info("Building Server for Guild %s (id: %s)...", g.Name, g.ID)

	// Create Server object
	server.BuildServer(g)

	// Build Discord Structure
	buildRoles(server, s, g)
	buildCategory(server, s, g)
	buildChannels(server, s, g)
	buildPermissions(server, s, g)
	// buildMessages(server, s, g)

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
	logger.Log.Info("Destroying Server %s (%s)...", g.Name, g.ID)

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
	logger.Log.Info("Building Roles for server %s (%s)...", g.Name, g.ID)

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
		if role.Level >= 0 {
			err := s.GuildRoleDelete(g.ID, role.ID)
			if err != nil {
				logger.Log.Error(err.Error())
			}
		}
	}

	logger.Log.Info("Roles from server %s (%s) successfully destroyed.", g.Name, g.ID)
}

func buildCategory(server *structure.Server, s *discordgo.Session, g *discordgo.Guild) {
	logger.Log.Info("Building Category for server %s (%s)", g.Name, g.ID)

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
	logger.Log.Info("Building Channels for server %s (%s)...", g.Name, g.ID)

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
			Position: i,
		})
		if err != nil {
			logger.Log.Error(err.Error())
			return
		}
		channels[i].ID = channel.ID
	}

	logger.Log.Info("Channels for server %s (%s) successfully built.", g.Name, g.ID)
}

func destroyChannels(server *structure.Server, s *discordgo.Session, g *discordgo.Guild) {
	logger.Log.Info("Destroying Channels from server %s (%s)...", g.Name, g.ID)

	channels := server.Channels.GetChannels()
	for _, channel := range channels {
		_, err := s.ChannelDelete(channel.ID)
		if err != nil {
			logger.Log.Error(err.Error())
		}
	}

	logger.Log.Info("Channels from server %s (%s) successfully destroyed.", g.Name, g.ID)
}

func buildPermissions(server *structure.Server, s *discordgo.Session, g *discordgo.Guild) {
	logger.Log.Info("Building Permissions for server %s (%s)...", g.Name, g.ID)

	channels := server.Channels.GetChannels()
	bot, err := s.User("@me")
	for i, channel := range channels {
		// Set Bot Permissions
		err = s.ChannelPermissionSet(channel.ID, bot.ID, "member", server.BotPerm, 0)
		if err != nil {
			logger.Log.Error(err.Error())
			return
		}
		channels[i].Permissions = append(channels[i].Permissions, structure.Perm{
			Role:  bot.ID,
			Allow: server.BotPerm,
			Deny:  0,
		})

		if channel.Level < 0 {
			// Hub Channel
			err := s.ChannelPermissionSet(channel.ID, server.Roles.Everyone.ID, "role", channel.Allow, channel.Deny)
			if err != nil {
				logger.Log.Error(err.Error())
				return
			}
			channels[i].Permissions = append(channels[i].Permissions, structure.Perm{
				Role:  server.Roles.Everyone.ID,
				Allow: channel.Allow,
				Deny:  channel.Deny,
			})
		} else {
			err := s.ChannelPermissionSet(channel.ID, server.Roles.Everyone.ID, "role", 0, server.BotPerm)
			if err != nil {
				logger.Log.Error(err.Error())
				return
			}
			channels[i].Permissions = append(channels[i].Permissions, structure.Perm{
				Role:  server.Roles.Everyone.ID,
				Allow: 0,
				Deny:  server.BotPerm,
			})
			// Game Channel
			for _, role := range server.Roles.GetRoles() {
				if channel.Level <= role.Level {
					err := s.ChannelPermissionSet(channel.ID, role.ID, "role", channel.Allow, channel.Deny)
					if err != nil {
						logger.Log.Error(err.Error())
						return
					}
					channels[i].Permissions = append(channels[i].Permissions, structure.Perm{
						Role:  role.ID,
						Allow: channel.Allow,
						Deny:  channel.Deny,
					})
				}
			}
		}
	}

	logger.Log.Info("Permissions for server %s (%s) successfully built.", g.Name, g.ID)
}

// func buildMessages() {
// 	// Send Messages
// 	for _, message := range c.Messages.GetMessages() {
// 		embed := discordgo.MessageEmbed{
// 			Title:       message.Title,
// 			Description: message.Description,
// 			Color:       message.Color,
// 			Footer: &discordgo.MessageEmbedFooter{
// 				Text: "Knights of Discord",
// 			},
// 			Author: &discordgo.MessageEmbedAuthor{
// 				Name:    bot.Username,
// 				IconURL: bot.AvatarURL("32"),
// 			},
// 		}
// 		for _, field := range message.Fields {
// 			embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
// 				Name:  field.Title,
// 				Value: field.Value,
// 			})
// 		}
// 		if message.Icon != "" {
// 			embed.Thumbnail = &discordgo.MessageEmbedThumbnail{
// 				URL: message.Icon,
// 			}
// 		}
// 		_, err := s.ChannelMessageSendEmbed(c.ID, &embed)
// 		if err != nil {
// 			logger.Log.Error(err.Error())
// 			return
// 		}
// 	}
// }
