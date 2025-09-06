package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/mark3labs/mcp-go/server"
	"gorm.io/gorm"
)

const BASE_URL = "https://office.warpdevelopment.com"

type User struct {
	gorm.Model
	Email         string `gorm:"uniqueIndex"`
	Token         string
	PersonId      int
	InitializedAt time.Time `gorm:"not null"`
}

type Database interface {
	CreateUser(user *User) error
	GetUserByEmail(email string) (*User, error)
	UpdateUserToken(email, token string) error
	GetUserByID(id uint) (*User, error)
}

type GormDatabase struct {
	db *gorm.DB
}

func NewGormDatabase(db *gorm.DB) *GormDatabase {
	return &GormDatabase{db: db}
}

func (g *GormDatabase) CreateUser(user *User) error {
	return g.db.Create(user).Error
}

func (g *GormDatabase) GetUserByEmail(email string) (*User, error) {
	var user User
	err := g.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (g *GormDatabase) UpdateUserToken(email, token string) error {
	return g.db.Model(&User{}).Where("email = ?").Update("token", token).Error
}

func (g *GormDatabase) GetUserByID(id uint) (*User, error) {
	var user User
	err := g.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

type Handler struct {
	db Database
}

func NewHandler(db Database) *Handler {
	return &Handler{db: db}
}

func main() {
	exePath, err := os.Executable()
	if err != nil {
		panic(fmt.Sprintf("Failed to get executable path: %v\n", err))
	}
	exeDir := filepath.Dir(exePath)
	dbPath := filepath.Join(exeDir, "timesheets.db")

	gormDB, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		panic("Database failed to initialize")
	}

	err = gormDB.AutoMigrate(&User{})
	if err != nil {
		panic("Error running migration")
	}

	db := NewGormDatabase(gormDB)

	handler := NewHandler(db)

	handler.updateUserToken("", "")

	s := server.NewMCPServer(
		"Demo 🚀",
		"1.0.0",
		server.WithToolCapabilities(false),
	)

	// authTool := mcp.NewTool("Get Timesheet Token",
	// 	mcp.WithDescription("Get the Token from Timesheets Endpoint Using Username and Password"),
	// )

	// s.AddTool(authTool, handler.updateUserToken)

	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}

func (h *Handler) updateUserToken(email string, password string) (token string, err error) {
	// email, exists := os.LookupEnv("EMAIL")
	// if !exists {
	// 	return mcp.NewToolResultError("Email does not exist in environment."), fmt.Errorf("email does not exist in environment")
	// }

	// password, exists := os.LookupEnv("PASSWORD")
	// if !exists {
	// 	return mcp.NewToolResultError("Password does not exist in environment"), fmt.Errorf("password does not exist in environment")
	// }

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
		BASE_URL+"/api/account/authorise",
		"application/json",
		bytes.NewBuffer(jsonData),
	)

	if err != nil {
		return "", fmt.Errorf("error making HTTP request to %s: %v", BASE_URL, err)
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

	err = h.db.UpdateUserToken(email, responseData.Token)
	if err != nil {
		return "", fmt.Errorf("error updating user token: %v", err)
	}

	return responseData.Token, nil
}
