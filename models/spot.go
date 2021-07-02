package models

import (
	"context"
	"database/sql"
	"minesweeper/dbhandler"
)

type Spot struct {
	SpotId    int64
	GameId    int64
	Value     string
	X         int
	Y         int
	NearSpots map[string]Spot
	Status    string
}

func (s *Spot) Insert(db *dbhandler.DbHandler, tx *sql.Tx, ctx *context.Context) error {
	args := make([]interface{}, 0)
	args = append(args, s.GameId)
	args = append(args, s.Value)
	args = append(args, s.X)
	args = append(args, s.Y)

	nearSpots := ""

	for key := range s.NearSpots {
		nearSpots += key + "|"
	}

	args = append(args, nearSpots)
	args = append(args, s.Status)

	id, err := db.ExecuteTransaction(dbhandler.CREATE_SPOT, args, tx, ctx)

	if err != nil {
		return err
	}

	s.SpotId = id

	return nil
}
