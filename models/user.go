package models

import (
	"database/sql"
	"fmt"
	"minesweeper/dbhandler"
	"minesweeper/guid"
	"time"
)

type User struct {
	UserId      guid.Guid
	Name        string
	LastName    string
	CreatedDate time.Time
}

func CreateUser(name, lastName string) (*User, error) {
	if name == "" {
		return nil, fmt.Errorf("User name required")
	}

	if lastName == "" {
		return nil, fmt.Errorf("User last name required")
	}

	user := User{
		Name:     name,
		LastName: lastName,
	}

	return &user, nil
}

func (u *User) Insert(db *dbhandler.DbHandler) error {
	var args []sql.NamedArg

	args[0] = sql.Named("UserId", *guid.New())
	args[1] = sql.Named("UserName", u.Name)
	args[2] = sql.Named("UserLastName", u.LastName)

	return db.Execute("INSERT INTO [dbo].[User] ([UserId],[UserName],[UserLastName],[CreatedDate]) VALUES (@UserId, @UserName, @UserLastName, GetDate()", args)
}
