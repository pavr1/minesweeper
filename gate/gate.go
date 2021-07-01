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

func (g *Gate) CreateUser(name, lastName, password string) error {
	if name == "" {
		return fmt.Errorf("User name required")
	}

	if lastName == "" {
		return fmt.Errorf("User last name required")
	}

	if password == "" {
		return fmt.Errorf("Password required")
	}

	u, err := models.CreateUser(name, lastName, password)

	if err != nil {
		return err
	}

	err = u.Insert(g.DbHandler)

	if err != nil {
		return err
	}

	return nil
}
