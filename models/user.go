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
}

func CreateUser(name, lastName, password string) (*User, error) {
	user := User{
		Name:     name,
		LastName: lastName,
		Password: password,
	}

	return &user, nil
}

func (u *User) Insert(db *dbhandler.DbHandler) error {
	args := []string{u.Name, u.LastName, u.Password}

	return db.Execute(dbhandler.INSERT_USER, args)
}
