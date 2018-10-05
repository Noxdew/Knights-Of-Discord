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
	ID      string `json:"id" bson:"id"`
	Power   int    `json:"power" bson:"power"`
	Checked bool   `json:"checked" bson:"checked"`
	Roles   []Role `json:"roles" bson:"roles"`
}

// Role represents a Discord Role in the game
type Role struct {
	ID   string `json:"id" bson:"id"`
	Type string `json:"type" bson:"type"`
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

// CreateServer uploads a new server to the DB
func CreateServer(s Server) error {
	client := connect()
	defer client.Disconnect(context.Background())
	collection := client.Database("knights-of-discord").Collection("servers")
	_, err := collection.InsertOne(context.Background(), s)
	return err
}

// RemoveServer removes a server from the DB
func RemoveServer(s string) error {
	client := connect()
	defer client.Disconnect(context.Background())
	collection := client.Database("knights-of-discord").Collection("servers")
	filter := bson.NewDocument(bson.EC.String("id", s))
	_, err := collection.DeleteOne(context.Background(), filter)
	return err
}

// FlagServers changes the `checked` value of all DB servers to `b`
func FlagServers(b bool) error {
	client := connect()
	defer client.Disconnect(context.Background())
	collection := client.Database("knights-of-discord").Collection("servers")
	filter := bson.NewDocument()
	replacement := bson.NewDocument(bson.EC.SubDocumentFromElements("$set", bson.EC.Boolean("checked", b)))
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
	filter := bson.NewDocument(bson.EC.String("id", s), bson.EC.String("roles.type", r.Type))
	replacement := bson.NewDocument(bson.EC.SubDocumentFromElements("$set", bson.EC.Interface("roles.$.id", r.ID)))
	_, err := collection.UpdateOne(context.Background(), filter, replacement)
	return err
}
