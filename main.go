package main

import (
	"fmt"
	"net/http"
	"time"

	"web/usfl-backend/handlers"
	"web/usfl-backend/mongohandler"

	"github.com/labstack/echo"
	"go.mongodb.org/mongo-driver/mongo"
)

var MongoClient *mongo.Client

func main() {
	err := mongohandler.InitialDatabase(MongoClient)
	defer mongohandler.Close(MongoClient)

	if err != nil {
		fmt.Println(err)
	}

	e := echo.New()
	s := &http.Server{
		Addr:         ":2000",
		ReadTimeout:  20 * time.Minute,
		WriteTimeout: 20 * time.Minute,
	}

	e.GET("/api/allbook", handlers.GetAllBook)


	e.Logger.Fatal(e.StartServer(s))
}