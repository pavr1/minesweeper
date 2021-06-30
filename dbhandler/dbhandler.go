package dbhandler

import (
	"context"
	"database/sql"
	"fmt"
)

// Replace with your own connection parameters
var server = "localhost"
var port = 1433
var user = "sa"
var password = "your_password"

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
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d", server, user, password, port)

	db, err := sql.Open("sqlserver", connString)

	if err != nil {
		return nil, fmt.Errorf("Error creating db instance: %s" + err.Error())
	}

	return db, nil
}

func (h *DbHandler) Execute(statement string, args []sql.NamedArg) error {
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

	//tsql := "INSERT INTO TestSchema.Employees (Name, Location) VALUES (@Name, @Location); select convert(bigint, SCOPE_IDENTITY());"

	stmt, err := db.Prepare(statement)
	if err != nil {
		return fmt.Errorf("Error preparing statement: %s" + err.Error())
	}
	defer stmt.Close()

	//args[0] = sql.Named("test", "test")

	/*row := */
	stmt.QueryRowContext(
		ctx, args)
	// sql.Named("Name", name),
	// sql.Named("Location", location))
	// var newID int64
	// err = row.Scan(&newID)
	// if err != nil {
	// 	return err
	// }

	return nil
}
