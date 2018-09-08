package bot

import (
	"github.com/noxdew/knights-of-discord/config"
	"github.com/noxdew/knights-of-discord/utils"
	"github.com/noxdew/knights-of-discord/builder"
	"fmt"
	"github.com/bwmarrin/discordgo"
)

var BotID string
var goBot *discordgo.Session

func Start() {
	goBot, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	u, err := goBot.User("@me")
	if err != nil {
		fmt.Println(err.Error())
	}

	BotID = u.ID

	goBot.AddHandler(messageHandler)
	goBot.AddHandler(serverJoin)

	err = goBot.Open()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Bot is running")
	<-make(chan struct{})
}

func serverJoin(s *discordgo.Session, g *discordgo.GuildCreate) {
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

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Content == "test" {
		_, _ = s.ChannelMessageSend(m.ChannelID, "yo mom")
	}
}