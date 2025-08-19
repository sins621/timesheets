package main

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	mcpServer := server.NewMCPServer(
		"HTTP Demo",
		"1.0.0",
		server.WithToolCapabilities(false),
		server.WithRecovery(),
	)

	httpTool := mcp.NewTool("http request",
		mcp.WithDescription("Perform a basic http request"),
		mcp.WithString("address",
			mcp.Required(),
			mcp.Description("The address to make the get request to"),
		),
	)

	mcpServer.AddTool(httpTool, httpToolHandler)

	if err := server.ServeStdio(mcpServer); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}

func httpToolHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	url, err := request.RequireString("address")

	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	res, err := requestHandler(url)

	if err != nil {
		return mcp.NewToolResultError("Error making http request"), nil
	}

	return mcp.NewToolResultText(res), nil
}

func requestHandler(requestUrl string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, requestUrl, nil)
	if err != nil {
		return "could not make request", nil
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "error making http request", nil
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Sprintf("could not read response body: %s\n", err), nil
	}

	return fmt.Sprintf("response body: %s\n", resBody), nil
}
