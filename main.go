package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type User struct {
	Email    string
	Password string
}

func main() {

}

func get(url string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return "could not make request", err
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return "error making http request", err
	}
	defer resp.Body.Close()

	resBody, err := io.ReadAll(resp.Body)

	if err != nil {
		return fmt.Sprintf("could not read response body: %s\n", err), err
	}

	return fmt.Sprintf("response body: %s\n", resBody), nil
}

func getWithToken(url string, token string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	if err != nil {
		return "could not make request", err
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return "error making http request", err
	}
	defer resp.Body.Close()

	resBody, err := io.ReadAll(resp.Body)

	if err != nil {
		return fmt.Sprintf("could not read response body: %s\n", err), err
	}

	return fmt.Sprintf("response body: %s\n", resBody), nil
}

func post(url string, data map[string]string) (string, error) {
	jsonData, err := json.Marshal(data)

	if err != nil {
		return "Error Converting Data into Json", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))

	resp, err := http.DefaultClient.Do(req)
	resp.Body.Close()

	statusCode := resp.StatusCode
	resBody, err := io.ReadAll(resp.Body)

	if err != nil {
		return fmt.Sprintf("could not read response body: %s\n", err), err
	}

	return fmt.Sprintf("response body: %s\n status code: %d\n", resBody, statusCode), nil
}

func postWithToken(url string, token string, data map[string]string) (string, error) {
	jsonData, err := json.Marshal(data)

	if err != nil {
		return "Error Converting Data into Json", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	resp.Body.Close()

	statusCode := resp.StatusCode
	resBody, err := io.ReadAll(resp.Body)

	if err != nil {
		return fmt.Sprintf("could not read response body: %s\n", err), err
	}

	return fmt.Sprintf("response body: %s\n status code: %d\n", resBody, statusCode), nil
}
