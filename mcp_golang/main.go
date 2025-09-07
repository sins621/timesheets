package main

import (
	"fmt"

	"github.com/mark3labs/mcp-go/server"

	"main.go/database"
	"main.go/handlers"
	"main.go/request"
)

func main() {
	gormDB := database.Init()
	db := database.NewGormDatabase(gormDB)
	tsr := request.NewTimeSheetRequest("https://office.warpdevelopment.com")
	dataHandler := handlers.NewDataHandler(db, tsr)
	toolHandler := NewToolHandler(dataHandler)

	s := server.NewMCPServer(
		"Demo 🚀",
		"1.0.0",
		server.WithToolCapabilities(false),
	)

	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
