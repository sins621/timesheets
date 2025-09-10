package interfaces

import (
	"time"
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
	PostTimeSheetEntry(token string, taskID int, personID int, costCodeID int, overtime bool, time int, date time.Time, description string) (err error)
}
