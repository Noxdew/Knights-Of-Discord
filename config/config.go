package config

import (
	"encoding/json"
	"io/ioutil"

	"github.com/Noxdew/Knights-Of-Discord/logger"
)

// Definition defines the structure of a config for this application
type Definition struct {
	Token      string          `json:"Token"`
	Prefix     string          `json:"Prefix"`
	DBUrl      string          `json:"DBUrl"`
	DBUser     string          `json:"DBUser"`
	DBPassword string          `json:"DBPassword"`
	Roles      []string        `json:"Roles"`
	RolePerm   int             `json:"RolePerm"`
	Channels   []ChannelConfig `json:"Channels"`
	BotPerm    PermConfig      `json:"BotPerm"`
	ActionPerm PermConfig      `json:"ActionPerm"`
	SocialPerm PermConfig      `json:"SocialPerm"`
	ClosedPerm PermConfig      `json:"ClosedPerm"`
}

// ChannelConfig defines config structure of a channel
type ChannelConfig struct {
	Name string   `json:"name"`
	Type string   `json:"type"`
	Role []string `json:"role"`
}

// PermConfig defines config structure of a channel
type PermConfig struct {
	Allow int `json:"allow"`
	Deny  int `json:"deny"`
}

// Config contains the configuration of this application
var config *Definition

// Get returns the config struct
func Get() *Definition {
	return config
}

// Load reads the config json file and parses it, then stores the values for the rest of the application to use
func Load() {
	logger.Log.Info("Reading config file...")
	file, err := ioutil.ReadFile("config.json")
	if err != nil {
		logger.Log.Panic(err.Error())
	}

	// Parse the json
	err = json.Unmarshal(file, &config)
	if err != nil {
		logger.Log.Panic(err.Error())
	}

	logger.Log.Info("Config loaded")
}
