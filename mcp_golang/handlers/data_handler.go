package handlers

import "main.go/interfaces"

type DataHandler struct {
	db interfaces.Database
	r  interfaces.Request
}

func NewDataHandler(db interfaces.Database, r interfaces.Request) *DataHandler {
	return &DataHandler{db: db, r: r}
}

type ToolHandler struct {
	sh *DataHandler
}

func NewToolHandler(sh *DataHandler) *ToolHandler {
	return &ToolHandler{sh: sh}
}
