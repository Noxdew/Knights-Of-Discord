package db

import (
	"context"

	"github.com/Noxdew/Knights-Of-Discord/config"
	"github.com/Noxdew/Knights-Of-Discord/logger"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
)

// NotFound represents empty query results
var NotFound = mongo.ErrNoDocuments

// Server represents a Discord Guild in the game
type Server struct {
	ID       string    `json:"id" bson:"id"`
	Checked  bool      `json:"checked" bson:"checked"`
	Playing  bool      `json:"playing" bson:"playing"`
	Power    int       `json:"power" bson:"power"`
	Roles    []Role    `json:"roles" bson:"roles"`
	Category string    `json:"category" bson:"category"`
	Channels []Channel `json:"channels" bson:"channels"`
}

// Role represents a Discord Role in the game
type Role struct {
	ID      string `json:"id" bson:"id"`
	DefName string `json:"defName" bson:"defName"`
}

// Channel represents a Discord Channel in the game
type Channel struct {
	ID      string   `json:"id" bson:"id"`
	DefName string   `json:"defName" bson:"defName"`
	Type    string   `json:"type" bson:"type"`
	Roles   []string `json:"roles" bson:"roles"`
	Perms   []Perm   `json:"perms" bson:"perms"`
}

// Perm represents a Discord PermissionOverwrite in the game
type Perm struct {
	ID    string `json:"id" bson:"id"`
	Type  string `json:"type" bson:"type"`
	Allow int    `json:"allow" bson:"allow"`
	Deny  int    `json:"deny" bson:"deny"`
}

func connect() *mongo.Client {
	client, err := mongo.NewClient("mongodb://" + config.Get().DBUser + ":" + config.Get().DBPassword + "@" + config.Get().DBUrl)
	if err != nil {
		logger.Log.Panic(err)
	}
	err = client.Connect(context.Background())
	if err != nil {
		logger.Log.Panic(err)
	}
	return client
}

// GetServer returns a Guild from the DB with ID s
func GetServer(s string) (*Server, error) {
	client := connect()
	defer client.Disconnect(context.Background())
	collection := client.Database("knights-of-discord").Collection("servers")
	filter := bson.NewDocument(bson.EC.String("id", s))
	server := Server{}
	doc := collection.FindOne(context.Background(), filter)
	err := doc.Decode(&server)
	return &server, err
}

// CreateServer uploads a new server `s` to the DB
func CreateServer(s Server) error {
	client := connect()
	defer client.Disconnect(context.Background())
	collection := client.Database("knights-of-discord").Collection("servers")
	_, err := collection.InsertOne(context.Background(), s)
	return err
}

// RemoveServer removes server with ID `s` from the DB
func RemoveServer(s string) error {
	client := connect()
	defer client.Disconnect(context.Background())
	collection := client.Database("knights-of-discord").Collection("servers")
	filter := bson.NewDocument(bson.EC.String("id", s))
	_, err := collection.DeleteOne(context.Background(), filter)
	return err
}

// UpdateServerChecked changes the `checked` value of server with ID `s` to `b`
func UpdateServerChecked(s string, b bool) error {
	client := connect()
	defer client.Disconnect(context.Background())
	collection := client.Database("knights-of-discord").Collection("servers")
	filter := bson.NewDocument(bson.EC.String("id", s))
	replacement := bson.NewDocument(bson.EC.SubDocumentFromElements("$set", bson.EC.Boolean("checked", b)))
	_, err := collection.UpdateMany(context.Background(), filter, replacement)
	return err
}

// UpdateServerPlaying changes the `playing` value for server with ID `s` to `b`
func UpdateServerPlaying(s string, b bool) error {
	client := connect()
	defer client.Disconnect(context.Background())
	collection := client.Database("knights-of-discord").Collection("servers")
	filter := bson.NewDocument(bson.EC.String("id", s))
	replacement := bson.NewDocument(bson.EC.SubDocumentFromElements("$set", bson.EC.Boolean("playing", b)))
	_, err := collection.UpdateMany(context.Background(), filter, replacement)
	return err
}

// CreateRole uploads new role `r` to server `s`
func CreateRole(r Role, s string) error {
	client := connect()
	defer client.Disconnect(context.Background())
	collection := client.Database("knights-of-discord").Collection("servers")
	filter := bson.NewDocument(bson.EC.String("id", s))
	replacement := bson.NewDocument(bson.EC.SubDocumentFromElements("$push", bson.EC.Interface("roles", r)))
	_, err := collection.UpdateOne(context.Background(), filter, replacement)
	return err
}

// UpdateRole updates an existing role in server `s`
func UpdateRole(r Role, s string) error {
	client := connect()
	defer client.Disconnect(context.Background())
	collection := client.Database("knights-of-discord").Collection("servers")
	filter := bson.NewDocument(bson.EC.String("id", s), bson.EC.String("roles.defName", r.DefName))
	replacement := bson.NewDocument(bson.EC.SubDocumentFromElements("$set", bson.EC.Interface("roles.$.id", r.ID)))
	_, err := collection.UpdateOne(context.Background(), filter, replacement)
	return err
}

// CreateCategory uploads a new category `c` to server `s`
func CreateCategory(c string, s string) error {
	client := connect()
	defer client.Disconnect(context.Background())
	collection := client.Database("knights-of-discord").Collection("servers")
	filter := bson.NewDocument(bson.EC.String("id", s))
	replacement := bson.NewDocument(bson.EC.SubDocumentFromElements("$set", bson.EC.String("category", c)))
	_, err := collection.UpdateOne(context.Background(), filter, replacement)
	return err
}

// CreateChannel uploads new channel `c` to server `s`
func CreateChannel(c Channel, s string) error {
	client := connect()
	defer client.Disconnect(context.Background())
	collection := client.Database("knights-of-discord").Collection("servers")
	filter := bson.NewDocument(bson.EC.String("id", s))
	replacement := bson.NewDocument(bson.EC.SubDocumentFromElements("$push", bson.EC.Interface("channels", c)))
	_, err := collection.UpdateOne(context.Background(), filter, replacement)
	return err
}

// UpdateChannel updates an existing channel in server `s`
func UpdateChannel(c Channel, s string) error {
	client := connect()
	defer client.Disconnect(context.Background())
	collection := client.Database("knights-of-discord").Collection("servers")
	filter := bson.NewDocument(bson.EC.String("id", s), bson.EC.String("channels.defName", c.DefName))
	replacement := bson.NewDocument(bson.EC.SubDocumentFromElements("$set", bson.EC.Interface("channels.$.id", c.ID)))
	_, err := collection.UpdateOne(context.Background(), filter, replacement)
	return err
}

// CreatePermission uploads a new permission `p` to channel `c` in server `s`
func CreatePermission(p Perm, c Channel, s string) error {
	client := connect()
	defer client.Disconnect(context.Background())
	collection := client.Database("knights-of-discord").Collection("servers")
	filter := bson.NewDocument(bson.EC.String("id", s), bson.EC.String("channels.id", c.ID))
	replacement := bson.NewDocument(bson.EC.SubDocumentFromElements("$push", bson.EC.Interface("channels.$.perms", p)))
	_, err := collection.UpdateOne(context.Background(), filter, replacement)
	return err
}
