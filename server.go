package main

import (
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func main() {

	r := mux.NewRouter()

	fmt.Printf("Hello")
	//http.Handle("/foo", fooHandler)

	r.HandleFunc("/authenticate", authenticate)

	http.HandleFunc("/bar", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func authenticate(response http.ResponseWriter, request *http.Request) {

	secret := os.Getenv("AUTH_SECRET")
	if secret == "" {
		log.Fatal("Secret key not found")
	}

	signingKey := []byte(secret)
	duration, _ := time.ParseDuration("24h")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp":  time.Now().Add(duration).Unix(),
		"test": "yaj",
	})
	ss, err := token.SignedString(signingKey)

	if err != nil {
		log.Error("Failed to authenticate due to: ")
		response.WriteHeader(302)
		log.Error(err)
		return
	}

	type AuthResponse struct {
		Token string
	}

	authResponse := AuthResponse{ss}
	data, jsonErr := json.Marshal(authResponse)

	if jsonErr != nil {
		log.Error("Failed to authenticate2 due to: ")
		log.Error(jsonErr)
		response.WriteHeader(302)
		return
	}

	response.WriteHeader(200)
	response.Write(data)
	os.Stdout.Write(data)
	return

}
