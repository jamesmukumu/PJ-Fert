package users

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jamesmukumu/backup/db"
	"github.com/jamesmukumu/backup/schema/admin"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"

	"gopkg.in/gomail.v2"

	"golang.org/x/crypto/bcrypt"
)
var myuser admin.User



//send mail

func  Sendmail(uniquestring string){

//load dotenv
dotenv := godotenv.Load()
if dotenv != nil {
	log.Fatal(dotenv)
}


Mypassword := os.Getenv("GmailPassword")
MyEmail := os.Getenv("Email")


	


mail := gomail.NewMessage()


mail.SetHeader("From",MyEmail)
mail.SetHeader("To",myuser.Email)
mail.SetHeader("Subject","Thank you for registering")
mail.SetBody("text/plain",uniquestring)


//set up dialer

dialer := gomail.NewDialer("smtp.gmail.com",587,MyEmail,Mypassword)
fmt.Println(dialer)


dialer.TLSConfig = &tls.Config{InsecureSkipVerify:false,
	ServerName: "smtp.gmail.com",
}



err := dialer.DialAndSend(mail)
if err != nil {
	log.Fatal(err)
}

}















//create token

func Createtoken(username string) string {
    dotenv := godotenv.Load()
    fmt.Println(dotenv)
    jwtsecret := os.Getenv("jwtSecret")
    secret := []byte(jwtsecret)

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "username": username,
        "exp":      time.Now().Add(time.Minute * 15).Unix(),
    })

    signedToken, err := token.SignedString(secret)

    if err != nil {
        panic(err)
    } else {
        return signedToken
    }
}


//generate string
func Generatestring(length int)(string, error){
buffer := make([]byte, length)

_,err := rand.Read(buffer)
if err != nil {
	return "",err
}else{
	return base64.URLEncoding.EncodeToString(buffer)[:length],nil
}



}














func Postuser(res http.ResponseWriter, req *http.Request) {
	db.Connectmongo()
   
	// grab the class User

 
	// decode the request.body
	err := json.NewDecoder(req.Body).Decode(&myuser)
 
	if err != nil {
		http.Error(res, "Error in getting values from the request", http.StatusBadRequest)
		return
	}


   if !myuser.Validatesex() {
    http.Error(res,"Sex must be female or male",http.StatusBadRequest)
    return
   }


   if !myuser.Checkemail() {
    http.Error(res,"Email is of invalid format",http.StatusBadRequest)
    return
   }
 

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(myuser.Password), 14)

	myuser.Password = string(hashedPassword)
	myuser.Uniquestring, _ = Generatestring(10)
	data, err := db.Mongodb.InsertOne(context.Background(), myuser)
	if err != nil {
		http.Error(res, fmt.Sprintf("Error inserting data: %v", err), http.StatusInternalServerError)
		return
	}
	fmt.Println(data)
  token := Createtoken(myuser.Username)
	// Respond with the desired message
    res.Header().Set("Authorization",token)
	json.NewEncoder(res).Encode(map[string]string{"message":"posted sucessfully"})


    

	

    defer Sendmail(myuser.Uniquestring)

}







//login 

func Loginuser(res http.ResponseWriter, req *http.Request){
    var myuser admin.User
    //make a connection to mongodb
    db.Connectmongo()

    //decode request
    err := json.NewDecoder(req.Body).Decode(&myuser)
    if err != nil {
        http.Error(res,"Error in decoding json",http.StatusBadRequest)
        return
    }

    var dbuser admin.User

    //find matching username
    err = db.Mongodb.FindOne(context.Background(),bson.M{"username":myuser.Username}).Decode(&dbuser)
    if err != nil {
        http.Error(res,"Username not found",http.StatusUnauthorized)
        return
    }

    // Compare the user's password with the stored hashed password
    err = bcrypt.CompareHashAndPassword([]byte(dbuser.Password), []byte(myuser.Password))
    if err != nil {
        json.NewEncoder(res).Encode(map[string]string{"message":"incorrect password"})
        return
    } else {
		token := Createtoken(myuser.Username)
		res.Header().Set("Authorization",token)
      json.NewEncoder(res).Encode(map[string]string{"message":"login sucess"})
    }
}




// login with sent mongodb
func Loginwithid(res http.ResponseWriter, req *http.Request){
	db.Connectmongo()
    var dbmyuser admin.User

    //get unique string
    inputUniquestring := json.NewDecoder(req.Body).Decode(&dbmyuser)
    if inputUniquestring != nil {
        log.Fatal(inputUniquestring)
        return
    }


  var testUser admin.User

    //create filter
    filter := bson.M{"uniquestring":dbmyuser.Uniquestring}
 
    //try to find  
    matchingUniquestring := db.Mongodb.FindOne(context.Background(),filter).Decode(&testUser)
    if matchingUniquestring != nil {
        res.Write([]byte("No matching unique string"))
        return
    }
    token := Createtoken(dbmyuser.Username)
    res.Header().Set("Authorization",token) 
    json.NewEncoder(res).Encode(map[string]string{"message":"login successfully"})
} 
 
 



//update password
func Updatepassword(res http.ResponseWriter , req *http.Request){
    var emaildb admin.User
    //query the email
    queryEmail := req.URL.Query().Get("Email")

    Newpassword := json.NewDecoder(req.Body).Decode(&emaildb)
    if Newpassword != nil {
        log.Fatal(Newpassword)     
        return
    }

    myhashedPassword, _ := bcrypt.GenerateFromPassword([]byte(emaildb.Password),14)
       emaildb.Password = string(myhashedPassword)
    update := bson.M{"$set": bson.M{"password":emaildb.Password}}

    filter := bson.M{
        "Email":queryEmail,  

    }
    updatedPassandfoundemail := db.Mongodb.FindOneAndUpdate(context.Background(),filter,update).Decode(&myuser)

    if updatedPassandfoundemail != nil {
        http.Error(res,"error",http.StatusUnauthorized)   
        return 
    }
  fmt.Println(updatedPassandfoundemail) 
    json.NewEncoder(res).Encode(map[string]string{"message":"Password changed"})

}
      
//delete account

func Deletdaccount(res http.ResponseWriter,req *http.Request){
    //first check if email exists if it exists we delete
    // var deletedACC admin.User
    deleteEmail := req.URL.Query().Get("Email")
     if deleteEmail=="" {
        http.Error(res,"Provide a query email",http.StatusBadRequest)
     }


    filter :=bson.M{
        "Email":deleteEmail,
    }


    actualDeletion, err := db.Mongodb.DeleteOne(context.Background(),filter)
  if actualDeletion.DeletedCount == 0 {
    log.Fatal(actualDeletion)
    http.Error(res,"No document deleted",http.StatusInternalServerError)
    return
  }
if err != nil {
    log.Fatal(err) 
}
 

    json.NewEncoder(res).Encode(map[string]string{"message":"Account deleted sucessfully"})
  


 




}
 

 




 