package admin

import "strings"

// import (
// 	"go.mongodb.org/mongo-driver/bson/primitive"
// )

type User struct {
	
	Username string             `json:"username" bson:"username"`
	Email    string             `json:"Email" bson:"Email"`
	Password string             `json:"password" bson:"password"`
	Sex    string               `json:"sex" bson:"sex"`
	Uniquestring string          `json:"uniquestring" bson:"uniquestring"`
}   


func (u *User)Validatesex()bool{
	return u.Sex == "female" || u.Sex == "male"
	} 



	// verify email format
func (email *User)Checkemail()bool{
	return strings.Contains(email.Email,"@")
}