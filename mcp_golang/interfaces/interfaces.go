package interfaces

import (
	"ts_mcp/models"
)

type Database interface {
	CreateUser(user *models.User) (*models.User, error)
	SelectUserByEmail(email string) (*models.User, error)
	UpdateUser(user *models.User) (*models.User, error)
}

type Request interface {
	GetUserToken(email string, password string) (token string, err error)
	GetPersonID(token string) (id int, err error)
	PostTimeSheetEntry(token string, personID int, workEntry models.WorkEntry) (err error)
}
