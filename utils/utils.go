package utils

import "github.com/bwmarrin/discordgo"

// GetChannelByName finds channel with name n in guild g
func GetChannelByName(g *discordgo.Guild, n string) *discordgo.Channel {
	for _, ch := range g.Channels {
		if ch.Name == n {
			return ch
		}
	}
	return nil
}

// GetChannelInCategory finds channel with name n nested in category p in guild g
func GetChannelInCategory(g *discordgo.Guild, p *discordgo.Channel, n string) *discordgo.Channel {
	for _, ch := range g.Channels {
		if ch.ParentID == p.ID && ch.Name == n {
			return ch
		}
	}
	return nil
}

// GetRoleByName finds role with name n in guild g
func GetRoleByName(g *discordgo.Guild, n string) *discordgo.Role {
	for _, r := range g.Roles {
		if r.Name == n {
			return r
		}
	}
	return nil
}

// GetRoleByID finds role with name n in guild g
func GetRoleByID(g *discordgo.Guild, n string) *discordgo.Role {
	for _, r := range g.Roles {
		if r.ID == n {
			return r
		}
	}
	return nil
}

// HasRole checks if member m has a role with name n
func HasRole(m *discordgo.Member, n string) bool {
	for _, role := range m.Roles {
		if role == n {
			return true
		}
	}
	return false
}
