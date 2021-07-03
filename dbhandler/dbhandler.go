package dbhandler

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
)

// Replace with your own connection parameters
var server = "127.0.0.1"
var port = 1434
var user = "minesweeper"
var password = "minesweeper"
var databasename = "minesweeper"

type DbHandler struct {
	Db *sql.DB
}

type DbUser struct {
	UserId      int64
	Name        string
	LastName    string
	Password    string
	CreatedDate time.Time
	Message     string
}

type DbGame struct {
	UserId       int64
	GameId       int64
	CreatedDate  time.Time
	TimeConsumed float32
	Status       string
	Rows         int
	Columns      int
	Mines        int
}

type DbSpot struct {
	SpotId    int64
	GameId    int64
	Value     string
	X         int
	Y         int
	NearSpots string //map[string]Spot
	Status    string
}

func GetInstance() (*DbHandler, error) {
	db, err := createDatabase()

	if err != nil {
		return nil, err
	}

	handler := DbHandler{
		Db: db,
	}

	return &handler, nil
}

func createDatabase() (*sql.DB, error) {
	connString := fmt.Sprintf("server=%s;database=%s;port=%d;Trusted_Connection=true", server, databasename, port)

	db, err := sql.Open("sqlserver", connString)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)

	if err != nil {
		return nil, fmt.Errorf("Error creating db instance: %s" + err.Error())
	}

	return db, nil
}

func (h *DbHandler) Execute(statement string, args []interface{}) (int64, error) {
	ctx := context.Background()
	var err error
	var db *sql.DB

	if h.Db == nil {
		db, err = createDatabase()

		if err != nil {
			return -1, err
		}

		h.Db = db
	}

	err = h.Db.PingContext(ctx)
	if err != nil {
		fmt.Println(err.Error())
		return -1, fmt.Errorf("Error pinging db server: %s" + err.Error())
	}

	result, err := h.Db.QueryContext(ctx, statement, args...)

	if err != nil {
		return -1, fmt.Errorf("Error executing statement: %s" + err.Error())
	}

	var id int64
	result.Next()
	err = result.Scan(&id)

	if err != nil {
		return -1, fmt.Errorf("Error retrieving latest id after insert: %s" + err.Error())
	}

	return id, err
}

func (h *DbHandler) ExecuteTransaction(statement string, args []interface{}, tx *sql.Tx, ctx *context.Context) (int64, error) {
	var err error

	result, err := tx.Query(statement, args...)

	if err != nil {
		return -1, fmt.Errorf("Error executing statement: %s" + err.Error())
	}

	var id int64
	result.Next()
	err = result.Scan(&id)

	if err != nil {
		fmt.Println(fmt.Errorf("Error retrieving latest id after insert: " + err.Error()))
		//return -1, fmt.Errorf("Error retrieving latest id after insert: " + err.Error())
	}

	return id, nil
}

