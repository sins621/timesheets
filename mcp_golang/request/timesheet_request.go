package request

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type TimeSheetRequest struct {
	BaseURL string
}

func NewTimeSheetRequest(BaseURL string) *TimeSheetRequest {
	return &TimeSheetRequest{BaseURL: BaseURL}
}

func (tsr *TimeSheetRequest) RequestUserToken(email string, password string) (token string, err error) {
	type RequestBody struct {
		Email    string `json:"Email"`
		Password string `json:"Password"`
	}

	type ResponseBody struct {
		Token string `json:"token"`
	}

	requestData := RequestBody{
		Email:    email,
		Password: password,
	}

	jsonData, err := json.Marshal(requestData)

	if err != nil {
		return "", fmt.Errorf("error marshaling JSON: %v", err)
	}

	resp, err := http.Post(
		tsr.BaseURL+"/api/account/authorise",
		"application/json",
		bytes.NewBuffer(jsonData),
	)

	if err != nil {
		return "", fmt.Errorf("error making HTTP request to %s: %v", tsr.BaseURL+"/api/account/authorise", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("request failed with status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %v", err)
	}

	var responseData ResponseBody
	err = json.Unmarshal(body, &responseData)

	if err != nil {
		return "", fmt.Errorf("error parsing json response: %v", err)
	}

	return responseData.Token, nil
}

func (tsr *TimeSheetRequest) RequestPersonID(token string) (id int, err error) {
	type Person struct {
		PersonID          int    `json:"PersonId"`
		FirstName         string `json:"FirstName"`
		Surname           string `json:"Surname"`
		Email             string `json:"Email"`
		TelephoneNumber   string `json:"TelephoneNumber"`
		IsAdmin           bool   `json:"is_admin"`
		PersonStatus      string `json:"PersonStatus"`
		CreatedOnUtc      string `json:"CreatedOnUtc"`
		ModifiedOnUtc     string `json:"ModifiedOnUtc"`
		ProfilePictureURL string `json:"ProfilePictureUrl"`
	}

	req, err := http.NewRequest("GET", tsr.BaseURL+"/api/users/me", nil)
	if err != nil {
		return 0, fmt.Errorf("error creating request: %v\n", err)
	}

	req.Header.Add("Authorization", "Bearer "+token)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("error sending request: %v\n", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("API returned non-OK status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("error reading response body: %v", err)
	}

	var p Person
	err = json.Unmarshal(body, &p)
	if err != nil {
		return 0, fmt.Errorf("error unmarshaling JSON: %v", err)
	}

	return p.PersonID, nil
}
