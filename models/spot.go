package models

import (
	"database/sql"
	"fmt"
	"minesweeper/dbhandler"
	"strconv"
)

type Spot struct {
	SpotId    int64
	GameId    int64
	Value     string
	X         int
	Y         int
	NearSpots map[string]*Spot
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

func GetSpotsByGameId(handler *dbhandler.DbHandler, gameId int64) (map[string]*Spot, error) {
	args := make([]interface{}, 0)
	args = append(args, gameId)

	r, err := handler.Select(dbhandler.SELECT_SPOTS_BY_GAME_ID, "Spot", args)

	if err != nil {
		return nil, err
	}

	spotsMap := make(map[string]*Spot)

	for _, s := range r {
		nearSpots := make(map[string]*Spot)
		dbspot := s.(dbhandler.DbSpot)
		id := strconv.Itoa(dbspot.X) + "," + strconv.Itoa(dbspot.Y)

		spot := Spot{
			SpotId:    dbspot.SpotId,
			GameId:    dbspot.GameId,
			Value:     dbspot.Value,
			X:         dbspot.X,
			Y:         dbspot.Y,
			NearSpots: nearSpots,
			Status:    dbspot.Status,
		}

		spotsMap[id] = &spot
	}

	return spotsMap, nil
}

func GetSpotById(handler *dbhandler.DbHandler, spotId int64) (*Spot, error) {
	args := make([]interface{}, 0)
	args = append(args, spotId)

	r, err := handler.Select(dbhandler.SELECT_SPOT_BY_ID, "Spot", args)

	if err != nil {
		return nil, err
	}

	var spot *Spot
	for _, s := range r {
		nearSpots := make(map[string]*Spot)
		dbspot := s.(dbhandler.DbSpot)
		spot = &Spot{
			SpotId:    dbspot.SpotId,
			GameId:    dbspot.GameId,
			Value:     dbspot.Value,
			X:         dbspot.X,
			Y:         dbspot.Y,
			NearSpots: nearSpots,
			Status:    dbspot.Status,
		}

		break
	}

	return spot, nil
}

func (s *Spot) GetMinesAround() string {
	if s.Value == "&#128163" {
		return s.Value
	}

	amountOfNearMines := 0

	for _, spot := range s.NearSpots {
		if spot.Value == "&#128163" {
			amountOfNearMines++
		}
	}

	if amountOfNearMines == 0 {
		return ""
	} else {
		return strconv.Itoa(amountOfNearMines)
	}
}

func (s *Spot) ProcessSpot(handler *dbhandler.DbHandler, rows, colums int, status string) error {
	args := make([]interface{}, 0)

	if status == "Open" {
		args = append(args, "Open")
		args = append(args, s.SpotId)

		_, err := handler.Execute(dbhandler.UPDATE_SPOT_STATUS, args)

		s.Status = "Open"

		if err != nil {
			return err
		}

		if len(s.NearSpots) == 0 {
			spots, err := GetSpotsByGameId(handler, s.GameId)

			if err != nil {
				return err
			}

			s.NearSpots = spots
			s.GetNearSpots(rows, colums, spots)
		}

		if s.Value == "" {
			for _, spot := range s.NearSpots {
				args := make([]interface{}, 0)
				args = append(args, "Open")
				args = append(args, spot.SpotId)

				_, err := handler.Execute(dbhandler.UPDATE_SPOT_STATUS, args)

				if err != nil {
					return err
				}
			}
		} else if s.Value == "&#128163" {
			args := make([]interface{}, 0)
			args = append(args, "E")
			args = append(args, s.SpotId)

			_, err := handler.Execute(dbhandler.UPDATE_SPOT_VALUE, args)
			if err != nil {
				return err
			}

			//game over
			spots, err := GetSpotsByGameId(handler, s.GameId)
			if err != nil {
				return err
			}

			for _, spot := range spots {
				args := make([]interface{}, 0)
				status := "Open"

				if spot.Status == "RedFlagged" {
					status = "ROpen"
				} else if spot.Status == "QuestionFlagged" {
					status = "QOpen"
				}

				args = append(args, status)
				args = append(args, spot.SpotId)

				_, err := handler.Execute(dbhandler.UPDATE_SPOT_STATUS, args)

				if err != nil {
					return err
				}
			}

			return fmt.Errorf("Game Over!")
		}
	} else {
		args = append(args, status)
		args = append(args, s.SpotId)

		_, err := handler.Execute(dbhandler.UPDATE_SPOT_STATUS, args)
		if err != nil {
			return err
		}

		return nil
	}

	return nil
}

func (s *Spot) GetNearSpots(rows, colums int, spots map[string]*Spot) {
	nearSpots := make(map[string]*Spot)

	var id string
	auxX := s.X - 1
	auxY := s.Y - 1
	if auxX >= 0 && auxY >= 0 {
		id = strconv.Itoa(auxX) + "," + strconv.Itoa(auxY)
		nearSpots[id] = spots[id]
	}

	auxX = s.X
	auxY = s.Y - 1
	if auxY >= 0 {
		id = strconv.Itoa(auxX) + "," + strconv.Itoa(auxY)
		nearSpots[id] = spots[id]
	}

	auxX = s.X + 1
	auxY = s.Y - 1
	if auxX < rows && auxY >= 0 {
		id = strconv.Itoa(auxX) + "," + strconv.Itoa(auxY)
		nearSpots[id] = spots[id]
	}

	//---

	auxX = s.X - 1
	auxY = s.Y
	if auxX >= 0 {
		id = strconv.Itoa(auxX) + "," + strconv.Itoa(auxY)
		nearSpots[id] = spots[id]
	}

	auxX = s.X + 1
	auxY = s.Y
	if auxX < colums {
		id = strconv.Itoa(auxX) + "," + strconv.Itoa(auxY)
		nearSpots[id] = spots[id]
	}
	//--

	auxX = s.X - 1
	auxY = s.Y + 1
	if auxX >= 0 && auxY < rows {
		id = strconv.Itoa(auxX) + "," + strconv.Itoa(auxY)
		nearSpots[id] = spots[id]
	}

	auxX = s.X
	auxY = s.Y + 1
	if auxY < rows {
		id = strconv.Itoa(auxX) + "," + strconv.Itoa(auxY)
		nearSpots[id] = spots[id]
	}

	auxX = s.X + 1
	auxY = s.Y + 1
	if auxX < colums && auxY < rows {
		id = strconv.Itoa(auxX) + "," + strconv.Itoa(auxY)
		nearSpots[id] = spots[id]
	}

	s.NearSpots = nearSpots
}
