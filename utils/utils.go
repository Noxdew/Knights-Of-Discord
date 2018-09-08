package utils

import (
	"github.com/bwmarrin/discordgo"
)

func GetChannelByName(g *discordgo.Guild, n string) *discordgo.Channel {
	for _, ch := range g.Channels {
		if ch.Name == n {
			return ch
		}
	}
	return nil
}

func GetChannelInCategory(g *discordgo.Guild, p *discordgo.Channel, n string) *discordgo.Channel {
	for _, ch := range g.Channels {
		if ch.ParentID == p.ID && ch.Name == n {
			return ch
		}
	}
	return nil
}

func GetRoleByName(g *discordgo.Guild, n string) *discordgo.Role {
	for _, r := range g.Roles {
		if r.Name == n {
			return r
		}
	}
	return nil
}

func HasRole(m *discordgo.Member, n string) bool {
	for _, role := range m.Roles {
		if role == n {
			return true
		}
	}
	return false
}