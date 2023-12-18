package db

import (
	"context"
	"log"
	"os"
     "fmt"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)
var Mongodbmenses *mongo.Collection 
func Mensesdb() {

	dotenv := godotenv.Load()
     if dotenv != nil {
		log.Fatal(dotenv)
	 }







connectionStringmenses := os.Getenv("connectionString")



//make actual connection
optionsforConnectionmenses := options.Client().ApplyURI(connectionStringmenses)

Conne,err := mongo.Connect(context.TODO(),optionsforConnectionmenses)
if err != nil{
	log.Fatal(err)
}else{
	fmt.Println("Connected to menses db sucessfully")
}



Mongodbmenses = Conne.Database("Menses").Collection("myusers menses")






}
