package db

import (
	"context"

	"github.com/Noxdew/Knights-Of-Discord/config"
	"github.com/Noxdew/Knights-Of-Discord/logger"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
)

var NotFound = mongo.ErrNoDocuments

// Server represents a Discord Guild in the game
type Server struct {
	ID    string `json:"id" bson:"id"`
	Power int    `json:"power" bson:"power"`
}

// Role represents a Discord Role in the game
type Role struct {
	ID       string `json:"id" bson:"id"`
	ServerID string `json:"serverId" bson:"serverId"`
	Type     string `json:"type" bson:"type"`
}

// Channel represents a Discord text Channel in the game
type Channel struct {
	ID       string `json:"id" bson:"id"`
	ServerID string `json:"serverId" bson:"serverId"`
	Type     string `json:"type" bson:"type"`
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

// GetRoles returns all game Roles associated with Server s
func GetRoles(s string) (*[]Role, error) {
	client := connect()
	collection := client.Database("knights-of-discord").Collection("roles")
	filter := bson.NewDocument(bson.EC.String("serverId", s))
	docs, err := collection.Find(context.Background(), filter)
	con := context.Background()
	defer docs.Close(con)
	var roles []Role
	for docs.Next(con) {
		role := Role{}
		err := docs.Decode(&role)
		if err != nil {
			logger.Log.Error(err.Error())
		}
		roles = append(roles, role)
	}
	return &roles, err
}

// GetRole returns a game Role with ID s
func GetRole(g string, s string) (*Role, error) {
	client := connect()
	defer client.Disconnect(context.Background())
	collection := client.Database("knights-of-discord").Collection("roles")
	filter := bson.NewDocument(bson.EC.String("serverId", g), bson.EC.String("type", s))
	role := Role{}
	doc := collection.FindOne(context.Background(), filter)
	err := doc.Decode(&role)
	return &role, err
}

// CreateRole uploads a new role to the DB
func CreateRole(r Role) error {
	client := connect()
	defer client.Disconnect(context.Background())
	collection := client.Database("knights-of-discord").Collection("roles")
	_, err := collection.InsertOne(context.Background(), r)
	return err
}

// UpdateRole updates an existing role in the DB
func UpdateRole(r Role) error {
	client := connect()
	defer client.Disconnect(context.Background())
	collection := client.Database("knights-of-discord").Collection("roles")
	filter := bson.NewDocument(bson.EC.String("serverId", r.ServerID), bson.EC.String("type", r.Type))
	replacement := bson.NewDocument(bson.EC.SubDocumentFromElements("$set", bson.EC.String("id", r.ID)))
	_, err := collection.UpdateOne(context.Background(), filter, replacement)
	return err
}

// RemoveRoles removes server s roles from the DB
func RemoveRoles(s string) error {
	client := connect()
	defer client.Disconnect(context.Background())
	collection := client.Database("knights-of-discord").Collection("roles")
	filter := bson.NewDocument(bson.EC.String("serverId", s))
	_, err := collection.DeleteMany(context.Background(), filter)
	return err
}
