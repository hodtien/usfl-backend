package mongohandler

import (
	"fmt"
	// "time"
	"context"
	// "encoding/json"
	
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	// "github.com/cavdy-play/go_mongo/controllers"
)

// MongoHost - Mongo URL
var	MongoHost = "mongodb+srv://hodtien:1712181HCMUS@usfl-ksxxv.gcp.mongodb.net/?retryWrites=true&w=majority"
// DBName - Database name
var DBName = "usfl-database"

// InitialDatabase - InitialDatabase
func InitialDatabase(MongoClient *mongo.Client) (err error) {
	clientOptions := options.Client().ApplyURI(MongoHost)

	// Connect to MongoDB
	MongoClient, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		fmt.Println(err)
	}

	// Check the connection
	err = MongoClient.Ping(context.TODO(), nil)

	if err != nil {
		fmt.Println(err)
	}

	return (err)
}

// Close - Close connection to mongoDB
func Close(MongoClient *mongo.Client) {
	err := MongoClient.Disconnect(context.TODO())

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Connection to MongoDB closed.")
}