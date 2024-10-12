package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/suleimanodetoro/Go-Bank-Pro/api"
	db "github.com/suleimanodetoro/Go-Bank-Pro/db/sqlc"
	"github.com/suleimanodetoro/Go-Bank-Pro/db/util"
)

// in order to create a server, we need to connect to the database and create a store
// set up dataabse conection to enable test for database (defined in _test.go files)

func main() {
	// Load variables from environment
	config, err := util.LoadConfig(".") //. is current folder, cause config file is in same directory as this file
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	// establish connection
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to database", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("Cannot't start server!", err)

	}

}
