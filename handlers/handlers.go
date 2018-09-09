package handlers

import (
	"github.com/Noxdew/Knights-Of-Discord/builder"
	"github.com/Noxdew/Knights-Of-Discord/logger"
	"github.com/Noxdew/Knights-Of-Discord/utils"
	"github.com/bwmarrin/discordgo"
)

// ServerJoinHandler function called when a server is joined
func ServerJoinHandler(s *discordgo.Session, g *discordgo.GuildCreate) {
	// Check for existing game
	channel := utils.GetChannelByName(g.Guild, "Knights of Discord")
	if channel != nil {
		logger.Log.Info("Checking game integrity in server %s (%s)", g.Guild.Name, g.Guild.ID)
		builder.BuildRoles(s, g.Guild)
		builder.BuildChannels(s, g.Guild, channel)
	} else {
		logger.Log.Info("Starting new game for %s (%s)", g.Guild.Name, g.Guild.ID)
		channel, err := s.GuildChannelCreate(g.Guild.ID, "Knights of Discord", "4")
		if err != nil {
			logger.Log.Error(err.Error())
		} else {
			builder.BuildRoles(s, g.Guild)
			builder.BuildChannels(s, g.Guild, channel)
		}
	}
	logger.Log.Infof("Game started for %s (%s)", g.Guild.Name, g.Guild.ID)
}
