package api

import (
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

// config gin to use test mode to shorten test result statements
// very similar to previous main_test file in db package

// Test main is the entry point for all unit tests in your package by default
func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	os.Exit(m.Run())

}
