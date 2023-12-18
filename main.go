package main

import (
	"fmt"

	"github.com/jamesmukumu/backup/db"
	"github.com/jamesmukumu/backup/routers"
)  

func main(){ 
	fmt.Println("hello world")
	  db.Connectmongo()
	  db.Mensesdb()
	 routers.Handleallroutes() 
	
} 


 