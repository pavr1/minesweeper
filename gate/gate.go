package gate

import (
	"minesweeper/dbhandler"
	"minesweeper/models"
)

var gate *Gate

type Gate struct {
	DbHandler *dbhandler.DbHandler
}

func Start() error {
	handler, err := dbhandler.GetInstance()

	if err != nil {
		return err
	}

	gate = &Gate{
		DbHandler: handler,
	}

	return nil
}

func CreateUser(name, lastName string) error {
	u, err := models.CreateUser(name, lastName)

	if err != nil {
		return err
	}

	gate.DbHandler.
}
