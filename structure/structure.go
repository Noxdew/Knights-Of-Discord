package structure

import (
	"encoding/json"
	"errors"
	"io/ioutil"

	"github.com/Noxdew/Knights-Of-Discord/logger"
	"github.com/bwmarrin/discordgo"
)

// Server contains game information for a Discord Guild
type Server struct {
	ID       string   `json:"id" bson:"id"`
	Playing  bool     `json:"playing" bson:"playing"`
	Power    int      `json:"power" bson:"power"`
	Roles    Roles    `json:"roles" bson:"roles"`
	Category Category `json:"category" bson:"category"`
	BotPerm  int      `json:"botPerm" bson:"botPerm"`
	Channels Channels `json:"channels" bson:"channels"`
	Messages Messages `json:"messages" bson:"messages"`
	Users    []User   `json:"users" bson:"users"`
}

// Roles contains game Roles
type Roles struct {
	Villager Role `json:"villager" bson:"villager"`
	Esquire  Role `json:"esquire" bson:"esquire"`
	Knight   Role `json:"knight" bson:"knight"`
	Everyone Role `json:"everyone" bson:"everyone"`
}

// Role contains game information for a Discord Role
type Role struct {
	ID         string `json:"id" bson:"id"`
	DefName    string `json:"defName" bson:"defName"`
	Permission int    `json:"permission" bson:"permission"`
	Level      int    `json:"level" bson:"level"`
}

// Category contains game information for a Discord Category Channel
type Category struct {
	ID      string `json:"id" bson:"id"`
	DefName string `json:"defName" bson:"defName"`
}

// Channels contains game Channels
type Channels struct {
	Rules         Channel `json:"rules" bson:"rules"`
	Announcements Channel `json:"announcements" bson:"announcements"`
	Outskirts     Channel `json:"outskirts" bson:"outskirts"`
	Tavern        Channel `json:"tavern" bson:"tavern"`
	InnerCity     Channel `json:"innerCity" bson:"innerCity"`
	Inn           Channel `json:"inn" bson:"inn"`
	Castle        Channel `json:"castle" bson:"castle"`
	MeadHall      Channel `json:"meadHall" bson:"meadHall"`
}

// Channel contains game information for a Discord Channel
type Channel struct {
	ID          string `json:"id" bson:"id"`
	DefName     string `json:"defName" bson:"defName"`
	Level       int    `json:"level" bson:"level"`
	Allow       int    `json:"allow" bson:"allow"`
	Deny        int    `json:"deny" bson:"deny"`
	Permissions []Perm `json:"perms" bson:"perms"`
}

// Perm conatins game information for a Discord PermissionOverwrite
type Perm struct {
	Role  string `json:"role" bson:"role"`
	Allow int    `json:"allow" bson:"allow"`
	Deny  int    `json:"deny" bson:"deny"`
}

// Messages contains game Messages
type Messages struct {
	Rules   Message `json:"rules,omitempty" bson:"rules,omitempty"`
	Farm    Message `json:"farm,omitempty" bson:"farm,omitempty"`
	Woods   Message `json:"woods,omitempty" bson:"woods,omitempty"`
	Quary   Message `json:"quary,omitempty" bson:"quary,omitempty"`
	Builder Message `json:"builder,omitempty" bson:"builder,omitempty"`
}

// Message contains game information for a Discord Message
type Message struct {
	ID          string  `json:"id" bson:"id"`
	Channel     string  `json:"channel" bson:"channel"`
	Title       string  `json:"title" bson:"title"`
	Type        string  `json:"type" bson:"type"`
	Description string  `json:"description" bson:"description"`
	Color       int     `json:"color" bson:"color"`
	Icon        string  `json:"icon" bson:"icon"`
	Footer      string  `json:"footer" bson:"footer"`
	Fields      []Field `json:"fields" bson:"fields"`
}

// Field contains game a Discord MessageEmbed Field
type Field struct {
	Title  string `json:"title" bson:"title"`
	Value  string `json:"value" bson:"value"`
	Inline bool   `json:"fInline" bson:"fInline"`
}

// User contains game information for a Discord User
type User struct {
	ID           string `json:"id" bson:"id"`
	Role         string `json:"role" bson:"role"`
	Contribution int    `json:"contribution" bson:"contribution"`
}

// BuildServer creates a new Server object to store Discord Guild information
func (s *Server) BuildServer(g *discordgo.Guild) {
	file, err := ioutil.ReadFile("structure.json")
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}
	err = json.Unmarshal(file, &s)
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}
	s.ID = g.ID
}

// FindUser returns User object if one exist
func (s *Server) FindUser(u string) (*User, error) {
	for _, user := range s.Users {
		if user.ID == u {
			return &user, nil
		}
	}
	return &User{}, errors.New("user doesn't exist")
}

// GetRoles returns all game Roles as a slice
func (r *Roles) GetRoles() []*Role {
	var arr []*Role
	arr = append(arr, &r.Villager)
	arr = append(arr, &r.Esquire)
	arr = append(arr, &r.Knight)
	arr = append(arr, &r.Everyone)
	return arr
}

// GetChannels returns all game Channels as a slice
func (c *Channels) GetChannels() []*Channel {
	var arr []*Channel
	arr = append(arr, &c.Rules)
	arr = append(arr, &c.Announcements)
	arr = append(arr, &c.Outskirts)
	arr = append(arr, &c.Tavern)
	arr = append(arr, &c.InnerCity)
	arr = append(arr, &c.Inn)
	arr = append(arr, &c.Castle)
	arr = append(arr, &c.MeadHall)
	return arr
}

// SetMessageChannel assigns a Channel ID to game Messages
func (s *Server) SetMessageChannel() {
	s.Messages.Rules.Channel = s.Channels.Rules.ID
	s.Messages.Farm.Channel = s.Channels.Outskirts.ID
	s.Messages.Woods.Channel = s.Channels.Outskirts.ID
	s.Messages.Quary.Channel = s.Channels.Outskirts.ID
	s.Messages.Builder.Channel = s.Channels.Outskirts.ID
}

// GetMessages returns all game Messages as a slice
func (m *Messages) GetMessages() []*Message {
	var arr []*Message
	arr = append(arr, &m.Rules)
	arr = append(arr, &m.Farm)
	arr = append(arr, &m.Woods)
	arr = append(arr, &m.Quary)
	arr = append(arr, &m.Builder)
	return arr
}
