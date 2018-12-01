package handlers

import (
	"strings"

	"github.com/Noxdew/Knights-Of-Discord/builder"
	"github.com/Noxdew/Knights-Of-Discord/command"
	"github.com/Noxdew/Knights-Of-Discord/config"
	"github.com/Noxdew/Knights-Of-Discord/db"
	"github.com/Noxdew/Knights-Of-Discord/logger"
	"github.com/Noxdew/Knights-Of-Discord/structure"

	"github.com/bwmarrin/discordgo"
)

// ReadyHandler is called when `Ready` event is triggered
func ReadyHandler(s *discordgo.Session, r *discordgo.Ready) {
	s.UpdateStatus(0, "Knights of Discord")
	structure.DefaultServer.BuildServer()
	logger.Log.Info("Knights of Discord has successfully started.")
}

// ServerJoinHandler is called when `GuildCreate` event is triggered
func ServerJoinHandler(s *discordgo.Session, g *discordgo.GuildCreate) {
	server, err := db.GetServer(g.Guild.ID)
	if err != nil && err != db.NotFound {
		logger.Log.Error(err.Error())
	} else if err == db.NotFound {
		builder.BuildServer(server, s, g.Guild)
	} else {
		// Server exists
	}
}

// MessageReceiveHandler function called when Message is sent
func MessageReceiveHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore bot messages
	if m.Author.Bot {
		return
	}

	// Check for command
	if strings.HasPrefix(m.Content, config.Get().Prefix) {
		// Get Server object
		c, err := s.Channel(m.ChannelID)
		if err != nil {
			logger.Log.Error(err.Error())
			return
		}
		server, err := db.GetServer(c.GuildID)
		if err != nil {
			logger.Log.Error(err.Error())
			return
		}

		// Execute command
		trigger := strings.Split(m.Content, config.Get().Prefix)[1]
		for _, cmd := range command.MessageCommands {
			if cmd.Trigger() == trigger {
				cmd.Execute(server, s, m)
				return
			}
		}

		// Fallback for wrong command
		// Create response message
		message := &structure.Message{
			Title:  "Unknown command. Try using `!kod-help` for a list of all game commands.",
			Type:   "system",
			Icon:   "https://cdn.discordapp.com/attachments/512302843437252611/512302951814004752/ac6918be09a389876ee5663d6b08b55a.png",
			Footer: "Command execution feedback.",
		}

		// Build Embed
		embed := builder.BuildEmbed(message)

		// Send response
		_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)
		if err != nil {
			logger.Log.Error(err.Error())
		}
	}
}

// ReactionAddHandler function called when a Reaction is sent
func ReactionAddHandler(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	// Check author for bot
	u, err := s.User(r.UserID)
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}
	if u.Bot {
		return
	}

	// Get Server object
	c, err := s.Channel(r.ChannelID)
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}
	server, err := db.GetServer(c.GuildID)
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}

	// Call reaction command
	for _, cmd := range command.ReactionCommands {
		if cmd.Trigger() == r.Emoji.ID {
			cmd.Execute(server, s, r)
			return
		}
	}
}

// RoleEditHandler function called when a Discord Role receives an update
func RoleEditHandler(s *discordgo.Session, r *discordgo.GuildRoleUpdate) {
	// Get Server object
	server, err := db.GetServer(r.GuildID)
	if err != nil && err != db.NotFound {
		logger.Log.Error(err.Error())
		return
	} else if err == db.NotFound {
		return
	}

	// Check for playing Server
	if !server.Playing {
		return
	}

	// Check Discord Role
	for _, role := range server.Roles {
		if role.ID == r.Role.ID {
			// Revert to game requirement
			if r.Role.Hoist != role.Hoist || r.Role.Mentionable != role.Mentionable || r.Role.Permissions != server.RolePerm {
				_, err = s.GuildRoleEdit(server.ID, r.Role.ID, r.Role.Name, r.Role.Color, role.Hoist, server.RolePerm, role.Mentionable)
				if err != nil {
					logger.Log.Error(err.Error())
				}
			}
			return
		}
	}
}

// ChannelEditHandler function called  when a Discord Channel receives an update
func ChannelEditHandler(s *discordgo.Session, c *discordgo.ChannelUpdate) {
	// Get Server object
	server, err := db.GetServer(c.GuildID)
	if err != nil && err != db.NotFound {
		logger.Log.Error(err.Error())
		return
	} else if err == db.NotFound {
		return
	}

	// Check for playing Server
	if !server.Playing {
		return
	}

	// Check Discord Category
	if c.ID == server.Category.ID {
		if len(c.PermissionOverwrites) > 0 {
			_, err = s.ChannelEditComplex(c.ID, &discordgo.ChannelEdit{
				PermissionOverwrites: []*discordgo.PermissionOverwrite{},
			})
			if err != nil {
				logger.Log.Error(err.Error())
			}
		}
		return
	}

	// Check Discord Channel
	for _, channel := range server.Channels {
		if channel.ID == c.ID {
			// Revert to game requirement
			if c.Position != channel.Position || c.ParentID != server.Category.ID {
				_, err = s.ChannelEditComplex(c.ID, &discordgo.ChannelEdit{
					ParentID: server.Category.ID,
					Position: channel.Position,
				})
				if err != nil {
					logger.Log.Error(err.Error())
				}
			}

			// Check Permissions
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
			return
		}
	}
}
