package builder

import (
	"github.com/bwmarrin/discordgo"
	"github.com/noxdew/knights-of-discord/utils"
	"fmt"
)

func BuildRoles(s *discordgo.Session, g *discordgo.Guild) {
	fmt.Println("Building roles...")
	role := utils.GetRoleByName(g, "KoD-King")
	if role == nil {
		role, err := s.GuildRoleCreate(g.ID)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			baseRole := utils.GetRoleByName(g, "@everyone")
			perm := baseRole.Permissions
			s.GuildRoleEdit(g.ID, role.ID, "KoD-King", 0, false, perm, true)
		}
	}
	king, err := s.GuildMember(g.ID, g.OwnerID)
	if err != nil {
		fmt.Println(err.Error())
	}
	if !utils.HasRole(king, role.Name) {
		s.GuildMemberRoleAdd(g.ID, king.User.ID, role.ID)
	}

	role = utils.GetRoleByName(g, "KoD-Knight")
	if role == nil {
		role, err := s.GuildRoleCreate(g.ID)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			baseRole := utils.GetRoleByName(g, "@everyone")
			perm := baseRole.Permissions
			s.GuildRoleEdit(g.ID, role.ID, "KoD-Knight", 0, false, perm, true)
		}
	}
	role = utils.GetRoleByName(g, "KoD-Esquire")
	if role == nil {
		role, err := s.GuildRoleCreate(g.ID)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			baseRole := utils.GetRoleByName(g, "@everyone")
			perm := baseRole.Permissions
			s.GuildRoleEdit(g.ID, role.ID, "KoD-Esquire", 0, false, perm, true)
		}
	}
	role = utils.GetRoleByName(g, "KoD-Villager")
	if role == nil {
		role, err := s.GuildRoleCreate(g.ID)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			baseRole := utils.GetRoleByName(g, "@everyone")
			perm := baseRole.Permissions
			s.GuildRoleEdit(g.ID, role.ID, "KoD-Villager", 0, false, perm, true)
		}
	}
	fmt.Println("Roles built.")
}

func BuildChannels(s *discordgo.Session, g *discordgo.Guild, c *discordgo.Channel) {
	fmt.Println("Building channels...")
	ch := utils.GetChannelInCategory(g, c, "announcements")
	if ch == nil {
		ch, err := s.GuildChannelCreate(g.ID, "announcements", "0")
		if err != nil {
			fmt.Println(err.Error())
		}
		s.ChannelEditComplex(ch.ID, &discordgo.ChannelEdit {
			Topic: "General Castle Information",
			Position: 0,
			ParentID: c.ID,
		})
	}
	ch = utils.GetChannelInCategory(g, c, "logs")
	if ch == nil {
		ch, err := s.GuildChannelCreate(g.ID, "logs", "0")
		if err != nil {
			fmt.Println(err.Error())
		}
		s.ChannelEditComplex(ch.ID, &discordgo.ChannelEdit {
			Topic: "Activity Log For The Castle",
			Position: 1,
			ParentID: c.ID,
		})
	}
	ch = utils.GetChannelInCategory(g, c, "small-tavern")
	if ch == nil {
		ch, err := s.GuildChannelCreate(g.ID, "small-tavern", "0")
		if err != nil {
			fmt.Println(err.Error())
		}
		s.ChannelEditComplex(ch.ID, &discordgo.ChannelEdit {
			Topic: "General Social Space For Villagers",
			Position: 2,
			ParentID: c.ID,
		})
	}
	ch = utils.GetChannelInCategory(g, c, "medium-tavern")
	if ch == nil {
		ch, err := s.GuildChannelCreate(g.ID, "medium-tavern", "0")
		if err != nil {
			fmt.Println(err.Error())
		}
		s.ChannelEditComplex(ch.ID, &discordgo.ChannelEdit {
			Topic: "General Social Space For Esquires",
			Position: 3,
			ParentID: c.ID,
		})
	}
	ch = utils.GetChannelInCategory(g, c, "large-tavern")
	if ch == nil {
		ch, err := s.GuildChannelCreate(g.ID, "large-tavern", "0")
		if err != nil {
			fmt.Println(err.Error())
		}
		s.ChannelEditComplex(ch.ID, &discordgo.ChannelEdit {
			Topic: "General Social Space For Knights",
			Position: 4,
			ParentID: c.ID,
		})
	}
	ch = utils.GetChannelInCategory(g, c, "outer-city")
	if ch == nil {
		ch, err := s.GuildChannelCreate(g.ID, "outer-city", "0")
		if err != nil {
			fmt.Println(err.Error())
		}
		s.ChannelEditComplex(ch.ID, &discordgo.ChannelEdit {
			Topic: "Activity Center For Villagers",
			Position: 5,
			ParentID: c.ID,
		})
	}
	ch = utils.GetChannelInCategory(g, c, "inner-city")
	if ch == nil {
		ch, err := s.GuildChannelCreate(g.ID, "inner-city", "0")
		if err != nil {
			fmt.Println(err.Error())
		}
		s.ChannelEditComplex(ch.ID, &discordgo.ChannelEdit {
			Topic: "Activity Center For Esquires",
			Position: 6,
			ParentID: c.ID,
		})
	}
	ch = utils.GetChannelInCategory(g, c, "castle")
	if ch == nil {
		ch, err := s.GuildChannelCreate(g.ID, "castle", "0")
		if err != nil {
			fmt.Println(err.Error())
		}
		s.ChannelEditComplex(ch.ID, &discordgo.ChannelEdit {
			Topic: "Activity Center For Knights",
			Position: 7,
			ParentID: c.ID,
		})
	}
	fmt.Println("Channels built.")
}