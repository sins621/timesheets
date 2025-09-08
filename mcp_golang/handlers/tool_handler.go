package handlers

import (
	"context"
	"time"
	"ts_mcp/constants"

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

	date := r.GetString(constants.ParamDate, t.Format("2006-01-02T15:04:05"))

	hours, err := r.RequireInt(constants.ParamHours)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	projectID, err := r.RequireString(constants.ParamProjectID)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	costCodeID, err := r.RequireString(constants.ParamCostCodeID)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText("Working"), nil
}
