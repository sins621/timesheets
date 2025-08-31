package main

import (
	"context"
	"fmt"
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

	// db := NewGormDatabase(gormDB)

	// handler := NewHandler(db)

	s := server.NewMCPServer(
		"Demo 🚀",
		"1.0.0",
		server.WithToolCapabilities(false),
	)

	// authTool := mcp.NewTool("Get Timesheet Token",
	// 	mcp.WithDescription("Get the Token from Timesheets Endpoint Using Username and Password"),
	// )

	// s.AddTool(authTool, handler.authHandler)

	nameTool := mcp.NewTool("get_name",
		mcp.WithDescription("Get the name to enter into the Greeting Tool"),
	)

	s.AddTool(nameTool, nameHandler)

	greetingTool := mcp.NewTool("greet",
		mcp.WithDescription("Greet today's name"),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("Name for today which can be gotten from get_name tool"),
		),
	)

	s.AddTool(greetingTool, greetingHandler)

	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}

func nameHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return mcp.NewToolResultText("John Doe"), nil
}

func greetingHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, err := request.RequireString("name")

	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("You did not provide the name for this tool: %v\n", err)), err
	}

	return mcp.NewToolResultText(fmt.Sprintf("Hello %s\n", name)), nil
}

// func (h *Handler) authHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
// 	email, exists := os.LookupEnv("EMAIL")
// 	if !exists {
// 		return mcp.NewToolResultError("Email does not exist in environment."), fmt.Errorf("email does not exist in environment")
// 	}

// 	password, exists := os.LookupEnv("PASSWORD")
// 	if !exists {
// 		return mcp.NewToolResultError("Password does not exist in environment"), fmt.Errorf("password does not exist in environment")
// 	}

// 	type RequestBody struct {
// 		Email    string `json:"Email"`
// 		Password string `json:"Password"`
// 	}

// 	type ResponseBody struct {
// 		Token string `json:"token"`
// 	}

// 	requestData := RequestBody{
// 		Email:    email,
// 		Password: password,
// 	}

// 	jsonData, err := json.Marshal(requestData)

// 	if err != nil {
// 		return mcp.NewToolResultError(fmt.Sprintf("Error marshaling JSON: %v\n", err)), err
// 	}

// 	resp, err := http.Post(
// 		BASE_URL+"/api/account/authorise",
// 		"application/json",
// 		bytes.NewBuffer(jsonData),
// 	)

// 	if err != nil {
// 		return mcp.NewToolResultError(fmt.Sprintf("Error making HTTP Request to %s: %v\n", BASE_URL, err)), err
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK {
// 		return mcp.NewToolResultError(fmt.Sprintf("Request Failed with status: %d\n", resp.StatusCode)), fmt.Errorf("request failed with status: %d", resp.StatusCode)
// 	}

// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		return mcp.NewToolResultError(fmt.Sprintf("Error reading response: %v\n", err)), err
// 	}

// 	var responseData ResponseBody
// 	err = json.Unmarshal(body, &responseData)

// 	if err != nil {
// 		return mcp.NewToolResultError(fmt.Sprintf("Error parsing json response: %v\n", err)), err
// 	}

// 	err = h.db.UpdateUserToken(email, responseData.Token)
// 	if err != nil {
// 		return mcp.NewToolResultError(fmt.Sprintf("Error updating user token: %v\n", err)), err
// 	}

// 	return mcp.NewToolResultText(fmt.Sprintf("The authorizatin token is: %s\n", responseData.Token)), nil
// }
