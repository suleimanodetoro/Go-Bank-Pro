package api

import (
	"database/sql"
	"net/http" // Package for HTTP utilities like status codes

	"github.com/gin-gonic/gin"                          // Gin framework for building web applications
	db "github.com/suleimanodetoro/Go-Bank-Pro/db/sqlc" // Importing the db package to access SQLC-generated code for database queries
)

// The `createAccountRequest` struct is used to represent the structure of the incoming JSON payload
// that clients will send when creating a new account. We use `binding` tags to enforce validation rules
// and ensure only valid data reaches our application. In this case, `Owner` and `Currency` are required fields.
type createAccountRequest struct {
	Owner    string `json:"owner" binding:"required"`                  // `Owner` represents the account owner's name, required field
	Currency string `json:"currency" binding:"required,oneof=USD EUR"` // `Currency` field restricted to USD or EUR using the `oneof` validator
}

// The `createAccount` function handles the creation of a new account.
// It's a method on the `Server` struct, allowing it to access `Server` fields, such as the `store` for database operations.
// The `ctx` parameter is a Gin context that provides request-specific information and handles the response back to the client.
func (server *Server) createAccount(ctx *gin.Context) {
	// `req` is an instance of `createAccountRequest`, which will store the parsed JSON request data
	var req createAccountRequest

	// `ShouldBindJSON` automatically parses the JSON from the incoming request into the `req` struct
	// If there's a validation error (such as missing `Owner` or invalid `Currency`), the function will return an error
	if err := ctx.ShouldBindJSON(&req); err != nil {
		// If the validation fails, return a `400 Bad Request` HTTP status code along with the error message
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// The `db.CreateAccountParams` struct is an SQLC-generated struct that defines the parameters required
	// to create a new account in the database. In this case, it requires the `Owner`, `Currency`, and an initial `Balance`.
	arg := db.CreateAccountParams{
		Owner:    req.Owner,    // Setting the owner of the account from the request data
		Currency: req.Currency, // Setting the currency of the account from the request data
		Balance:  0,            // The initial balance is always set to 0
	}

	// `server.store.CreateAccount` is the actual function that interacts with the database to create the account.
	// The `ctx` is passed to allow for context-based cancellations and timeouts, which can be useful in high-load scenarios.
	account, err := server.store.CreateAccount(ctx, arg)
	if err != nil {
		// If the database query fails, return a `500 Internal Server Error` HTTP status code along with the error message.
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// If everything is successful, return the created account with a `200 OK` status code.
	// The account object will be automatically marshaled into JSON format by Gin.
	ctx.JSON(http.StatusOK, account)
}

type getAccountRequest struct {
	// Bind the `id` from the URI (URL) to this field.
	// The `uri:"id"` tag tells Gin to extract the `id` from the URL and bind it to this field.
	// The `binding:"required,min=1"` ensures that the ID must be present and greater than or equal to 1.
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getAccount(ctx *gin.Context) {
	var req getAccountRequest

	// Bind the URI parameters (in this case, the ID) to the `getAccountRequest` struct.
	// `ShouldBindUri` extracts the `id` parameter from the URL (e.g., /accounts/:id) and binds it to `req.ID`.
	// If binding fails (e.g., if `id` is missing or invalid), it returns a `400 Bad Request`.
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Fetch the account from the database using the bound ID.
	account, err := server.store.GetAccount(ctx, req.ID)
	if err != nil {
		// If the error indicates that no account was found with the given ID, return `404 Not Found`.
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		// For all other database-related errors, return `500 Internal Server Error`.
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// If the account is found, return it with a `200 OK` response.
	ctx.JSON(http.StatusOK, account)
}

type listAccountRequest struct {
	// Bind the `page_id` from the query string (e.g., /accounts?page_id=1) to this field.
	// The `form:"page_id"` tag tells Gin to extract `page_id` from the query string.
	// The `binding:"required,min=1"` ensures that the page ID must be provided and be at least 1.
	PageID int32 `form:"page_id" binding:"required,min=1"`

	// Bind the `page_size` from the query string (e.g., /accounts?page_size=10) to this field.
	// The `binding:"required,min=5,max=10"` ensures that the page size must be between 5 and 10,
	// controlling how many results are returned per page.
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listAccount(ctx *gin.Context) {
	var req listAccountRequest

	// Bind the query parameters from the URL query string (e.g., /accounts?page_id=1&page_size=10)
	// to the `listAccountRequest` struct. This will automatically extract and validate the
	// `page_id` and `page_size` parameters. If the query parameters are invalid or missing,
	// it returns a `400 Bad Request` response.
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Create the parameters for the database query, including pagination.
	// `Limit` is the number of results to return per page.
	// `Offset` calculates where the page starts by multiplying the (PageID - 1) by the PageSize.
	arg := db.ListAccountsParams{
		Limit:  req.PageSize,                    // Limit the number of results returned to the specified page size.
		Offset: (req.PageID - 1) * req.PageSize, // Calculate the offset for pagination.
	}

	// Fetch the paginated list of accounts from the database.
	accounts, err := server.store.ListAccounts(ctx, arg)
	if err != nil {
		// If there is a database error, return a `500 Internal Server Error` response.
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// If the accounts are successfully retrieved, return them with a `200 OK` status.
	// The accounts will be automatically serialized into JSON and sent back in the response.
	ctx.JSON(http.StatusOK, accounts)
}

/*
## Concepts to Master:

1. **Go Structs and Tags (`createAccountRequest`):**
   - Go structs are user-defined types that group fields together.
   - Struct tags (e.g., `json:"owner"`, `binding:"required"`) are metadata that tell how fields should be processed.
   - In this example, `json:"owner"` means that the `Owner` field will be serialized to/from JSON as "owner".
   - The `binding` tag tells Gin to validate incoming requests and ensure that the required fields are present.

2. **Gin Context (`ctx`):**
   - The `gin.Context` object is central to how Gin handles requests and responses.
   - It holds all the information about the HTTP request (like headers, body, query parameters) and provides methods
     for sending responses back to the client.
   - For example, `ctx.ShouldBindJSON` binds the incoming JSON to a Go struct, and `ctx.JSON` sends a JSON response.

3. **Error Handling in Go:**
   - Go uses explicit error handling. If an error occurs, it must be handled manually rather than through exceptions.
   - In the case of binding errors or database errors, the `if err != nil` pattern checks if the error occurred and handles it appropriately.
   - `errorResponse` is a helper function (defined elsewhere) that converts errors into a consistent JSON structure for the client.

4. **SQLC and Database Operations:**
   - SQLC is a tool that generates type-safe Go code from SQL queries.
   - In this code, `db.CreateAccountParams` is an SQLC-generated struct that corresponds to the parameters required by the `CREATE ACCOUNT` SQL statement.
   - The `server.store.CreateAccount` method is an SQLC-generated function that runs the SQL query to insert a new account into the database.

5. **HTTP Status Codes:**
   - The `ctx.JSON` method sends an HTTP response back to the client. The first argument is the HTTP status code:
     - `http.StatusBadRequest` (400) indicates the client sent invalid data.
     - `http.StatusInternalServerError` (500) indicates a server-side issue.
     - `http.StatusOK` (200) indicates success.
   - Returning the correct status code is crucial for API clients to understand how to handle the response.

6. **Context (`ctx`) for Cancellations and Timeouts:**
   - The `ctx` parameter (passed to `server.store.CreateAccount`) allows the request to be canceled or timed out.
   - This is particularly useful in production systems where long-running requests may need to be terminated to free up resources.

## Algorithm Explanation (Step-by-Step):

1. **Receive HTTP Request:**
   - A POST request is sent to create an account with a JSON payload that contains the `Owner` and `Currency` fields.

2. **Bind and Validate Request:**
   - Gin automatically parses the incoming JSON into the `createAccountRequest` struct.
   - The `binding:"required"` tag ensures that the `Owner` and `Currency` fields are present. If they are missing or invalid, an error is returned.

3. **Prepare Database Query:**
   - Using the parsed data (`Owner`, `Currency`), a `CreateAccountParams` struct is populated. This struct is used by SQLC to construct the database query for creating the account.

4. **Execute Database Query:**
   - The `CreateAccount` method on `server.store` (generated by SQLC) is called to execute the SQL query and insert the new account into the database.
   - If this query fails (e.g., due to a database error), the server returns a `500 Internal Server Error` to the client.

5. **Return Success or Failure:**
   - If the account creation succeeds, a `200 OK` response is returned, and the created account is sent back to the client as JSON.
   - If any step fails, an appropriate error response (`400` for bad request, `500` for internal server errors) is returned to the client.
*/
