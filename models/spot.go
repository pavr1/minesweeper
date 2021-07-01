package models

import "minesweeper/dbhandler"

type Spot struct {
	SpotId    int64
	GameId    int64
	Value     string
	X         int
	Y         int
	NearSpots map[string]*Spot
	Status    string
}

func (s *Spot) Insert(db *dbhandler.DbHandler) error {
	args := make([]interface{}, 3)
	args = append(args, s.GameId)
	args = append(args, s.Value)
	args = append(args, s.X)
	args = append(args, s.Y)

	nearSpots := ""

	for key, _ := range s.NearSpots {
		nearSpots += key + "|"
	}

	args = append(args, nearSpots)
	args = append(args, s.Status)

	id, err := db.Execute(dbhandler.CREATE_GAME, args)

	if err != nil {
		return err
	}

	s.SpotId = id

	return nil
}
