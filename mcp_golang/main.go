package main

import (
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"ts_mcp/database"
	"ts_mcp/handlers"
	"ts_mcp/request"
)

func main() {
	gormDB := database.InitializeGormDB()
	db := database.NewGormDatabase(gormDB)
	tsr := request.NewTimeSheetRequest("https://office.warpdevelopment.com")
	serviceHandler := handlers.NewServiceHandler(db, tsr)
	toolHandler := handlers.NewToolHandler(serviceHandler)

	s := server.NewMCPServer(
		"Timesheet MCP",
		"0.0.1",
		server.WithToolCapabilities(false),
	)

	tool := mcp.NewTool("test",
		mcp.WithDescription("test"),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("test"),
		),
	)

	s.AddTool(tool, toolHandler.LogWork)

	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
