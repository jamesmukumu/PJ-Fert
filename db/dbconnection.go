package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Mongodb *mongo.Collection
var Client *mongo.Client

func Connectmongo() {
	var db = "Users"
	var mycollection = "myusers"

	dotenv := godotenv.Load()

	if dotenv != nil {
		log.Fatal(dotenv)
	}

	connectionString := os.Getenv("connectionString")

	// Create a MongoDB client
	connectionOptions := options.Client().ApplyURI(connectionString)
	var err error  
	Client, err = mongo.Connect(context.TODO(), connectionOptions)  
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Connection to MongoDB has been made")
	}

	// Set up the MongoDB collection
	Mongodb = Client.Database(db).Collection(mycollection)

	// Create unique index for username
	usernameIndexOptions := options.Index().SetUnique(true)
	usernameIndexModel := mongo.IndexModel{
		Keys:    bson.M{"username": 1},
		Options: usernameIndexOptions,
	}
	_, err = Mongodb.Indexes().CreateOne(context.Background(), usernameIndexModel)
	if err != nil {
		log.Fatal(err)
	}

	// Create unique index for email
	emailIndexOptions := options.Index().SetUnique(true)
	emailIndexModel := mongo.IndexModel{
		Keys:    bson.M{"email": 1},
		Options: emailIndexOptions,
	}
	_, err = Mongodb.Indexes().CreateOne(context.Background(), emailIndexModel)
	if err != nil {
		log.Fatal(err)
	}
}
