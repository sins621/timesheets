package handlers

import (
	"context"
	"strconv"
	"strings"
	"time"

	"ts_mcp/constants"
	"ts_mcp/types"

	"github.com/mark3labs/mcp-go/mcp"
)

type McpUser struct {
	Email string
	Password string
}

type ToolHandler struct {
	sh *ServiceHandler
	user McpUser
}

func NewToolHandler(sh *ServiceHandler, user McpUser) *ToolHandler {
	return &ToolHandler{sh: sh, user: user}
}

func (th *ToolHandler) LogWork(ctx context.Context, r mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	description, err := r.RequireString("description")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	t := time.Now()

	stringDate := r.GetString(constants.ParamDate, t.Format(constants.TimeFormat))

	date, err := time.Parse(constants.TimeFormat, stringDate)

	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	hours, err := r.RequireInt(constants.ParamHours)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	taskId, err := r.RequireInt(constants.ParamTaskID)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	costCodeIDString, err := r.RequireString(constants.ParamCostCodeID)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	costCodeID, err := strconv.Atoi(strings.Split(costCodeIDString, ".")[0])

	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	workEntry := types.WorkEntry{
		Description: description,
		Date:        date,
		Hours:       hours,
		TaskID:      taskID,
		CostCodeID:  costCodeID,
	}

	err = th.sh.logWorkService(th.user.Email, th.user.Password, workEntry)

	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText("Working"), nil
}

func (th *ToolHandler) GetCostCodeIDs(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	    return []mcp.ResourceContents{
	        mcp.TextResourceContents{
	            URI:      "info://CostCodeIDs",
	            MIMEType: "text/plain",
	            Text:     strings.Join(constants.CostCodeIDs, "\n"),
	        },
	    }, nil
	}
