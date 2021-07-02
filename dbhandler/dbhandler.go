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
	UserId      int
	Name        string
	LastName    string
	Password    string
	CreatedDate time.Time
	Message     string
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

	result, err := tx.QueryContext(*ctx, statement, args...)

	if err != nil {
		return -1, fmt.Errorf("Error executing statement: %s" + err.Error())
	}

	var id int64
	result.Next()
	err = result.Scan(&id)

	if err != nil {
		fmt.Println(err.Error())
		return -1, fmt.Errorf("Error retrieving latest id after insert: %s" + err.Error())
	}

	return -1, err
}

func (h *DbHandler) Select(statement, structType string, args []string) ([]interface{}, error) {
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
			var userId int
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
	case "Spot":
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
