package request

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"ts_mcp/constants"
	"ts_mcp/utils"
)

type TimeSheetRequest struct {
	BaseURL string
}

func NewTimeSheetRequest(BaseURL string) *TimeSheetRequest {
	return &TimeSheetRequest{BaseURL: BaseURL}
}

func (tsr *TimeSheetRequest) GetUserToken(email string, password string) (token string, err error) {
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

func (tsr *TimeSheetRequest) GetPersonID(token string) (id int, err error) {
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

func (tsr *TimeSheetRequest) PostTimeSheetEntry(token string, taskID int, personID int, costCodeID int, overtime bool, time int, date time.Time, description string) (err error) {
	type TimeSheetEntry struct {
		TaskID       int    `json:"TaskId"`
		PersonID     int    `json:"PersonId"`
		CostCodeID   int    `json:"CostCodeId"`
		DepartmentID int    `json:"DepartmentId"`
		Overtime     int    `json:"Overtime"`
		Time         int    `json:"Time"`
		EntryDate    string `json:"EntryDate"`
		Description  string `json:"Comments"`
		WorklogID    int    `json:"WorklogId"`
		Audited      int    `json:"Audited"`
	}

	body := TimeSheetEntry{
		TaskID:       taskID,
		PersonID:     personID,
		CostCodeID:   costCodeID,
		DepartmentID: 1,
		Overtime:     utils.Bool2int(overtime),
		Time:         time,
		EntryDate:    date.Format(constants.TimeFormat),
		Description:  description,
		WorklogID:    0,
		Audited:      0,
	}

	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(body)

	req, err := http.NewRequest("POST", tsr.BaseURL+"/api/entry/create", payloadBuf)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %v\n", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API returned non-OK status: %s", resp.Status)
	}

	return nil
}
