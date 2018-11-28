package command

import (
	"github.com/Noxdew/Knights-Of-Discord/builder"
	"github.com/Noxdew/Knights-Of-Discord/structure"
	"github.com/bwmarrin/discordgo"
)

// Action interface for parsing reaction commands
type Action interface {
	Trigger() string
	Execute(*structure.Server, *discordgo.Session, *discordgo.Guild, *discordgo.MessageReactionAdd)
}

// Response interface for parsing message commands
type Response interface {
	Trigger() string
	Execute(*structure.Server, *discordgo.Session, *discordgo.Guild, *discordgo.MessageCreate)
}

// AddUser command structure
type AddUser struct{}

// Execute AddUser command
func (*AddUser) Execute(server *structure.Server, s *discordgo.Session, g *discordgo.Guild, m *discordgo.MessageReactionAdd) {
	if server.Messages.Rules.ID != m.MessageID {
		return
	}

	_, err := server.FindUser(m.UserID)
	if err != nil {
		builder.AddUser(server, s, g, m.UserID)
	}
}

// Trigger for AddUser command
func (*AddUser) Trigger() string {
	return "514099648949125153"
}

// CloseGame command structure
type CloseGame struct{}

// Execute CloseGame command
func (*CloseGame) Execute(server *structure.Server, s *discordgo.Session, g *discordgo.Guild, m *discordgo.MessageCreate) {
	s.ChannelMessageSend(m.ChannelID, "Closing Game...")
	builder.DestroyServer(server, s, g)
}

// Trigger for CloseGame command
func (*CloseGame) Trigger() string {
	return "closeGame"
}

// ReactionCommands array
var ReactionCommands = []Action{
	&AddUser{},
}

// MessageCommands array
var MessageCommands = []Response{
	&CloseGame{},
}
