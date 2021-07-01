package gate

import (
	"fmt"
	"math/rand"
	"minesweeper/ccache"
	"minesweeper/dbhandler"
	"minesweeper/models"
	"strconv"
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
		return fmt.Errorf("User name required")
	}

	if user.LastName == "" {
		return fmt.Errorf("User last name required")
	}

	if user.Password == "" {
		return fmt.Errorf("Password required")
	}

	err := user.CreateUser(g.DbHandler)

	if err != nil {
		return err
	}

	return nil
}

func (g *Gate) ValidateLogin(user models.User) (int, error) {
	if user.Name == "" {
		return -1, fmt.Errorf("User name required")
	}

	if user.Password == "" {
		return -1, fmt.Errorf("Password required")
	}

	userId, err := user.ValidateUser(g.DbHandler)

	if err != nil {
		return -1, err
	}

	return userId, nil
}

func (g *Gate) CreateGame(game models.Game) error {
	if game.UserId == 0 {
		return fmt.Errorf("User id required")
	}

	err := game.Create(g.DbHandler)

	if err != nil {
		return err
	}

	g.GenerateGrid(&game)

	if err != nil {
		return err
	}

	return nil
}

func (g *Gate) GenerateGrid(game *models.Game) error {
	rows := game.Rows
	columns := game.Columns
	mines := game.Mines

	totalSpots := rows * columns

	if mines >= totalSpots {
		return fmt.Errorf("Too many mines, please decrease mine amount!")
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

	//var wg sync.WaitGroup
	for _, spot := range spots {
		// wg.Add(1)
		// go func(spot *models.Spot, spots map[string]*models.Spot, wg *sync.WaitGroup) {
		spot.NearSpots = LoadNearSpots(spot.X, spot.Y, spots)
		//wg.Done()
		// }(value, spots, &wg)
	}

	//wg.Wait()

	return nil
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

func LoadNearSpots(x, y int, spots map[string]*models.Spot) map[string]*models.Spot {
	const maxIndex = 3
	auxX := x - 1
	auxY := y - 1

	nearSpots := make(map[string]*models.Spot)

	for {
		if auxX < 0 || auxY < 0 {
			auxY++

			if auxY >= maxIndex {
				auxY = y - 1
				auxX++
			}

			continue
		}

		if auxX != x || auxY != y {
			id := strconv.Itoa(auxX) + "," + strconv.Itoa(auxY)

			nearSpots[id] = spots[id]
		}

		auxY++

		if auxY >= maxIndex {
			auxY = y - 1
			auxX++
		}

		if auxX >= maxIndex {
			break
		}
	}

	return nearSpots
}
