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
	Spots        map[string]*Spot
	Message      string
}

type UIGame struct {
	UserId       int64
	GameId       int64
	CreatedDate  time.Time
	TimeConsumed float32
	Status       string
	Rows         int
	Columns      int
	Mines        int
	Spots        map[string]*UISpot
	Message      string
}

func (g *Game) Create(ctx context.Context, handler *dbhandler.DbHandler) (*UIGame, error) {
	args := make([]interface{}, 0)
	args = append(args, g.UserId)
	args = append(args, g.TimeConsumed)
	args = append(args, g.Rows)
	args = append(args, g.Columns)
	args = append(args, g.Mines)

	conn, err := handler.DB.Conn(ctx)
	if err != nil {
		return nil, err
	}

	tx, err := conn.BeginTx(ctx, nil)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	_, err = handler.ExecuteTransaction(tx, dbhandler.CREATE_GAME, args)

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	game, err := GetLatestGame(tx, handler, g.UserId)

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	ui := UIGame{
		GameId:       game.GameId,
		UserId:       game.UserId,
		TimeConsumed: game.TimeConsumed,
		Status:       game.Status,
		Rows:         game.Rows,
		Columns:      game.Columns,
		Mines:        game.Mines,
	}

	spots, err := game.GenerateGrid()

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = insertSpots(tx, handler, spots)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	uiSpots, err := GetUISpotsByGameId(tx, handler, ui.GameId)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	ui.Spots = uiSpots
	tx.Commit()

	return &ui, nil
}

func insertSpots(tx *sql.Tx, handler *dbhandler.DbHandler, spots *map[string]*Spot) error {
	for key, spot := range *spots {
		err := spot.Insert(tx, handler)

		if err != nil {
			return err
		}

		s := *spots
		s[key] = spot
	}

	return nil
}

func (g *Game) GenerateGrid() (*map[string]*Spot, error) {
	rows := g.Rows
	columns := g.Columns
	mines := g.Mines

	totalSpots := rows * columns

	if mines >= totalSpots {
		return nil, fmt.Errorf("too many mines, please decrease mine amount")
	}

	spots := make(map[string]*Spot)
	x := 0
	y := 0

	for {
		nearSpots := make(map[string]*Spot)

		spot := Spot{
			GameId:    g.GameId,
			Value:     "",
			X:         x,
			Y:         y,
			NearSpots: nearSpots,
			Status:    "Closed",
		}

		spots[strconv.Itoa(x)+","+strconv.Itoa(y)] = &spot

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
		go func(spot *Spot, spots map[string]*Spot, wg *sync.WaitGroup) {
			spot.GetNearSpots(rows, columns, spots)
			spot.Value = spot.GetMinesAround()
			wg.Done()
		}(value, spots, &wg)
	}
	wg.Wait()

	return &spots, nil
}

func (g *Game) setupMines(rows, coulmns, mines int, spots *map[string]*Spot) {
	for {
		randX := rand.Intn(rows)
		randY := rand.Intn(coulmns)

		id := strconv.Itoa(randX) + "," + strconv.Itoa(randY)

		val := *spots
		spot := val[id]

		if spot.Value == "" {
			spot.Value = "&#128163"

			val[id] = spot

			mines--
		}

		if mines <= 0 {
			break
		}
	}
}

func GetPendingGames(handler *dbhandler.DbHandler, userId int64) ([]UIGame, error) {
	args := make([]interface{}, 0)
	args = append(args, userId)

	results := make([]UIGame, 0)
	r, err := handler.Select(dbhandler.SELECT_GAMES_BY_USER, "Game", args)

	if err != nil {
		return nil, err
	}

	for _, game := range r {
		dbgame := game.(dbhandler.DbGame)
		results = append(results, UIGame{
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

func GetSingleGame(handler *dbhandler.DbHandler, gameId int64) (*UIGame, error) {
	args := make([]interface{}, 0)
	args = append(args, gameId)

	var game = UIGame{}

	r, err := handler.Select(dbhandler.SELECT_GAME_BY_ID, "Game", args)

	if err != nil {
		return &game, err
	}

	for _, g := range r {
		dbgame := g.(dbhandler.DbGame)
		game = UIGame{
			GameId:       dbgame.GameId,
			UserId:       dbgame.UserId,
			CreatedDate:  dbgame.CreatedDate,
			TimeConsumed: dbgame.TimeConsumed,
			Status:       dbgame.Status,
			Rows:         dbgame.Rows,
			Columns:      dbgame.Columns,
			Mines:        dbgame.Mines,
		}

		uiSpots, err := GetUISpotsByGameId(nil, handler, game.GameId)
		if err != nil {
			return nil, err
		}

		game.Spots = uiSpots
		break
	}

	return &game, nil
}

func GetLatestGame(tx *sql.Tx, handler *dbhandler.DbHandler, userId int64) (*Game, error) {
	args := make([]interface{}, 0)
	args = append(args, userId)

	r, err := handler.SelectTransaction(tx, dbhandler.SELECT_LATEST_GAME, "Game", args)

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

func (g *Game) UpdateStatus(handler *dbhandler.DbHandler) (int64, error) {
	args := make([]interface{}, 0)
	args = append(args, g.Status)
	args = append(args, g.GameId)

	id, err := handler.Execute(dbhandler.UPDATE_GAME_STATUS, args)

	if err != nil {
		return -1, err
	}

	return id, nil
}
