package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/suleimanodetoro/Go-Bank-Pro/db/util"
)

var testQueries *Queries
var testDB *sql.DB

// Test main is the entry point for all unit tests in your package by default
func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("wahala ti burst o:", err)
	}
	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to database to perform tests", err)
	}
	// If there are no errors, initialize testQueries with the connection object
	testQueries = New(testDB)
	os.Exit(m.Run())

}
