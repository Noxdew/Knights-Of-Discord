package structure

import (
	"encoding/json"
	"io/ioutil"

	"github.com/Noxdew/Knights-Of-Discord/logger"
)

// Server contains game information for a Discord Guild
type Server struct {
	ID           string               `json:"-" bson:"id"`
	Playing      bool                 `json:"-" bson:"playing"`
	Resources    map[string]*Resource `json:"resources" bson:"resources"`
	BotPerm      int                  `json:"botPerm" bson:"-"`
	SocialPerm   int                  `json:"socialPerm" bson:"-"`
	ActionPerm   int                  `json:"actionPerm" bson:"-"`
	RolePerm     int                  `json:"rolePerm" bson:"-"`
	EveryoneRole string               `json:"-" bson:"everyoneRole"`
	Roles        map[string]*Role     `json:"roles" bson:"roles"`
	Category     *Category            `json:"category" bson:"category"`
	Channels     map[string]*Channel  `json:"channels" bson:"channels"`
	Messages     map[string]*Message  `json:"messages" bson:"messages"`
	Actions      map[string]string    `json:"actions" bson:"-"`
	Users        []*User              `json:"-" bson:"users"`
}

// BuildServer creates a new Server object to store Discord Guild information
func (s *Server) BuildServer() {
	file, err := ioutil.ReadFile("structure.json")
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}
	err = json.Unmarshal(file, s)
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}
}

// Resource contains game information for a Server Resource
type Resource struct {
	Name  string `json:"name" bson:"-"`
	Count int    `json:"count" bson:"count"`
	Icon  string `json:"icon" bson:"-"`
}

// Role contains game information for a Discord Role
type Role struct {
	ID          string `json:"-" bson:"id"`
	DefaultName string `json:"defaultName" bson:"-"`
	Mentionable bool   `json:"mentionable" bson:"-"`
	Hoist       bool   `json:"hoist" bson:"-"`
	Tier        int    `json:"tier" bson:"-"`
}

// Category contains game information for a Discord Category Channel
type Category struct {
	ID          string `json:"-" bson:"id"`
	DefaultName string `json:"defaultName" bson:"-"`
}

// Channel contains game information for a Discord Channel
type Channel struct {
	ID          string `json:"-" bson:"id"`
	DefaultName string `json:"defaultName" bson:"-"`
	Topic       string `json:"topic" bson:"-"`
	Tier        int    `json:"tier" bson:"-"`
	Position    int    `json:"position" bson:"-"`
	Type        string `json:"type" bson:"-"`
}

// Message contains game information for a Discord Message
type Message struct {
	ID          string   `json:"-" bson:"id"`
	ChannelID   string   `json:"-" bson:"channelID"`
	Type        string   `json:"type" bson:"-"`
	Title       string   `json:"title" bson:"-"`
	Description string   `json:"description" bson:"-"`
	Icon        string   `json:"icon" bson:"-"`
	Footer      string   `json:"footer" bson:"-"`
	Fields      []*Field `json:"fields" bson:"-"`
}

// Field contains game a Discord MessageEmbed Field
type Field struct {
	Title string `json:"title" bson:"-"`
	Value string `json:"value" bson:"-"`
}

// User contains game information for a Discord User
type User struct {
	ID           string `json:"-" bson:"id"`
	Role         string `json:"-" bson:"role"`
	Contribution int    `json:"-" bson:"contribution"`
}

// DefaultServer object
var DefaultServer Server
