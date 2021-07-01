package models

import (
	"fmt"
	"minesweeper/dbhandler"
	"strconv"
	"time"
)

type Game struct {
	UserId       int
	GameId       int
	CreatedDate  time.Time
	TimeConsumed float32
	Status       string
	Rows         int
	Columns      int
	Mines        int
}

func (g *Game) Create(db *dbhandler.DbHandler) error {
	args := []string{strconv.Itoa(g.UserId), g.CreatedDate.String(), fmt.Sprintf("%f", g.TimeConsumed), "P"}

	return db.Execute(dbhandler.CREATE_GAME, args)
}
