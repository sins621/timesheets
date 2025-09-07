package interfaces

import "main.go/models"

type Database interface {
	CreateUser(user *models.User) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	UpdateUser(user *models.User) (*models.User, error)
}

type Request interface {
	RequestUserToken(email string, password string) (token string, err error)
	RequestPersonID(token string) (id int, err error)
}