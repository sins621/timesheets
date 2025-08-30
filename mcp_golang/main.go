package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"gorm.io/gorm"
)

const BASE_URL = "https://office.warpdevelopment.com"

type User struct {
	gorm.Model
	Email         string `gorm:"uniqueIndex"`
	Token         string
	InitializedAt time.Time `gorm:"not null"`
}

func main() {
	exePath, err := os.Executable()

	if err != nil {
		panic(fmt.Sprintf("Failed to get executable path: %v\n", err))
	}
	exeDir := filepath.Dir(exePath)
	dbPath := filepath.Join(exeDir, "timesheets.db")

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})

	if err != nil {
		panic("Database failed to initialize")
	}

	err = db.AutoMigrate(&User{})

	if err != nil {
		panic("Error running migration")
	}

	s := server.NewMCPServer(
		"Demo 🚀",
		"1.0.0",
		server.WithToolCapabilities(false),
	)

	authTool := mcp.NewTool("Get Timesheet Token",
		mcp.WithDescription("Get the Token from Timesheets Endpoint Using Username and Password"),
	)

	s.AddTool(authTool, authHandler)

	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}

func authHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	email, exists := os.LookupEnv("EMAIL")
	if !exists {
		return mcp.NewToolResultError("Email does not exist in environment."), nil
	}

	password, exists := os.LookupEnv("PASSWORD")
	if !exists {
		return mcp.NewToolResultError("Password does not exist in environment"), nil
	}

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
		return mcp.NewToolResultError(fmt.Sprintf("Error marshaling JSON: %v\n", err)), nil
	}

	resp, err := http.Post(
		BASE_URL+"/api/account/authorise",
		"application/json",
		bytes.NewBuffer(jsonData),
	)

	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Error making HTTP Request to %s: %v\n", BASE_URL, err)), nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return mcp.NewToolResultError(fmt.Sprintf("Request Failed with status: %d\n", resp.StatusCode)), nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Error reading response: %v\n", err)), nil
	}

	var responseData ResponseBody
	err = json.Unmarshal(body, &responseData)

	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Error parsing json response: %v\n", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("The authorizatin token is: %s\n", responseData.Token)), nil
}
