package mongohandler

import (
	"gopkg.in/mgo.v2"
)

// Mgo - Mongo DB Struct
type Mgo struct{}

//DB - DB
var db *mgo.Session

// Mgodb - MongoDB
var Mgodb Mgo

// MongoHost - Mongo URL
var	MongoHost = "localhost:27017"

// DBName - Database name
var DBName = "usfl-database"

//Limit - Limit
const Limit = 5