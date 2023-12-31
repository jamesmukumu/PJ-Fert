package routers

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	middlewares "github.com/jamesmukumu/backup/Middlewares"
	cycles "github.com/jamesmukumu/backup/controllers/Cycles"
	users "github.com/jamesmukumu/backup/controllers/Users"
	"github.com/jamesmukumu/backup/db"
)

func Handleallroutes() {
	fmt.Println("Listening for requests at 7000")
	db.Connectmongo()
	dotenv := godotenv.Load()
	
      if dotenv != nil {
		log.Fatal(dotenv.Error())
	  }

     port := os.Getenv("PORT")
	Router := mux.NewRouter()
	Router.HandleFunc("/post/user",users.Postuser).Methods("POST")
    Router.HandleFunc("/login/user",users.Loginuser).Methods("POST")
    Router.HandleFunc("/login/uniquestring",users.Loginwithid).Methods("POST")
    Router.HandleFunc("/change/password",users.Updatepassword).Methods("PUT")
    Router.HandleFunc("/delete/account",users.Deletdaccount).Methods("DELETE")
 
 




    //menses
  Router.HandleFunc("/fetch/allmenses",cycles.Fetchallmenses).Methods("GET")
   Router.HandleFunc("/post/menses",middlewares.TokenMiddleware(cycles.Postmenses)).Methods("POST")
   Router.HandleFunc("/get/menses/email",middlewares.TokenMiddleware(cycles.Getmensesprediction)).Methods("GET")
   Router.HandleFunc("/delete/menses",middlewares.TokenMiddleware(cycles.Deletemenses)).Methods("DELETE")
   Router.HandleFunc("/fetch/pastsafedays",cycles.FetchAccountsduesafe).Methods("GET")


 server :=http.ListenAndServe(":"+port,Router)
if  server !=nil {
	log.Fatal(server)
}
	




} 






