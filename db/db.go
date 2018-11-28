package db

import (
	"context"

	"github.com/Noxdew/Knights-Of-Discord/config"
	"github.com/Noxdew/Knights-Of-Discord/logger"
	"github.com/Noxdew/Knights-Of-Discord/structure"
	"github.com/bwmarrin/discordgo"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
)

// NotFound represents empty query results
var NotFound = mongo.ErrNoDocuments

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

// GetServer returns a Server object for the Discord Guild
func GetServer(g *discordgo.Guild) (*structure.Server, error) {
	client := connect()
	defer client.Disconnect(context.Background())
	collection := client.Database("knights-of-discord").Collection("servers")
	filter := bson.NewDocument(bson.EC.String("id", g.ID))
	server := structure.Server{}
	doc := collection.FindOne(context.Background(), filter)
	err := doc.Decode(&server)
	return &server, err
}

// CreateServer uploads a Server object
func CreateServer(s *structure.Server) error {
	client := connect()
	defer client.Disconnect(context.Background())
	collection := client.Database("knights-of-discord").Collection("servers")
	_, err := collection.InsertOne(context.Background(), s)
	return err
}

// UpdateServerPlaying runs or stops a game for given Server
func UpdateServerPlaying(s *structure.Server, b bool) error {
	client := connect()
	defer client.Disconnect(context.Background())
	collection := client.Database("knights-of-discord").Collection("servers")
	filter := bson.NewDocument(bson.EC.String("id", s.ID))
	replacement := bson.NewDocument(bson.EC.SubDocumentFromElements("$set", bson.EC.Boolean("playing", b)))
	_, err := collection.UpdateOne(context.Background(), filter, replacement)
	return err
}

// DeleteServer removes a Server object
func DeleteServer(s *structure.Server) error {
	client := connect()
	defer client.Disconnect(context.Background())
	collection := client.Database("knights-of-discord").Collection("servers")
	filter := bson.NewDocument(bson.EC.String("id", s.ID))
	_, err := collection.DeleteOne(context.Background(), filter)
	return err
}

// AddUser uploads a User object
func AddUser(s *structure.Server, u *structure.User) error {
	client := connect()
	defer client.Disconnect(context.Background())
	collection := client.Database("knights-of-discord").Collection("servers")
	filter := bson.NewDocument(bson.EC.String("id", s.ID))
	replacement := bson.NewDocument(bson.EC.SubDocumentFromElements("$push", bson.EC.Interface("users", u)))
	_, err := collection.UpdateOne(context.Background(), filter, replacement)
	return err
}
