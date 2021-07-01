package models

import (
	"minesweeper/dbhandler"
	"time"
)

type User struct {
	UserId      string
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

func (u *User) ValidateUser(db *dbhandler.DbHandler) (bool, error) {
	args := []string{u.Name, u.Password}

	result, err := db.Select(dbhandler.VALIDATE_LOGIN, "User", args)

	if err != nil {
		return false, err
	}

	if result == nil {
		return false, nil
	}

	return len(result) > 0, nil
}
