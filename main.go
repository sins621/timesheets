package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func main() {
	const warpBaseUrl = "https://office.warpdevelopment.com"

	type User struct {
		Email    string
		Password string
	}

	user := User{os.Getenv("email"), os.Getenv("password")}
	signInData := map[string]string{"Email": user.Email, "Password": user.Password}
	signInBody, err := json.Marshal(signInData)

	if err != nil {
		fmt.Printf("Error marshalling sign in body: %s", err)
	}

	req, err := http.NewRequest(http.MethodPost, warpBaseUrl+"api/account/Authorise", bytes.NewReader(signInBody))

	if err != nil {
		fmt.Printf("Error making request object")
	}

	req.Header.Set("Content-Type", "application/json")

	if err != nil {
		fmt.Printf("Error making http request: %s", err)
		os.Exit(1)
	}
}
