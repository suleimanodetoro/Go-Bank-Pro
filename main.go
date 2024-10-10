package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/suleimanodetoro/Go-Bank-Pro/api"
	db "github.com/suleimanodetoro/Go-Bank-Pro/db/sqlc"
)

// in order to create a server, we need to connect to the database and create a store
// set up dataabse conection to enable test for database (defined in _test.go files)
const (
	dbDriver      = "postgres"
	dbSource      = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
	serverAddress = "0.0.0.0:8080"
)

func main() {
	// establish connection
	var err error
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to database", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(serverAddress)
	if err != nil {
		log.Fatal("Cannot't start server!", err)

	}

}
