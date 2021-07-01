package models

import (
	"minesweeper/dbhandler"
	"time"
)

type User struct {
	UserId      int64
	Name        string
	LastName    string
	Password    string
	CreatedDate time.Time
	Message     string
}

func (u *User) Insert(db *dbhandler.DbHandler) error {
	args := make([]interface{}, 3)
	args = append(args, u.Name)
	args = append(args, u.LastName)
	args = append(args, u.Password)

	id, err := db.Execute(dbhandler.INSERT_USER, args)

	if err != nil {
		return err
	}

	u.UserId = id

	return nil
}

func (u *User) ValidateUser(db *dbhandler.DbHandler) (int, error) {
	args := []string{u.Name, u.Password}

	result, err := db.Select(dbhandler.VALIDATE_LOGIN, "User", args)

	if err != nil {
		return -1, err
	}

	if result == nil {
		return -1, nil
	}

	if len(result) == 0 {
		return -1, nil
	}

	return result[0].(dbhandler.DbUser).UserId, nil
}
