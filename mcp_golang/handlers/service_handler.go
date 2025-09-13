package handlers

import (
	"ts_mcp/interfaces"
	"ts_mcp/models"
)

type ServiceHandler struct {
	db interfaces.Database
	r  interfaces.Request
}

func NewServiceHandler(db interfaces.Database, r interfaces.Request) *ServiceHandler {
	return &ServiceHandler{db: db, r: r}
}

func (sh *ServiceHandler) logWorkService (email string, password string, workEntry models.WorkEntry) (err error) {
	err = sh.r.PostTimeSheetEntry("test", 0, workEntry)
	return err
}
