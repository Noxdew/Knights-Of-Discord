package command

import (
	"github.com/Noxdew/Knights-Of-Discord/builder"
	"github.com/Noxdew/Knights-Of-Discord/db"
	"github.com/Noxdew/Knights-Of-Discord/logger"
	"github.com/Noxdew/Knights-Of-Discord/structure"
	"github.com/bwmarrin/discordgo"
)

// Action interface for parsing reaction commands
type Action interface {
	Trigger() string
	Execute(*structure.Server, *discordgo.Session, *discordgo.MessageReactionAdd)
}

// Response interface for parsing message commands
type Response interface {
	Trigger() string
	Description() string
	Execute(*structure.Server, *discordgo.Session, *discordgo.MessageCreate)
}

// AddUser command structure
type AddUser struct{}

// Execute method for AddUser command
func (*AddUser) Execute(server *structure.Server, s *discordgo.Session, m *discordgo.MessageReactionAdd) {
	// Check for right message
	if server.Messages["rules"].ID != m.MessageID {
		return
	}

	// Check is User exists
	for _, user := range server.Users {
		if user.ID == m.UserID {
			return
		}
	}

	// Assign Game Role to User
	err := s.GuildMemberRoleAdd(server.ID, m.UserID, server.Roles["r1"].ID)
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}

	// Update Server object
	err = db.AddServerUser(server, &structure.User{
		ID:           m.UserID,
		Role:         server.Roles["r1"].ID,
		Contribution: 0,
	})
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}

	// Send Message to the user
	channel, err := s.UserChannelCreate(m.UserID)
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}
	s.ChannelMessageSend(channel.ID, "Welcome to Knights of Discord!\nYou have joined a new Guild!")
}

// Trigger for AddUser command
func (*AddUser) Trigger() string {
	return "514099648949125153"
}

// ReactionCommands array
var ReactionCommands = []Action{
	&AddUser{},
}

// CloseGame command
type CloseGame struct{}

// Execute method for CloseGame command
func (*CloseGame) Execute(server *structure.Server, s *discordgo.Session, m *discordgo.MessageCreate) {
	// Check if Guild admin issued the command
	g, err := s.Guild(server.ID)
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}
	if g.OwnerID != m.Author.ID {
		return
	}

	// Create response message
	message := &structure.Message{
		Title:  "Closing game...",
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
	builder.DestroyServer(server, s, g)
}

// Trigger for CloseGame command
func (*CloseGame) Trigger() string {
	return "closeGame"
}

// Description for CloseGame command
func (*CloseGame) Description() string {
	return "Close the game running on this server.\nOnly the server's owner can execute this command.\n**WARNING!** Closing the game will result in losing all progress and resources of your server!\n"
}

// Help command
type Help struct{}

// Execute method for Help command
func (*Help) Execute(server *structure.Server, s *discordgo.Session, m *discordgo.MessageCreate) {
	// Create response message
	message := &structure.Message{
		Title:  "Command List",
		Type:   "system",
		Icon:   "https://cdn.discordapp.com/attachments/512302843437252611/512302951814004752/ac6918be09a389876ee5663d6b08b55a.png",
		Footer: "Command execution feedback.",
		Fields: []*structure.Field{},
	}

	// Add fields
	for _, cmd := range MessageCommands {
		message.Fields = append(message.Fields, &structure.Field{
			Title: "`!kod-" + cmd.Trigger() + "`",
			Value: cmd.Description(),
		})
	}

	// Build Embed
	embed := builder.BuildEmbed(message)

	// Send response
	_, err := s.ChannelMessageSendEmbed(m.ChannelID, embed)
	if err != nil {
		logger.Log.Error(err.Error())
	}
}

// Trigger for Help command
func (*Help) Trigger() string {
	return "help"
}

// Description for Help command
func (*Help) Description() string {
	return "Request a list of all game commands.\n"
}

// LeaveServer command
type LeaveServer struct{}

// Execute method for LeaveServer command
func (*LeaveServer) Execute(server *structure.Server, s *discordgo.Session, m *discordgo.MessageCreate) {
	// Check if user is playing
	for _, user := range server.Users {
		if user.ID == m.Author.ID {
			err := db.RemoveServerUser(server, user)
			if err != nil {
				logger.Log.Error(err.Error())
				return
			}

			// Create response message
			message := &structure.Message{
				Title:  "User " + m.Author.Username + " successfully removed from the server",
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
			return
		}
	}
}

// Trigger for LeaveServer command
func (*LeaveServer) Trigger() string {
	return "leaveGame"
}

// Description for LeaveServer command
func (*LeaveServer) Description() string {
	return "Removes the user from the game running on this server.\n**WARNING** Leaving the game will cause you to lose all your progression in the server!\n"
}

// MessageCommands array
var MessageCommands = []Response{
	&CloseGame{},
	&Help{},
	&LeaveServer{},
}
