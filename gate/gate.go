package gate

import (
	"fmt"
	"minesweeper/dbhandler"
	"minesweeper/models"
)

type Gate struct {
	DbHandler *dbhandler.DbHandler
}

func Start() (*Gate, error) {
	handler, err := dbhandler.GetInstance()

	if err != nil {
		return nil, err
	}

	gate := &Gate{
		DbHandler: handler,
	}

	return gate, nil
}

func (g *Gate) CreateUser(user models.User) error {
	if user.Name == "" {
		return fmt.Errorf("User name required")
	}

	if user.LastName == "" {
		return fmt.Errorf("User last name required")
	}

	if user.Password == "" {
		return fmt.Errorf("Password required")
	}

	err := user.Insert(g.DbHandler)

	if err != nil {
		return err
	}

	return nil
}
