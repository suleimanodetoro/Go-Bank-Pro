package api

import (
	"github.com/gin-gonic/gin"                          // Gin framework for HTTP routing and middleware
	db "github.com/suleimanodetoro/Go-Bank-Pro/db/sqlc" // SQLC-generated package for database interaction
	// other imports
)

// The `Server` struct represents the core of the application. It holds the necessary dependencies
// (like the database store and the router) required to handle HTTP requests and interact with the database.
type Server struct {
	store  *db.Store   // The `store` is the interface to the database, where SQLC-generated methods are available to interact with the DB.
	router *gin.Engine // The `router` is used to define HTTP routes and handle incoming HTTP requests using the Gin framework.
}

// `NewServer` is a constructor function that creates a new instance of the `Server` struct.
// It takes in a `store` (the database connection handler) and sets up the routing for the HTTP server.
func NewServer(store *db.Store) *Server {
	// Initialize the `Server` struct with the given `store`.
	server := &Server{
		store: store, // This allows the server to access the database through SQLC-generated methods.
	}

	// Create a new instance of the Gin router using `gin.Default()`. This includes middleware for logging and recovery from panics.
	router := gin.Default()

	// Define an HTTP POST route for creating accounts. The handler for this route is the `createAccount` method of the `server`.
	// When a POST request is sent to `/accounts`, the `createAccount` method is called.
	router.POST("/accounts", server.createAccount)

	// Store the router in the `server` struct, so it can be used when starting the server.
	server.router = router
	return server
}

/*
## Concepts to Master:
1. **Structs in Go (`Server`):**
   - In Go, structs are user-defined types that group fields together.
   - The `Server` struct contains dependencies like the database store (`store`) and the HTTP router (`router`), making them accessible throughout the server's lifecycle.

2. **Dependency Injection (`NewServer`):**
   - `NewServer` is an example of dependency injection in Go. It accepts a `store` (the database interface) as an argument and injects it into the `Server` struct.
   - This pattern makes it easier to test and scale applications since dependencies can be swapped out (e.g., using a mock database for testing).

3. **Gin Router Setup:**
   - `gin.Default()` initializes a new instance of the Gin router with built-in middleware for logging and recovery from panics.
   - The router is responsible for handling HTTP requests and mapping them to specific handler functions, like `createAccount`.

4. **Routing in Gin:**
   - The `router.POST("/accounts", server.createAccount)` line defines a POST route.
   - This tells Gin to call `server.createAccount` when a client sends a POST request to `/accounts`.
   - Routing is a key part of any web framework, connecting HTTP methods (GET, POST, etc.) to specific handlers.
*/

// `Start` is a method on the `Server` struct that starts the HTTP server.
// It takes an `address` string, which defines the network address and port on which the server will listen for requests (e.g., "localhost:8080").
func (server *Server) Start(address string) error {
	// `server.router.Run(address)` starts the HTTP server using the Gin router.
	// This function blocks and listens for incoming HTTP requests until the server is manually stopped.
	// If any errors occur during startup, they will be returned from `Run()`, which is why `error` is returned here.
	return server.router.Run(address)
}

/*
## Concepts to Master:
1. **HTTP Server Startup:**
   - The `Start` method is responsible for starting the HTTP server. It binds the server to a specified address (e.g., localhost:8080) and listens for incoming HTTP requests.
   - Gin’s `Run` method starts the server and handles incoming requests based on the routes defined earlier.

2. **Addressing:**
   - The `address` parameter determines where the server will be accessible. It typically consists of an IP address and port number (e.g., `127.0.0.1:8080`).
   - This is crucial for local development, testing, or production environments where servers listen on specific addresses.

3. **Error Handling:**
   - Since running a server can fail (e.g., due to the port being in use), the method returns an error. Proper error handling is vital to gracefully manage unexpected conditions.
*/

// `errorResponse` is a helper function that takes an `error` object and returns a standardized JSON response in the form of a map.
// This is used to return consistent error messages in the API responses.
func errorResponse(err error) gin.H {
	// Gin's `H` type is just a shortcut for `map[string]interface{}`.
	// It allows you to return JSON responses with key-value pairs, where the key is a string (in this case, "error") and the value is dynamic (in this case, the error message).
	return gin.H{"error": err.Error()}
}

/*
## Concepts to Master:
1. **Helper Functions:**
   - `errorResponse` is a reusable function for returning standardized error responses to the client.
   - In this case, it wraps the error message in a JSON format with a key of `"error"`.

2. **Returning JSON Responses:**
   - The Gin framework simplifies returning JSON by using `gin.H` (which is essentially a map) to build a key-value response.
   - The `ctx.JSON` method (seen in other parts of the code) will use this helper function to return consistent error messages.

## Algorithm Explanation:

1. **Creating the Server:**
   - When `NewServer` is called, a new server instance is created. This instance contains:
     - A reference to the database store (`store`), allowing it to execute database queries.
     - A Gin router (`router`) that handles incoming HTTP requests and routes them to the appropriate handlers (e.g., `createAccount`).

2. **Defining Routes:**
   - The route `/accounts` is defined using `router.POST`. When the client sends a POST request to `/accounts`, Gin calls the `createAccount` handler on the `server`.

3. **Starting the Server:**
   - The `Start` method binds the server to a specific address (e.g., `localhost:8080`) and listens for incoming requests.
   - The Gin framework handles incoming HTTP requests and matches them to the defined routes. For example, if a client sends a POST request to `/accounts`, the `createAccount` handler is called.

4. **Error Handling:**
   - If an error occurs (e.g., a request validation failure or a database error), the `errorResponse` function is used to return a standardized JSON error message to the client.

By structuring the server in this way, you create a scalable architecture where new routes can be easily added, and all the core server logic (starting, routing, and error handling) is centralized within the `Server` struct.
*/
