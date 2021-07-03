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

func (g *Gate) CreateGame(game models.Game) (map[string]models.Spot, error) {
	if game.UserId == 0 {
		return nil, fmt.Errorf("user id required")
	}

	ctx := context.Background()
	tx, err := g.DbHandler.Db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	err = game.Create(g.DbHandler, tx, &ctx)

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	gamePtr, err := models.GetLatestGame(game.UserId, g.DbHandler, tx, &ctx)

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	game = *gamePtr
	spots, err := game.GenerateGrid()

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = g.insertSpots(spots, tx, &ctx)

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	return *spots, nil
}

func (g *Gate) insertSpots(spots *map[string]models.Spot, tx *sql.Tx, ctx *context.Context) error {
	for key, spot := range *spots {
		err := spot.Insert(g.DbHandler, tx, ctx)

		if err != nil {
			return err
		}

		s := *spots
		s[key] = spot
	}

	return nil
}

func (g *Gate) GetPendingGames(userId int64) ([]models.Game, error) {
	games, err := models.GetPendingGames(userId, g.DbHandler)

	if err != nil {
		return nil, err
	}

	return games, nil
}

func (g *Gate) GetSingleGame(gameId int64) (*models.Game, error) {
	game, err := models.GetSingleGame(gameId, g.DbHandler)

	if err != nil {
		return nil, err
	}

	game.Spots = make(map[string]models.Spot)
	spotList, err := models.GetSpotsByGameId(gameId, g.DbHandler)

	if err != nil {
		return nil, err
	}

	for i := range spotList {
		spot := spotList[i]
		id := strconv.Itoa(spot.X) + "," + strconv.Itoa(spot.Y)

		game.Spots[id] = spot
	}

	if err != nil {
		return nil, err
	}

	return game, nil
}
