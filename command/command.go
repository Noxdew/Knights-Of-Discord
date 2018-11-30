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
	g, err := s.Guild(server.ID)
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}
	s.ChannelMessageSend(m.ChannelID, "Closing Game...")
	builder.DestroyServer(server, s, g)
}

// Trigger for CloseGame command
func (*CloseGame) Trigger() string {
	return "closeGame"
}

// MessageCommands array
var MessageCommands = []Response{
	&CloseGame{},
}
