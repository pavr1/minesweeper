package models

import (
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

func (s *Spot) Insert(handler *dbhandler.DbHandler, tx *sql.Tx) error {
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

	id, err := handler.ExecuteTransaction(tx, dbhandler.CREATE_SPOT, args)

	if err != nil {
		return err
	}

	s.SpotId = id

	return nil
}

func GetSpotsByGameId(handler *dbhandler.DbHandler, gameId int64) ([]Spot, error) {
	args := make([]interface{}, 0)
	args = append(args, gameId)

	results := make([]Spot, 0)
	r, err := handler.Select(dbhandler.SELECT_SPOTS_BY_GAME_ID, "Spot", args)

	if err != nil {
		return nil, err
	}

	for _, spot := range r {
		dbspot := spot.(dbhandler.DbSpot)
		results = append(results, Spot{
			SpotId:    dbspot.SpotId,
			GameId:    dbspot.GameId,
			Value:     dbspot.Value,
			X:         dbspot.X,
			Y:         dbspot.Y,
			NearSpots: make(map[string]Spot),
			Status:    dbspot.Status,
		})
	}

	return results, nil
}
