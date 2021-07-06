package models

import (
	"fmt"
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

type UIUser struct {
	UserId       int64
	Name         string
	LastName     string
	Password     string
	CreatedDate  time.Time
	Message      string
	PendingGames []UIGame
}

func (u *User) CreateUser(handler *dbhandler.DbHandler) (*UIUser, error) {
	args := make([]interface{}, 0)
	args = append(args, u.Name)
	args = append(args, u.LastName)
	args = append(args, u.Password)

	_, err := handler.Execute(dbhandler.CREATE_USER, args)

	if err != nil {
		return nil, err
	}

	uiuser, err := u.ValidateUser(handler)

	if err != nil {
		return nil, err
	}

	return uiuser, nil
}

func (u *User) ValidateUser(handler *dbhandler.DbHandler) (*UIUser, error) {
	args := make([]interface{}, 0)
	args = append(args, u.Name)
	args = append(args, u.Password)

	result, err := handler.Select(dbhandler.VALIDATE_LOGIN, "User", args)

	if err != nil {
		return nil, err
	}

	if result == nil || len(result) == 0 {
		return nil, fmt.Errorf("User or password invalid")
	}

	dbUser := result[0].(dbhandler.DbUser)

	user := UIUser{
		UserId:      dbUser.UserId,
		Name:        dbUser.Name,
		LastName:    dbUser.LastName,
		Password:    dbUser.Password,
		CreatedDate: dbUser.CreatedDate,
	}

	pendingGames, err := GetPendingGames(handler, user.UserId)
	if err != nil {
		return nil, err
	}

	uigames := make([]UIGame, len(pendingGames))
	for i, game := range pendingGames {
		uigames[i] = UIGame{
			GameId:       game.GameId,
			UserId:       game.UserId,
			CreatedDate:  game.CreatedDate,
			TimeConsumed: game.TimeConsumed,
			Status:       game.Status,
			Rows:         game.Rows,
			Columns:      game.Columns,
			Mines:        game.Mines,
		}
	}

	user.PendingGames = uigames

	return &user, nil
}
