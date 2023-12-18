package cycles

import (
	"context"
	"crypto/tls"
	"os"

	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"time"

	"github.com/jamesmukumu/backup/db"
	"github.com/jamesmukumu/backup/schema/menses"
	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"

	"go.mongodb.org/mongo-driver/bson"
)

var clientmenses menses.Menses

func Postmenses(res http.ResponseWriter, req *http.Request) {
	db.Mensesdb()

	mensesPosted := json.NewDecoder(req.Body).Decode(&clientmenses)
	if mensesPosted != nil {
		log.Fatal(mensesPosted)
	}

	lastcyledatetime, err := time.Parse("2006-01-02", clientmenses.Lastcycledate)
	clientmenses.Lastcycledatetime = lastcyledatetime
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < 5; i++ {
		Nextcycledate := clientmenses.Lastcycledatetime.AddDate(0, 0, clientmenses.Normalcycleday)
		clientmenses.Nextexpectedperioddate = Nextcycledate

		safedays := clientmenses.Lastcycledatetime.AddDate(0,0,7)
		clientmenses.Safedays = safedays

		inserteddata, err := db.Mongodbmenses.InsertOne(context.Background(), clientmenses)
		fmt.Println(inserteddata)

		if err != nil {
			json.NewEncoder(res).Encode(err)
			return
		}

		clientmenses.Lastcycledatetime = clientmenses.Nextexpectedperioddate

		//find the inseretdobject because insert one only returns
		filter := bson.M{
			"_id": inserteddata.InsertedID, 
		}

		data := db.Mongodbmenses.FindOne(context.Background(), filter)
		fmt.Println(data)

		//create a new instance
		var testmenses menses.Menses

		decodedData := data.Decode(&testmenses)
		fmt.Println(decodedData)
		if decodedData != nil {
			log.Fatal(decodedData)
		}

	}

	json.NewEncoder(res).Encode(map[string]string{"message": "Menses added"})

}

var fetchedMenses menses.Menses

// get menses
func Getmensesprediction(res http.ResponseWriter, req *http.Request) {
	db.Mensesdb()

	// Get email query parameter from the URL
	email := req.URL.Query().Get("Email.Email")
	if email == "" {
		http.Error(res, "Email parameter is missing", http.StatusBadRequest)
		return
	}

	filter := bson.M{
		"Email.Email": email,
	}

	queriedData := db.Mongodbmenses.FindOne(context.Background(), filter)
	if queriedData.Err() != nil {
		log.Println(queriedData.Err())
		http.Error(res, "Error querying data", http.StatusInternalServerError)
		return
	}

	err := queriedData.Decode(&fetchedMenses)
	if err != nil {
		log.Println(err)
		http.Error(res, "Error decoding data", http.StatusInternalServerError)
		return
	}

	jsonFormated, _ := json.Marshal(fetchedMenses)
	defer Sendemailsinglemenses(string(jsonFormated))
	res.Header().Set("Content-Type", "application/json")
	res.Write(jsonFormated)
}

//get all menses based

var mensesData []menses.Menses

func Fetchallmenses(res http.ResponseWriter, req *http.Request) {

	//query all the menses through

	queryFormenses := req.URL.Query().Get("Email.Email")
	if queryFormenses == "" {
		res.Write([]byte("No query has been passed"))
	}

	//actual filter
	filter := bson.M{
		"Email.Email": queryFormenses,
	}

	queriedData, err := db.Mongodbmenses.Find(context.Background(), filter)
	if err != nil {
		log.Println(err)
		http.Error(res, "Error querying data", http.StatusInternalServerError)
		return
	}

	if err = queriedData.All(context.Background(), &mensesData); err != nil {
		log.Println(err)
		http.Error(res, "Error decoding data", http.StatusInternalServerError)
		return
	}

	//  jsonFormated, _ := json.Marshal(mensesData)
	// fmt.Println(jsonFormated)
	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(map[string]string{"message": "Data fetched"})
	json.NewEncoder(res).Encode(map[string]interface{}{
		"message": "fetched data from mongodb",
		"data":    mensesData,
	})
	// defer Sendemailmenses(string(jsonFormated))

}

func Sendemailsinglemenses(contentEmail string) {
	dotenv := godotenv.Load()
	if dotenv != nil {
		log.Fatal(dotenv)
	}

	Mypassword := os.Getenv("GmailPassword")
	MyEmail := os.Getenv("Email")

	mail := gomail.NewMessage()

	mail.SetHeader("From", MyEmail)
	mail.SetHeader("To", fetchedMenses.Email.Email)
	mail.SetHeader("Subject", "Thank you for registering")
	mail.SetBody("text/plain", contentEmail)

	//set up dialer

	dialer := gomail.NewDialer("smtp.gmail.com", 587, MyEmail, Mypassword)
	fmt.Println(dialer)

	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: false,
		ServerName: "smtp.gmail.com",
	}

	err := dialer.DialAndSend(mail)
	if err != nil {
		log.Fatal(err)
	}

}

//fetch Menses where safedaya has reached and send an alert 
func FetchAccountsduesafe(res http.ResponseWriter,req *http.Request){
// var Duemenses menses.Menses


filteronDuetime := bson.M{
	"safedays":bson.M{"$lte":time.Now()},
}

resultsoffilter,err := db.Mongodbmenses.Find(context.Background(),filteronDuetime)
if err != nil {
log.Fatal(err)	
}






jsonFormatedresults,_ := json.Marshal(resultsoffilter)
res.Write(jsonFormatedresults)


}










//delete a menses
func Deletemenses(res http.ResponseWriter,req *http.Request){


deletionQuery := req.URL.Query().Get("Email.Email")
if deletionQuery == "" {
	json.NewEncoder(res).Encode(map[string]string{"message":"Please provide a deletionquery"})
}


//filter where it appears

filter := bson.M{
	"Email.Email":deletionQuery,
}




deletedmenses,err :=  db.Mongodbmenses.DeleteOne(context.Background(),filter)
if err != nil {
	log.Fatal(err)
}


if deletedmenses.DeletedCount == 0 {
json.NewEncoder(res).Encode("No document deleted")
}



jsonDeleted, _ := json.Marshal(deletedmenses)

res.Write(jsonDeleted)



}