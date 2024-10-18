package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	db "github.com/suleimanodetoro/Go-Bank-Pro/db/sqlc"
)

// The `transferRequest` struct represents the structure of the incoming JSON payload
// for creating a new transfer. We use `binding` tags to enforce validation rules
// and ensure only valid data reaches our application.
type transferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`      // ID of the account to transfer from, must be positive
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`        // ID of the account to transfer to, must be positive
	Amount        int64  `json:"amount" binding:"required,gt=0"`                // Amount to transfer, must be greater than 0
	Currency      string `json:"currency" binding:"required,oneof=USD EUR CAD"` // Currency of the transfer, limited to USD, EUR, or CAD
}

// The `createTransfer` function handles the creation of a new transfer.
// It's a method on the `Server` struct, allowing it to access `Server` fields, such as the `store` for database operations.
// The `ctx` parameter is a Gin context that provides request-specific information and handles the response back to the client.
func (server *Server) createTransfer(ctx *gin.Context) {
	// `req` is an instance of `transferRequest`, which will store the parsed JSON request data
	var req transferRequest

	// `ShouldBindJSON` automatically parses the JSON from the incoming request into the `req` struct
	// If there's a validation error (such as invalid account IDs or amount), the function will return an error
	if err := ctx.ShouldBindJSON(&req); err != nil {
		// If the validation fails, return a `400 Bad Request` HTTP status code along with the error message
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if !server.validAccount(ctx, req.FromAccountID, req.Currency) {
		return
	}

	if !server.validAccount(ctx, req.ToAccountID, req.Currency) {
		return
	}

	// The `db.TransferTxParams` struct is an SQLC-generated struct that defines the parameters required
	// to create a new transfer in the database. It includes the source and destination account IDs and the transfer amount.
	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	// `server.store.TransferTx` is the actual function that interacts with the database to perform the transfer.
	// The `ctx` is passed to allow for context-based cancellations and timeouts, which can be useful in high-load scenarios.
	result, err := server.store.TransferTx(ctx, arg)
	if err != nil {
		// If the database transaction fails, return a `500 Internal Server Error` HTTP status code along with the error message.
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// If the transfer is successful, return the result with a `200 OK` status code.
	// The result object will be automatically marshaled into JSON format by Gin.
	ctx.JSON(http.StatusOK, result)
}

// validAccount checks if an account exists and has the correct currency.
// It returns true if the account is valid, false otherwise.
// This function also handles sending appropriate error responses via the Gin context.
func (server *Server) validAccount(ctx *gin.Context, accountID int64, currency string) bool {
	// Attempt to retrieve the account from the database
	account, err := server.store.GetAccount(ctx, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			// If the account doesn't exist, return a 404 Not Found error
			ctx.JSON(http.StatusNotFound, errorResponse(fmt.Errorf("account %d not found", accountID)))
			return false
		}
		// For any other database error, return a 500 Internal Server Error
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("error fetching account %d: %v", accountID, err)))
		return false
	}

	// Check if the account's currency matches the transfer currency
	// Using EqualFold for case-insensitive string comparison
	if !strings.EqualFold(account.Currency, currency) {
		// If currencies don't match, return a 400 Bad Request error
		err := fmt.Errorf("account [%d] currency mismatch: %s vs %s", accountID, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return false
	}

	// If all checks pass, return true indicating a valid account
	return true
}
