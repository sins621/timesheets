package handlers

import "ts_mcp/interfaces"

type ServiceHandler struct {
	db interfaces.Database
	r  interfaces.Request
}

func NewServiceHandler(db interfaces.Database, r interfaces.Request) *ServiceHandler {
	return &ServiceHandler{db: db, r: r}
}
