package gate

import (
	"context"
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
	gate := &Gate{
		DbHandler: &dbhandler.DbHandler{},
		Cache:     ccache.New(),
	}

	gate.DbHandler.CheckConnection()

	return gate, nil
}

func (g *Gate) CreateUser(name, lastName, password string) (models.UIUser, error) {
	if name == "" {
		return models.UIUser{}, fmt.Errorf("user name required")
	}

	if lastName == "" {
		return models.UIUser{}, fmt.Errorf("user last name required")
	}

	if password == "" {
		return models.UIUser{}, fmt.Errorf("password required")
	}

	user := models.User{
		Name:     name,
		LastName: lastName,
		Password: password,
	}

	ui, err := user.CreateUser(g.DbHandler)

	if err != nil {
		return models.UIUser{}, err
	}

	return *ui, nil
}

func (g *Gate) ValidateLogin(name, password string) (models.UIUser, error) {
	if name == "" {
		return models.UIUser{}, fmt.Errorf("user name required")
	}

	if password == "" {
		return models.UIUser{}, fmt.Errorf("password required")
	}

	user := models.User{
		Name:     name,
		Password: password,
	}

	result, err := user.ValidateUser(g.DbHandler)

	if err != nil {
		return models.UIUser{}, err
	}

	return *result, nil
}

func (g *Gate) CreateGame(ctx context.Context, userId int64, rows, columns, mines int) (models.UIGame, error) {
	if userId == 0 {
		return models.UIGame{}, fmt.Errorf("user id required")
	}

	game := &models.Game{
		UserId:       userId,
		TimeConsumed: 0,
		Status:       "Pending",
		Rows:         rows,
		Columns:      columns,
		Mines:        mines,
	}

	ui, err := game.Create(ctx, g.DbHandler)

	if err != nil {
		return models.UIGame{}, err
	}

	return *ui, nil
}

func (g *Gate) GetPendingGames(userId int64) ([]models.UIGame, error) {
	games, err := models.GetPendingGames(g.DbHandler, userId)

	if err != nil {
		return nil, err
	}

	return games, nil
}

func (g *Gate) GetSingleGame(gameId int64) (models.UIGame, error) {
	ui, err := models.GetSingleGame(g.DbHandler, gameId)

	if err != nil {
		return models.UIGame{}, err
	}

	return *ui, nil
}

func (g *Gate) ProcessSpot(gameId, spotId int64, status string) (models.UIGame, error) {
	ui, err := models.GetSingleGame(g.DbHandler, int64(gameId))
	if err != nil {
		return models.UIGame{}, err
	}

	spot, err := models.GetSpotById(g.DbHandler, int64(spotId))
	if err != nil {
		return models.UIGame{}, err
	}

	err = spot.ProcessSpot(g.DbHandler, ui.Rows, ui.Columns, status)

	game := models.Game{
		GameId: ui.GameId,
	}

	gameOver := false
	if err != nil {
		ui.Message = err.Error()

		if ui.Message == "Game Over!" {
			gameOver = true
			game.Status = "Lost"
			game.UpdateStatus(g.DbHandler)
		}
	}

	ui, err = models.GetSingleGame(g.DbHandler, ui.GameId)
	if err != nil {
		return models.UIGame{}, err
	}

	allSpotsOpen := true

	for _, spot := range ui.Spots {
		if spot.Status == "Closed" {
			allSpotsOpen = false
			break
		}
	}

	if allSpotsOpen {
		if !gameOver {
			ui.Message = "Congratulations, game finished!"
			game.Status = "Won"
			game.UpdateStatus(g.DbHandler)
		} else {
			game.Message = "Game Over!"
		}
	}

	return *ui, nil
}
