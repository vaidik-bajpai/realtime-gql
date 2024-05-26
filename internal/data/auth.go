package data

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/pascaldekloe/jwt"
)

func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")

		authToken := r.Header.Get("Authorization")
		if authToken == "" {
			r = ContextSetUser(r, AnonymousUser)
			next.ServeHTTP(w, r)
			return
		}

		headerParts := strings.Split(authToken, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			log.Fatal("invalid token 1")
			return
		}

		token := headerParts[1]

		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			http.Error(w, "server configuration error", http.StatusInternalServerError)
			return
		}

		claims, err := jwt.HMACCheck([]byte(token), []byte(secret))
		if err != nil {
			log.Fatal("invalid token 2")
			return
		}

		if !claims.Valid(time.Now()) {
			log.Fatal("invalid token 3")
			return
		}

		if claims.Issuer != "bajpai" {
			log.Fatal("invalid token 3")
			return
		}

		if !claims.AcceptAudience("bajpai") {
			log.Fatal("invalid token 4")
			return
		}

		fmt.Println("claims.Subject:", claims.Subject)
		userID, err := strconv.ParseInt(claims.Subject, 10, 64)
		fmt.Println("Parsed userID:", userID, "Error:", err)
		if err != nil || userID == 0 {
			log.Fatal("invalid token 4")
			return
		}

		user, err := Model.Users.Get(int(userID)) // Replace with your user fetching logic
		if err != nil {
			return
		}

		r = ContextSetUser(r, user)

		next.ServeHTTP(w, r)
	})
}
