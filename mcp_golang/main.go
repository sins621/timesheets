package main

import (
	"fmt"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"ts_mcp/constants"
	"ts_mcp/database"
	"ts_mcp/handlers"
	"ts_mcp/request"
)

func main() {
	email, present := os.LookupEnv("EMAIL")
	if !present {
		panic("email not present in environment")
	}

	password, present := os.LookupEnv("PASSWORD")
	if !present {
		panic("password not present in environment")
	}

	gormDB := database.InitializeGormDB()
	db := database.NewGormDatabase(gormDB)
	tsr := request.NewTimeSheetRequest("https://office.warpdevelopment.com")
	serviceHandler := handlers.NewServiceHandler(db, tsr)
	toolHandler := handlers.NewToolHandler(serviceHandler, handlers.McpUser{
		Email:    email,
		Password: password,
	})

	s := server.NewMCPServer(
		"Timesheet MCP",
		"0.0.1",
		server.WithToolCapabilities(true),
	)

	logTool := mcp.NewTool("log_tool",
		mcp.WithDescription("Allows for logging work done on timesheets."),
		mcp.WithString(constants.ParamDescription,
			mcp.Required(),
			mcp.Description("Detailed Description of work done, gather information from git history if necessary"),
		),
		mcp.WithString(constants.ParamDate,
			mcp.Description("Date of the entry with format 'YYYY-MM-DDTHH:MM:SS`, default to today's date if not provided"),
		),
		mcp.WithNumber(constants.ParamHours,
			mcp.Required(),
			mcp.Description("Hours the user has worked on this entry"),
		),
		mcp.WithNumber(constants.ParamTaskID,
			mcp.Required(),
			mcp.Description("The ID of the project the user is working on. Can be retrieved from Project ID Tool"),
		),
		mcp.WithString(constants.ParamCostCodeID,
			mcp.Required(),
			mcp.Description("The relevant Cost Code ID of the work done. Call the Get Cost Code IDs and related descriptions tool for this information and try to gather which codes to use from the context of the conversation. Only ask the user to provide cost code IDs explicitly if confidence is low on assumptions."),
			mcp.WithStringEnumItems(constants.CostCodeIDs),
		),
	)

	costCodeTool := mcp.NewTool("get_cost_code_tool",
		mcp.WithDescription("Allows to look up Cost Code IDs along with descriptions of them."))

	s.AddTool(logTool, toolHandler.LogWork)
	s.AddTool(costCodeTool, toolHandler.GetCostCodeIDs)

	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
