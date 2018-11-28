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
	buildMessages(server, s, g)

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

	// Set Message Channels
	server.SetMessageChannel()

	logger.Log.Info("Permissions for server %s (%s) successfully built.", g.Name, g.ID)
}

func buildMessages(server *structure.Server, s *discordgo.Session, g *discordgo.Guild) {
	logger.Log.Info("Building Messages for Guild %s (id: %s)...", g.Name, g.ID)

	// Send Messages
	for _, message := range server.Messages.GetMessages() {
		// Create embed
		embed := buildEmbed(message)

		// Send Message
		m, err := s.ChannelMessageSendEmbed(message.Channel, embed)
		if err != nil {
			logger.Log.Error(err.Error())
			return
		}

		// Update Message object
		message.ID = m.ID

		// Add reactions
		buildReactions(message, s, g)
	}

	logger.Log.Info("Messages for server %s (%s) successfully built.", g.Name, g.ID)
}

func buildEmbed(m *structure.Message) *discordgo.MessageEmbed {
	embed := discordgo.MessageEmbed{
		Title:       m.Title,
		Description: m.Description,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: m.Icon,
		},
		Footer: &discordgo.MessageEmbedFooter{
			IconURL: "https://cdn.discordapp.com/attachments/512302843437252611/512302951814004752/ac6918be09a389876ee5663d6b08b55a.png",
			Text:    m.Footer,
		},
	}

	// Set Color
	if m.Type == "info" {
		embed.Color = 16098851
	} else if m.Type == "card" {
		embed.Color = 6711705
	} else if m.Type == "active" {
		embed.Color = 8311585
	} else if m.Type == "danger" {
		embed.Color = 13632027
	} else if m.Type == "shop" {
		embed.Color = 12390624
	} else {
		embed.Color = 4868682
	}

	// Set Fields
	for _, field := range m.Fields {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   field.Title,
			Value:  field.Value,
			Inline: field.Inline,
		})
	}

	return &embed
}

func buildReactions(message *structure.Message, s *discordgo.Session, g *discordgo.Guild) {
	if message.Type == "info" {
		err := s.MessageReactionAdd(message.Channel, message.ID, ":kod:514099648949125153")
		if err != nil {
			logger.Log.Error(err.Error())
		}
	}
}

// AddUser assigns game roles to a user in the guild
func AddUser(server *structure.Server, s *discordgo.Session, g *discordgo.Guild, u string) {
	logger.Log.Info("User %s joining server %s (%s)...", u, g.Name, g.ID)

	err := s.GuildMemberRoleAdd(server.ID, u, server.Roles.Villager.ID)
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}
	err = db.AddUser(server, &structure.User{
		ID:           u,
		Role:         server.Roles.Villager.ID,
		Contribution: 0,
	})
	if err != nil {
		logger.Log.Error(err.Error())
	}

	logger.Log.Info("User %s successfully joined %s (%s).", u, g.Name, g.ID)
}
