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

	// Build Discord Structure
	buildRoles(server, s, g)
	buildCategory(server, s, g)
	buildChannels(server, s, g)
	buildPermissions(server, s, g)
	buildMessages(server, s, g)

	// Update Server object
	server.ID = g.ID
	server.Playing = true
	server.Users = []*structure.User{}

	// Upload to DB
	err := db.CreateServer(server)
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}

	logger.Log.Info("Server for Guild %s (id: %s) successfully built.", g.Name, g.ID)
	return
}

// DestroyServer removes game instance from guild `g`
func DestroyServer(server *structure.Server, s *discordgo.Session, g *discordgo.Guild) {
	logger.Log.Info("Destroying Server %s (%s)...", g.Name, g.ID)

	// Update Server Object
	server.Playing = false
	err := db.UpdateServerPlaying(server)
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

	for _, r := range server.Roles {
		// Create Discord Role
		role, err := s.GuildRoleCreate(g.ID)
		if err != nil {
			logger.Log.Error(err.Error())
			return
		}
		_, err = s.GuildRoleEdit(g.ID, role.ID, r.DefaultName, 0, r.Hoist, server.RolePerm, r.Mentionable)
		if err != nil {
			logger.Log.Error(err.Error())
			return
		}

		// Update Server object
		r.ID = role.ID
	}

	// Update Server Object
	for _, role := range g.Roles {
		if role.Name == "@everyone" {
			server.EveryoneRole = role.ID
			break
		}
	}

	logger.Log.Info("Roles for server %s (%s) successfully built.", g.Name, g.ID)
}

func destroyRoles(server *structure.Server, s *discordgo.Session, g *discordgo.Guild) {
	logger.Log.Info("Destroying Roles from server %s (%s)...", g.Name, g.ID)

	for _, role := range server.Roles {
		err := s.GuildRoleDelete(g.ID, role.ID)
		if err != nil {
			logger.Log.Error(err.Error())
		}
	}

	logger.Log.Info("Roles from server %s (%s) successfully destroyed.", g.Name, g.ID)
}

func buildCategory(server *structure.Server, s *discordgo.Session, g *discordgo.Guild) {
	logger.Log.Info("Building Category for server %s (%s)", g.Name, g.ID)

	// Create Discord Category
	category, err := s.GuildChannelCreate(g.ID, server.Category.DefaultName, "4")
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}

	// Update Server Object
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

	for _, c := range server.Channels {
		// Create Discord Channel
		channel, err := s.GuildChannelCreate(g.ID, c.DefaultName, "0")
		if err != nil {
			logger.Log.Error(err.Error())
			return
		}
		_, err = s.ChannelEditComplex(channel.ID, &discordgo.ChannelEdit{
			ParentID: server.Category.ID,
			Position: c.Position,
			Topic:    c.Topic,
		})
		if err != nil {
			logger.Log.Error(err.Error())
			return
		}

		// Update Server object
		c.ID = channel.ID
	}

	logger.Log.Info("Channels for server %s (%s) successfully built.", g.Name, g.ID)
}

func destroyChannels(server *structure.Server, s *discordgo.Session, g *discordgo.Guild) {
	logger.Log.Info("Destroying Channels from server %s (%s)...", g.Name, g.ID)

	for _, channel := range server.Channels {
		_, err := s.ChannelDelete(channel.ID)
		if err != nil {
			logger.Log.Error(err.Error())
		}
	}

	logger.Log.Info("Channels from server %s (%s) successfully destroyed.", g.Name, g.ID)
}

func buildPermissions(server *structure.Server, s *discordgo.Session, g *discordgo.Guild) {
	logger.Log.Info("Building Permissions for server %s (%s)...", g.Name, g.ID)

	bot, err := s.User("@me")
	for _, channel := range server.Channels {
		// Set Discord Permissions for Bot
		err = s.ChannelPermissionSet(channel.ID, bot.ID, "member", server.BotPerm, 0)
		if err != nil {
			logger.Log.Error(err.Error())
			return
		}

		if channel.Tier == 0 {
			// Set Discord Permission for Hub Channel
			err := s.ChannelPermissionSet(channel.ID, server.EveryoneRole, "role", server.ActionPerm, (server.BotPerm - server.ActionPerm))
			if err != nil {
				logger.Log.Error(err.Error())
				return
			}
		} else {
			// Set Discord Permissions for Game/Social Channel
			// @everyone Permissions
			err := s.ChannelPermissionSet(channel.ID, server.EveryoneRole, "role", 0, server.BotPerm)
			if err != nil {
				logger.Log.Error(err.Error())
				return
			}

			// Game Role Permissions
			for _, role := range server.Roles {
				if role.Tier >= channel.Tier {
					if channel.Type == "social" {
						// Social Channel
						err := s.ChannelPermissionSet(channel.ID, role.ID, "role", server.SocialPerm, (server.BotPerm - server.SocialPerm))
						if err != nil {
							logger.Log.Error(err.Error())
							return
						}
					} else if channel.Type == "action" {
						// Game Channel
						err := s.ChannelPermissionSet(channel.ID, role.ID, "role", server.ActionPerm, (server.BotPerm - server.ActionPerm))
						if err != nil {
							logger.Log.Error(err.Error())
							return
						}
					}
				}
			}
		}
	}

	logger.Log.Info("Permissions for server %s (%s) successfully built.", g.Name, g.ID)
}

func buildMessages(server *structure.Server, s *discordgo.Session, g *discordgo.Guild) {
	logger.Log.Info("Building Messages for Guild %s (id: %s)...", g.Name, g.ID)

	// Send Messages
	for _, message := range server.Messages {
		// Create embed
		embed := BuildEmbed(message)

		// Send Message
		m, err := s.ChannelMessageSendEmbed(server.Channels["rules"].ID, embed)
		if err != nil {
			logger.Log.Error(err.Error())
			return
		}

		// Update Message object
		message.ID = m.ID
		message.ChannelID = m.ChannelID

		// Add reactions
		if message.Type == "info" {
			err := s.MessageReactionAdd(message.ChannelID, message.ID, server.Actions["join"])
			if err != nil {
				logger.Log.Error(err.Error())
			}
		}
	}

	logger.Log.Info("Messages for server %s (%s) successfully built.", g.Name, g.ID)
}

// BuildEmbed creates a new Discord Embed
func BuildEmbed(m *structure.Message) *discordgo.MessageEmbed {
	// Create Discord Embed
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
	} else if m.Type == "system" {
		embed.Color = 10372089
	} else {
		embed.Color = 4868682
	}

	// Set Fields
	for _, field := range m.Fields {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:  field.Title,
			Value: field.Value,
		})
	}

	return &embed
}
