package models

import (
	"minesweeper/dbhandler"
	"time"
)

type User struct {
	UserId       int64
	Name         string
	LastName     string
	Password     string
	CreatedDate  time.Time
	Message      string
	PendingGames []Game
}

func (u *User) CreateUser(db *dbhandler.DbHandler) error {
	args := make([]interface{}, 0)
	args = append(args, u.Name)
	args = append(args, u.LastName)
	args = append(args, u.Password)

	id, err := db.Execute(dbhandler.CREATE_USER, args)

	if err != nil {
		return err
	}

	u.UserId = id

	return nil
}

func (u *User) ValidateUser(db *dbhandler.DbHandler) (*User, error) {
	args := make([]interface{}, 0)
	args = append(args, u.Name)
	args = append(args, u.Password)

	result, err := db.Select(dbhandler.VALIDATE_LOGIN, "User", args)

	if err != nil {
		return nil, err
	}

	if result == nil {
		return nil, nil
	}

	if len(result) == 0 {
		return nil, nil
	}

	dbUser := result[0].(dbhandler.DbUser)

	user := User{
		UserId:      dbUser.UserId,
		Name:        dbUser.Name,
		LastName:    dbUser.LastName,
		Password:    dbUser.Password,
		CreatedDate: dbUser.CreatedDate,
	}

	return &user, nil
}
