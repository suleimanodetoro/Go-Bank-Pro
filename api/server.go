package api

import (
	"github.com/gin-gonic/gin" // Gin framework for HTTP routing and middleware
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/suleimanodetoro/Go-Bank-Pro/db/sqlc" // SQLC-generated package for database interaction
)

// The `Server` struct represents the core of the application, holding dependencies like the database store and the router.
// It serves HTTP requests, delegating the actual work to the underlying database via the store interface.
type Server struct {
	store  db.Store    // Store is the interface to the database where SQLC-generated methods are available for interaction.
	router *gin.Engine // Router is used to define HTTP routes and handle incoming HTTP requests using the Gin framework.
}

// `NewServer` is a constructor function that creates a new instance of the `Server` struct.
// It takes in a `store` (the database handler) and sets up the HTTP routing for the server.
func NewServer(store db.Store) *Server {
	server := &Server{
		store: store, // Inject the database store into the server.
	}

	router := gin.Default() // Initialize a new Gin router with logging and recovery middleware.

	// Force the validator to initialize
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	// Define the routes for the server, mapping HTTP methods to handler functions.
	router.POST("/accounts", server.createAccount)   // Route for creating an account.
	router.GET("/accounts/:id", server.getAccount)   // Route for fetching a single account by ID.
	router.GET("/accounts", server.listAccount)      // Route for listing accounts with optional pagination.
	router.POST("/transfers", server.createTransfer) // Route for creating a transfer

	server.router = router // Assign the router to the server instance.
	return server
}

// `Start` is a method on the `Server` struct that starts the HTTP server, binding it to a specific address and port.
func (server *Server) Start(address string) error {
	// Start the HTTP server using the Gin router and listen for incoming requests.
	return server.router.Run(address)
}

// `errorResponse` is a helper function that formats error messages into a standardized JSON response.
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()} // Return a JSON object with an "error" field containing the error message.
}
