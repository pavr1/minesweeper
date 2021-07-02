package models

import (
	"context"
	"database/sql"
	"minesweeper/dbhandler"
	"time"
)

type Game struct {
	UserId       int64
	GameId       int64
	CreatedDate  time.Time
	TimeConsumed float32
	Status       string
	Rows         int
	Columns      int
	Mines        int
	Message      string
}

func (g *Game) Create(db *dbhandler.DbHandler, tx *sql.Tx, ctx *context.Context) error {
	args := make([]interface{}, 0)
	args = append(args, g.UserId)
	args = append(args, g.TimeConsumed)

	id, err := db.ExecuteTransaction(dbhandler.CREATE_GAME, args, tx, ctx)

	if err != nil {
		return err
	}

	g.GameId = id

	return nil
}