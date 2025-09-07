package handlers

import "main.go/interfaces"

type DataHandler struct {
	db interfaces.Database
	r  interfaces.Request
}

func NewDataHandler(db interfaces.Database, r interfaces.Request) *DataHandler {
	return &DataHandler{db: db, r: r}
}
