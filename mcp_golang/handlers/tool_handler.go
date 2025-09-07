package handlers

import (
	"context"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
)

type ToolHandler struct {
	sh *ServiceHandler
}

func NewToolHandler(sh *ServiceHandler) *ToolHandler {
	return &ToolHandler{sh: sh}
}

func (th *ToolHandler) LogWork(ctx context.Context, r mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	description, err := r.RequireString("description")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	t := time.Now()

	date := r.GetString("date", t.Format("2006-01-02T15:04:05"))

	hours, err := r.RequireInt("hours")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	projectID, err := r.RequireString("projectID")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	costCodeID, err := r.RequireString("costCodeID")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText("Working"), nil
}
