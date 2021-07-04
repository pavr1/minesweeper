package gate

import (
	"context"
	"database/sql"
	"fmt"
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
	handler, err := dbhandler.InitConnection()

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

func (g *Gate) ValidateLogin(user models.User) (*models.User, error) {
	if user.Name == "" {
		return nil, fmt.Errorf("user name required")
	}

	if user.Password == "" {
		return nil, fmt.Errorf("password required")
	}

	result, err := user.ValidateUser(g.DbHandler)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (g *Gate) CreateGame(game *models.Game) (*map[string]*models.Spot, error) {
	if game.UserId == 0 {
		return nil, fmt.Errorf("user id required")
	}

	gameId, err := game.Create(g.DbHandler)

	if err != nil {
		return nil, err
	}

	game.GameId = gameId
	spots, err := game.GenerateGrid()

	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	conn, err := g.DbHandler.DB.Conn(ctx)
	tx, err := conn.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	err = g.insertSpots(spots, tx)

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	return spots, nil
}

func (g *Gate) insertSpots(spots *map[string]*models.Spot, tx *sql.Tx) error {
	for key, spot := range *spots {
		err := spot.Insert(g.DbHandler, tx)

		if err != nil {
			return err
		}

		s := *spots
		s[key] = spot
	}

	return nil
}

func (g *Gate) GetPendingGames(userId int64) ([]models.Game, error) {
	games, err := models.GetPendingGames(g.DbHandler, userId)

	if err != nil {
		return nil, err
	}

	return games, nil
}

func (g *Gate) GetSingleGame(gameId int64) (*models.Game, error) {
	game, err := models.GetSingleGame(g.DbHandler, gameId)

	if err != nil {
		return nil, err
	}

	spotList, err := models.GetSpotsByGameId(g.DbHandler, gameId)

	if err != nil {
		return nil, err
	}

	spots := make(map[string]models.Spot)
	spotListVal := *spotList
	for i := range spotListVal {
		spot := spotListVal[i]
		id := strconv.Itoa(spot.X) + "," + strconv.Itoa(spot.Y)

		spots[id] = spot
	}
	game.Spots = spots

	return game, nil
}
