package gate

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"
	"minesweeper/ccache"
	"minesweeper/dbhandler"
	"minesweeper/models"
	"strconv"
	"sync"
)

type Gate struct {
	DbHandler *dbhandler.DbHandler
	Cache     *ccache.CCache
}

func Start() (*Gate, error) {
	handler, err := dbhandler.GetInstance()

	if err != nil {
		return nil, err
	}

	gate := &Gate{
		DbHandler: handler,
		Cache:     ccache.New(),
	}

	return gate, nil
}

func (g *Gate) CreateUser(user models.User) error {
	if user.Name == "" {
		return fmt.Errorf("user name required")
	}

	if user.LastName == "" {
		return fmt.Errorf("user last name required")
	}

	if user.Password == "" {
		return fmt.Errorf("password required")
	}

	err := user.CreateUser(g.DbHandler)

	if err != nil {
		return err
	}

	return nil
}

func (g *Gate) ValidateLogin(user models.User) (int, error) {
	if user.Name == "" {
		return -1, fmt.Errorf("user name required")
	}

	if user.Password == "" {
		return -1, fmt.Errorf("password required")
	}

	userId, err := user.ValidateUser(g.DbHandler)

	if err != nil {
		return -1, err
	}

	return userId, nil
}

func (g *Gate) CreateGame(game models.Game) error {
	if game.UserId == 0 {
		return fmt.Errorf("user id required")
	}

	ctx := context.Background()
	tx, err := g.DbHandler.Db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	err = game.Create(g.DbHandler, tx, &ctx)

	if err != nil {
		tx.Rollback()
		return err
	}

	spots, err := g.GenerateGrid(&game)

	if err != nil {
		tx.Rollback()
		return err
	}

	err = g.InsertSpots(spots, tx, &ctx)

	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}

func (g *Gate) GenerateGrid(game *models.Game) (*map[string]*models.Spot, error) {
	rows := game.Rows
	columns := game.Columns
	mines := game.Mines

	totalSpots := rows * columns

	if mines >= totalSpots {
		return nil, fmt.Errorf("too many mines, please decrease mine amount")
	}

	spots := make(map[string]*models.Spot)
	x := 0
	y := 0

	for {
		spot := models.Spot{
			GameId:    game.GameId,
			Value:     "", //pending for now. We need to loop through all near spots to count
			X:         x,
			Y:         y,
			NearSpots: make(map[string]*models.Spot), //pending, loop throug all near spots to get ids
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

	setupMines(rows, columns, mines, &spots)

	var wg sync.WaitGroup
	for _, value := range spots {
		wg.Add(1)
		go func(spot *models.Spot, spots map[string]*models.Spot, wg *sync.WaitGroup) {
			spot.NearSpots = LoadNearSpots(spot.X, spot.Y, rows, columns, spots)
			wg.Done()
		}(value, spots, &wg)
	}
	wg.Wait()

	return &spots, nil
}

func setupMines(rows, coulmns, mines int, spots *map[string]*models.Spot) {
	for {
		randX := rand.Intn(rows)
		randY := rand.Intn(coulmns)

		id := strconv.Itoa(randX) + "," + strconv.Itoa(randY)

		val := *spots
		spot := val[id]

		if spot.Value == "" {
			spot.Value = "MINE"
			mines--
		}

		if mines <= 0 {
			break
		}
	}
}

func LoadNearSpots(x, y, rows, colums int, spots map[string]*models.Spot) map[string]*models.Spot {
	nearSpots := make(map[string]*models.Spot)

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

func (g *Gate) InsertSpots(spots *map[string]*models.Spot, tx *sql.Tx, ctx *context.Context) error {
	for _, spot := range *spots {
		err := spot.Insert(g.DbHandler, tx, ctx)

		if err != nil {
			return err
		}
	}

	return nil
}
