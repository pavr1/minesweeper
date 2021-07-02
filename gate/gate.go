package gate

import (
	"context"
	"database/sql"
	"fmt"
	"minesweeper/ccache"
	"minesweeper/dbhandler"
	"minesweeper/models"
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

	spots, err := game.GenerateGrid()

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = g.InsertSpots(spots, tx, &ctx)

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	return *spots, nil
}

func (g *Gate) InsertSpots(spots *map[string]models.Spot, tx *sql.Tx, ctx *context.Context) error {
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
