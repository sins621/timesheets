package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

type MockDatabase struct {
}

func (m *MockDatabase) CreateUser(*User) error {
	return nil
}

func (m *MockDatabase) GetUserByEmail(email string) (*User, error) {
	var user User

	return &user, nil
}

func (m *MockDatabase) UpdateUserToken(email, token string) error {
	return nil
}

func (m *MockDatabase) GetUserByID(id uint) (*User, error) {
	var user User

	return &user, nil
}

func TestUpdateToken(t *testing.T) {
	mockDb := &MockDatabase{}
	handler := NewHandler(mockDb)

	err := godotenv.Load()

	if err != nil {
		t.Fatal("failed to load environment.")
	}

	email, exists := os.LookupEnv("EMAIL")
	if !exists {
		t.Fatal("email does not exist in environment")
	}

	password, exists := os.LookupEnv("PASSWORD")
	if !exists {
		t.Fatal("password does not exist in environment")
	}

	token, err := handler.updateUserToken(email, password)

	if err != nil {
		t.Fatalf("updating user token failed with err: %v\n", err)
	}

	fmt.Printf("user token: %s\n", token)
}
