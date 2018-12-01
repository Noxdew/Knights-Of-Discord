package db

import (
	"context"

	"github.com/Noxdew/Knights-Of-Discord/config"
	"github.com/Noxdew/Knights-Of-Discord/logger"
	"github.com/Noxdew/Knights-Of-Discord/structure"
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
func GetServer(g string) (*structure.Server, error) {
	client := connect()
	defer client.Disconnect(context.Background())
	collection := client.Database("knights-of-discord").Collection("servers")
	filter := bson.NewDocument(bson.EC.String("id", g))
	server := structure.DefaultServer
	dbServer := structure.Server{}
	doc := collection.FindOne(context.Background(), filter)
	err := doc.Decode(&dbServer)

	if err == nil {
		server.ID = dbServer.ID
		server.Playing = dbServer.Playing
		for key, resource := range server.Resources {
			resource.Count = dbServer.Resources[key].Count
		}
		server.EveryoneRole = dbServer.EveryoneRole
		for key, role := range server.Roles {
			role.ID = dbServer.Roles[key].ID
		}
		server.Category.ID = dbServer.Category.ID
		for key, channel := range server.Channels {
			channel.ID = dbServer.Channels[key].ID
		}
		for key, message := range server.Messages {
			message.ID = dbServer.Messages[key].ID
			message.ChannelID = dbServer.Messages[key].ChannelID
		}
		server.Users = dbServer.Users
	}

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
func UpdateServerPlaying(s *structure.Server) error {
	client := connect()
	defer client.Disconnect(context.Background())
	collection := client.Database("knights-of-discord").Collection("servers")
	filter := bson.NewDocument(bson.EC.String("id", s.ID))
	replacement := bson.NewDocument(bson.EC.SubDocumentFromElements("$set", bson.EC.Boolean("playing", s.Playing)))
	_, err := collection.UpdateOne(context.Background(), filter, replacement)
	return err
}

// AddServerUser adds a new User to the game
func AddServerUser(s *structure.Server, u *structure.User) error {
	client := connect()
	defer client.Disconnect(context.Background())
	collection := client.Database("knights-of-discord").Collection("servers")
	filter := bson.NewDocument(bson.EC.String("id", s.ID))
	replacement := bson.NewDocument(bson.EC.SubDocumentFromElements("$push", bson.EC.Interface("users", u)))
	_, err := collection.UpdateOne(context.Background(), filter, replacement)
	return err
}

// RemoveServerUser removes an existing User from the game
func RemoveServerUser(s *structure.Server, u *structure.User) error {
	client := connect()
	defer client.Disconnect(context.Background())
	collection := client.Database("knights-of-discord").Collection("servers")
	filter := bson.NewDocument(bson.EC.String("id", s.ID))
	replacement := bson.NewDocument(bson.EC.SubDocumentFromElements("$pull", bson.EC.Interface("users", u)))
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
