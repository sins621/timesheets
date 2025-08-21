package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type User struct {
	Email    string
	Password string
}

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

	res, err := getWithToken(url)

	if err != nil {
		return mcp.NewToolResultError("Error making http request"), nil
	}

	return mcp.NewToolResultText(res), nil
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
