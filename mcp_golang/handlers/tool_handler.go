package handlers

type ToolHandler struct {
	sh *DataHandler
}

func NewToolHandler(sh *DataHandler) *ToolHandler {
	return &ToolHandler{sh: sh}
}
