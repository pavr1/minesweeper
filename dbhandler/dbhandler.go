package dbhandler

import (
	"context"
	"database/sql"
	"fmt"

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

	if err != nil {
		return nil, fmt.Errorf("Error creating db instance: %s" + err.Error())
	}

	return db, nil
}

func (h *DbHandler) Execute(statement string, args []string) error {
	ctx := context.Background()
	var err error
	var db *sql.DB

	if h.Db == nil {
		db, err = createDatabase()

		if err != nil {
			return err
		}

		h.Db = db
	}

	err = h.Db.PingContext(ctx)
	if err != nil {
		return fmt.Errorf("Error pinging db server: %s" + err.Error())
	}

	params := make([]interface{}, len(args))
	for i := range args {
		params[i] = args[i]
	}

	_, err = h.Db.Exec(statement, params...)

	if err != nil {
		return fmt.Errorf("Error executing statement: %s" + err.Error())
	}

	return nil
}
