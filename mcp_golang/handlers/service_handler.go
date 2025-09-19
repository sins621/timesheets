package handlers

import (
	"fmt"
	"time"

	"ts_mcp/interfaces"
	"ts_mcp/models"
	"ts_mcp/types"
)

type ServiceHandler struct {
	db interfaces.Database
	r  interfaces.Request
}

func NewServiceHandler(db interfaces.Database, r interfaces.Request) *ServiceHandler {
	return &ServiceHandler{db: db, r: r}
}

func (sh *ServiceHandler) logWorkService(email string, password string, workEntry types.WorkEntry) (err error) {
	user, err := sh.db.SelectUserByEmail(email)

	if err != nil {
		token, err := sh.r.GetUserToken(email, password)
		if err != nil {
			return fmt.Errorf("error getting user token for email: %s\n, err: %v\n", email, err)
		}

		personID, err := sh.r.GetPersonID(token)
		if err != nil {
			return fmt.Errorf("error getting user person id for email: %s\n, err: %v\n", email, err)
		}

		user, err = sh.db.CreateUser(&models.User{
			Email:         email,
			Token:         token,
			PersonID:      personID,
			InitializedAt: time.Now(),
		})
	} else if time.Since(user.InitializedAt) > time.Hour*24*7 {
		token, err := sh.r.GetUserToken(email, password)
		if err != nil {
			return fmt.Errorf("error getting user token for email: %s\n, err: %v\n", email, err)
		}

		user.Token = token
		user, err = sh.db.UpdateUser(user)
	}

	err = sh.r.PostTimeSheetEntry(user.Token, user.PersonID, workEntry)
	return err
}