func (h *DbHandler) Select(statement, structType string, args []interface{}) ([]interface{}, error) {
	ctx := context.Background()

	params := make([]interface{}, len(args))
	for i := range args {
		params[i] = args[i]
	}

	rows, err := h.Db.QueryContext(ctx, statement, params...)
	if err != nil {
		return nil, err
	}

	var result []interface{}
	defer rows.Close()

	switch structType {
	case "User":
		var users = []DbUser{}

		for rows.Next() {
			var userId int64
			var name string
			var lastName string
			var password string
			var createdDate time.Time

			rows.Scan(&userId, &name, &lastName, &password, &createdDate)

			users = append(users, DbUser{
				UserId:      userId,
				Name:        name,
				LastName:    lastName,
				Password:    password,
				CreatedDate: createdDate,
			})
		}

		result = make([]interface{}, len(users))
		for i, v := range users {
			result[i] = v
		}

		return result, nil
	case "Game":
		var games = []DbGame{}

		for rows.Next() {
			var gameId int64
			var userId int64
			var createdDate time.Time
			var timeConsumed int
			var status string
			var rows2 int
			var columns int
			var mines int

			rows.Scan(&gameId, &userId, &createdDate, &timeConsumed, &status, &rows2, &columns, &mines)

			games = append(games, DbGame{
				GameId:       gameId,
				UserId:       userId,
				CreatedDate:  createdDate,
				TimeConsumed: float32(timeConsumed),
				Status:       status,
				Rows:         rows2,
				Columns:      columns,
				Mines:        mines,
			})
		}

		result = make([]interface{}, len(games))
		for i, v := range games {
			result[i] = v
		}

		return result, nil
	case "Spot":
		var spots = []DbSpot{}

		for rows.Next() {
			var spotId int64
			var gameId int64
			var value string
			var x int
			var y int
			var nearSpots string //map[string]Spot
			var status string

			rows.Scan(&spotId, &gameId, &value, &x, &y, &nearSpots, &status)

			spots = append(spots, DbSpot{
				SpotId:    spotId,
				GameId:    gameId,
				Value:     value,
				X:         x,
				Y:         y,
				NearSpots: nearSpots,
				Status:    status,
			})
		}

		result = make([]interface{}, len(spots))
		for i, v := range spots {
			result[i] = v
		}

		return result, nil
	}

	// If the database is being written to ensure to check for Close
	// errors that may be returned from the driver. The query may
	// encounter an auto-commit error and be forced to rollback changes.
	err = rows.Close()
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (h *DbHandler) SelectTransaction(statement, structType string, args []interface{}, tx *sql.Tx, ctx *context.Context) ([]interface{}, error) {
	params := make([]interface{}, len(args))
	for i := range args {
		params[i] = args[i]
	}

	rows, err := tx.QueryContext(*ctx, statement, params...)
	if err != nil {
		return nil, err
	}

	var result []interface{}
	defer rows.Close()

	switch structType {
	case "User":
		var users = []DbUser{}

		for rows.Next() {
			var userId int64
			var name string
			var lastName string
			var password string
			var createdDate time.Time

			rows.Scan(&userId, &name, &lastName, &password, &createdDate)

			users = append(users, DbUser{
				UserId:      userId,
				Name:        name,
				LastName:    lastName,
				Password:    password,
				CreatedDate: createdDate,
			})
		}

		result = make([]interface{}, len(users))
		for i, v := range users {
			result[i] = v
		}

		return result, nil
	case "Game":
		var games = []DbGame{}

		for rows.Next() {
			var gameId int64
			var userId int64
			var createdDate time.Time
			var timeConsumed int
			var status string
			var rows2 int
			var columns int
			var mines int

			rows.Scan(&gameId, &userId, &createdDate, &timeConsumed, &status, &rows2, &columns, &mines)

			games = append(games, DbGame{
				GameId:       gameId,
				UserId:       userId,
				CreatedDate:  createdDate,
				TimeConsumed: float32(timeConsumed),
				Status:       status,
				Rows:         rows2,
				Columns:      columns,
				Mines:        mines,
			})
		}

		result = make([]interface{}, len(games))
		for i, v := range games {
			result[i] = v
		}

		return result, nil
	case "Spot":
		var spots = []DbSpot{}

		for rows.Next() {
			var spotId int64
			var gameId int64
			var value string
			var x int
			var y int
			var nearSpots string //map[string]Spot
			var status string

			rows.Scan(&spotId, &gameId, &value, &x, &y, &nearSpots, &status)

			spots = append(spots, DbSpot{
				SpotId:    spotId,
				GameId:    gameId,
				Value:     value,
				X:         x,
				Y:         y,
				NearSpots: nearSpots,
				Status:    status,
			})
		}

		result = make([]interface{}, len(spots))
		for i, v := range spots {
			result[i] = v
		}

		return result, nil
	}

	// If the database is being written to ensure to check for Close
	// errors that may be returned from the driver. The query may
	// encounter an auto-commit error and be forced to rollback changes.
	err = rows.Close()
	if err != nil {
		return nil, err
	}

	return result, nil
}
