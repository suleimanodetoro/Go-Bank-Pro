package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

// set up dtaabse conection to enable test for database (defined in _test.go files)
const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
)

var testQueries *Queries
var testDB *sql.DB

// Test main is the entry point for all unit tests in your package by default
func TestMain(m *testing.M) {
	var err error
	testDB, err = sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to database to perform tests", err)
	}
	// If there are no errors, initialize testQueries with the connection object
	testQueries = New(testDB)
	os.Exit(m.Run())

}
