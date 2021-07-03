package models

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"
	"minesweeper/dbhandler"
	"strconv"
	"sync"
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
	Spots        map[string]Spot
	Message      string
}

func (g *Game) Create(db *dbhandler.DbHandler, tx *sql.Tx, ctx *context.Context) error {
	args := make([]interface{}, 0)
	args = append(args, g.UserId)
	args = append(args, g.TimeConsumed)
	args = append(args, g.Status)
	args = append(args, g.Rows)
	args = append(args, g.Columns)
	args = append(args, g.Mines)

	id, err := db.ExecuteTransaction(dbhandler.CREATE_GAME, args, tx, ctx)

	if err != nil {
		return err
	}

	g.GameId = id

	return nil
}

func (g *Game) GenerateGrid() (*map[string]Spot, error) {
	rows := g.Rows
	columns := g.Columns
	mines := g.Mines

	totalSpots := rows * columns

	if mines >= totalSpots {
		return nil, fmt.Errorf("too many mines, please decrease mine amount")
	}

	spots := make(map[string]Spot)
	x := 0
	y := 0

	for {
		spot := Spot{
			GameId:    g.GameId,
			Value:     "", //pending for now. We need to loop through all near spots to count
			X:         x,
			Y:         y,
			NearSpots: make(map[string]Spot), //pending, loop throug all near spots to get ids
			Status:    "Closed",
		}

		spots[strconv.Itoa(x)+","+strconv.Itoa(y)] = spot

		y++

		if y >= columns {
			y = 0
			x++
		}

		if x >= rows {
			break
		}
	}

	g.setupMines(rows, columns, mines, &spots)

	var wg sync.WaitGroup
	for _, value := range spots {
		wg.Add(1)
		go func(spot Spot, spots map[string]Spot, wg *sync.WaitGroup) {
			spot.NearSpots = loadNearSpots(spot.X, spot.Y, rows, columns, spots)
			wg.Done()
		}(value, spots, &wg)
	}
	wg.Wait()

	return &spots, nil
}

func (g *Game) setupMines(rows, coulmns, mines int, spots *map[string]Spot) {
	for {
		randX := rand.Intn(rows)
		randY := rand.Intn(coulmns)

		id := strconv.Itoa(randX) + "," + strconv.Itoa(randY)

		val := *spots
		spot := val[id]

		if spot.Value == "" {
			spot.Value = "M"

			val[id] = spot

			mines--
		}

		if mines <= 0 {
			break
		}
	}
}

func loadNearSpots(x, y, rows, colums int, spots map[string]Spot) map[string]Spot {
	nearSpots := make(map[string]Spot)

	var id string
	auxX := x - 1
	auxY := y - 1
	if auxX >= 0 && auxY >= 0 {
		id = strconv.Itoa(auxX) + "," + strconv.Itoa(auxY)
		nearSpots[id] = spots[id]
	}

	auxX = x
	auxY = y - 1
	if auxY >= 0 {
		id = strconv.Itoa(auxX) + "," + strconv.Itoa(auxY)
		nearSpots[id] = spots[id]
	}

	auxX = x + 1
	auxY = y - 1
	if auxX < rows && auxY >= 0 {
		id = strconv.Itoa(auxX) + "," + strconv.Itoa(auxY)
		nearSpots[id] = spots[id]
	}

	//---

	auxX = x - 1
	auxY = y
	if auxX >= 0 {
		id = strconv.Itoa(auxX) + "," + strconv.Itoa(auxY)
		nearSpots[id] = spots[id]
	}

	auxX = x + 1
	auxY = y
	if auxX < colums {
		id = strconv.Itoa(auxX) + "," + strconv.Itoa(auxY)
		nearSpots[id] = spots[id]
	}
	//--

	auxX = x - 1
	auxY = y + 1
	if auxX >= 0 && auxY < rows {
		id = strconv.Itoa(auxX) + "," + strconv.Itoa(auxY)
		nearSpots[id] = spots[id]
	}

	auxX = x
	auxY = y + 1
	if auxY < rows {
		id = strconv.Itoa(auxX) + "," + strconv.Itoa(auxY)
		nearSpots[id] = spots[id]
	}

	auxX = x + 1
	auxY = y + 1
	if auxX < colums && auxY < rows {
		id = strconv.Itoa(auxX) + "," + strconv.Itoa(auxY)
		nearSpots[id] = spots[id]
	}

	return nearSpots
}

func GetPendingGames(userId int64, db *dbhandler.DbHandler) ([]Game, error) {
	args := make([]interface{}, 0)
	args = append(args, userId)

	results := make([]Game, 0)
	r, err := db.Select(dbhandler.SELECT_GAMES_BY_USER, "Game", args)

	if err != nil {
		return nil, err
	}

	for _, game := range r {
		dbgame := game.(dbhandler.DbGame)
		results = append(results, Game{
			GameId:       dbgame.GameId,
			UserId:       dbgame.UserId,
			CreatedDate:  dbgame.CreatedDate,
			TimeConsumed: dbgame.TimeConsumed,
			Status:       dbgame.Status,
			Rows:         dbgame.Rows,
			Columns:      dbgame.Columns,
			Mines:        dbgame.Mines,
		})
	}

	return results, nil
}

func GetSingleGame(gameId int64, db *dbhandler.DbHandler) (*Game, error) {
	args := make([]interface{}, 0)
	args = append(args, gameId)

	r, err := db.Select(dbhandler.SELECT_GAME_BY_ID, "Game", args)

	if err != nil {
		return nil, err
	}

	var game Game

	for _, g := range r {
		dbgame := g.(dbhandler.DbGame)
		game = Game{
			GameId:       dbgame.GameId,
			UserId:       dbgame.UserId,
			CreatedDate:  dbgame.CreatedDate,
			TimeConsumed: dbgame.TimeConsumed,
			Status:       dbgame.Status,
			Rows:         dbgame.Rows,
			Columns:      dbgame.Columns,
			Mines:        dbgame.Mines,
		}

		break
	}

	return &game, nil
}

func GetLatestGame(userId int64, db *dbhandler.DbHandler, tx *sql.Tx, ctx *context.Context) (*Game, error) {
	args := make([]interface{}, 0)
	args = append(args, userId)

	r, err := db.SelectTransaction(dbhandler.SELECT_LATEST_GAME, "Game", args, tx, ctx)

	if err != nil {
		return nil, err
	}

	var game Game

	for _, g := range r {
		dbgame := g.(dbhandler.DbGame)
		game = Game{
			GameId:       dbgame.GameId,
			UserId:       dbgame.UserId,
			CreatedDate:  dbgame.CreatedDate,
			TimeConsumed: dbgame.TimeConsumed,
			Status:       dbgame.Status,
			Rows:         dbgame.Rows,
			Columns:      dbgame.Columns,
			Mines:        dbgame.Mines,
		}

		break
	}

	return &game, nil
}
