package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("No environment file")
	}

	const warpBaseUrl = "https://office.warpdevelopment.com"

	type User struct {
		Email    string
		Password string
	}

	type Reponse struct {
		Token string `json:"token"`
	}

	user := User{os.Getenv("email"), os.Getenv("password")}
	signInData := map[string]string{"Email": user.Email, "Password": user.Password}
	signInBody, err := json.Marshal(signInData)

	if err != nil {
		panic(fmt.Sprintf("Error marshalling sign in body: %s", err))
	}

	req, err := http.NewRequest(http.MethodPost, warpBaseUrl+"/api/account/Authorise", bytes.NewReader(signInBody))

	if err != nil {
		panic(fmt.Sprintf("Error making request object: %s", err))
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	token := &Reponse{}

	derr := json.NewDecoder(res.Body).Decode(token)

	if derr != nil {
		panic(derr)
	}

	if res.StatusCode != http.StatusOK {
		panic(res.Status)
	}

	fmt.Printf("Token %s", token.Token)
}
