package multitool

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {

	s := server.NewMCPServer(
		"Demo 🚀",
		"1.0.0",
		server.WithToolCapabilities(false),
	)

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
