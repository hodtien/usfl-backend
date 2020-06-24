package handlers

import (
	"context"
	"fmt"
	"web/usfl-backend/models"

	"github.com/labstack/echo"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
	// "web/usfl-backend/mongohandler"
)

var db = "usfl-database"
var MongoClient *mongo.Client

func GetAllBook(c echo.Context) (err error) {
	collectionNames := "all@book" 
	collection := MongoClient.Database(db).Collection(collectionNames)


	findOptions := options.Find()
	findOptions.SetLimit(2)

	// Here's an array in which you can store the decoded documents
	var results []*models.Book

	// Passing bson.D{{}} as the filter matches all documents in the collection
	cur, err := collection.Find(context.TODO(), bson.D{{}}, findOptions)
	if err != nil {
		fmt.Println(err)
	}

	// Finding multiple documents returns a cursor
// Iterating through the cursor allows us to decode documents one at a time
for cur.Next(context.TODO()) {
    
    // create a value into which the single document can be decoded
    var book models.Book
    err := cur.Decode(&book)
    if err != nil {
		fmt.Println(err)
    }

    results = append(results, &book)
}

	if err := cur.Err(); err != nil {
		fmt.Println(err)
	}

	// Close the cursor once finished
	cur.Close(context.TODO())

	// resp := "abc"
	fmt.Println("All book:", results)
	return c.JSON(200, results)

}