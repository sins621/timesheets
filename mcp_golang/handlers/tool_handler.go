package handlers

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
)


type ToolHandler struct {
	sh *ServiceHandler
}

func NewToolHandler(sh *ServiceHandler) *ToolHandler {
	return &ToolHandler{sh: sh}
}

func (th *ToolHandler) LogWork(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return mcp.NewToolResultText("Working"), nil
}
