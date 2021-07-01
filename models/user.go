package models

import (
	"minesweeper/dbhandler"
	"time"

	uuid "github.com/nu7hatch/gouuid"
)

type User struct {
	UserId      uuid.UUID
	Name        string
	LastName    string
	Password    string
	CreatedDate time.Time
	Message     string
}

func (u *User) Insert(db *dbhandler.DbHandler) error {
	args := []string{u.Name, u.LastName, u.Password}

	return db.Execute(dbhandler.INSERT_USER, args)
}
