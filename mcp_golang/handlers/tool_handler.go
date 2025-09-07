package handlers

type ToolHandler struct {
	sh *ServiceHandler
}

func NewToolHandler(sh *ServiceHandler) *ToolHandler {
	return &ToolHandler{sh: sh}
}
