package api

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	db "github.com/suleimanodetoro/Go-Bank-Pro/db/sqlc"
	"github.com/suleimanodetoro/Go-Bank-Pro/db/util"
)

// createUserRequest represents the structure of the incoming JSON payload for creating a new user.
type createUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"` // Username must be alphanumeric and is required
	Password string `json:"password" binding:"required,min=6"`    // Password must have a minimum length of 6 characters
	FullName string `json:"full_name" binding:"required"`         // Full name is required
	Email    string `json:"email" binding:"required,email"`       // Email must be a valid email address and is required
}

// customer struct to omit hashed password from api response
type createUserResponse struct {
	Username string `json:"username"`
	// HashedPassword    string    `json:"hashed_password"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

// createUser handles the creation of a new user.
func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest

	// Bind the incoming JSON to the request struct
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Hash the user's password
	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Prepare the parameters to create the user in the database
	arg := db.CreateUserParams{
		Username:       req.Username,
		HashedPassword: hashedPassword,
		FullName:       req.FullName,
		Email:          req.Email,
	}

	// Create the user in the database
	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		// Handle PostgreSQL-specific errors
		if pqErr, ok := err.(*pq.Error); ok {
			log.Printf("PostgreSQL Error Code: %s, Message: %s", pqErr.Code, pqErr.Message)

			// Handle unique constraint violations
			if pqErr.Code.Name() == "unique_violation" {
				ctx.JSON(http.StatusConflict, errorResponse(err)) // Use 409 Conflict for unique constraint violations
				return
			}
		}

		// For all other errors, return a 500 Internal Server Error
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// create the customer api response without hashedPassword
	response := createUserResponse{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}

	// Return the created user with a 200 OK status
	ctx.JSON(http.StatusOK, response)
}
