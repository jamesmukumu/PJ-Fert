package middlewares

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)









func verifyToken(tokenString string) error {
   dotenv := godotenv.Load()
   fmt.Println(dotenv)
   jwtsecret := os.Getenv("jwtSecret")
	var secretKey = []byte(jwtsecret)



	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}






func TokenMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		authorizationHeader := req.Header.Get("Authorization")
		if authorizationHeader == "" {
			http.Error(res, "Unauthorized", http.StatusUnauthorized)
			return
		}

		tokenString := authorizationHeader

		err := verifyToken(tokenString)
		if err != nil {
			json.NewEncoder(res).Encode(map[string]string{"error":err.Error()})
			return
		} 

		next.ServeHTTP(res, req)
	}
}
