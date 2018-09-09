package bot

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/noxdew/knights-of-discord/builder"
	"github.com/noxdew/knights-of-discord/config"
	"github.com/noxdew/knights-of-discord/utils"
)

var goBot *discordgo.Session

func Start() {
	// Login the bot client and save its user ID
	goBot, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Add event handlers
	goBot.AddHandler(messageHandler)
	goBot.AddHandler(serverHandler)

	// Start the bot's session
	err = goBot.Open()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("Bot is running")

	<-make(chan struct{})
}

// Handler function called when a server is joined
func serverHandler(s *discordgo.Session, g *discordgo.GuildCreate) {
	// Check for existing game
	channel := utils.GetChannelByName(g.Guild, "Knights of Discord")
	if channel != nil {
		fmt.Println("Checking game integrity...")
		builder.BuildRoles(s, g.Guild)
		builder.BuildChannels(s, g.Guild, channel)
	} else {
		fmt.Println("Starting new game on the server...")
		channel, err := s.GuildChannelCreate(g.Guild.ID, "Knights of Discord", "4")
		if err != nil {
			fmt.Println(err.Error())
		} else {
			builder.BuildRoles(s, g.Guild)
			builder.BuildChannels(s, g.Guild, channel)
		}
	}
	fmt.Println("Game started.")
}

// Handler function called when a message is received
func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
		return
	}

	if m.Content == "ping" {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}
}
