package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/suleimanodetoro/Go-Bank-Pro/db/util"
)

// CreateRandomUser generates and inserts a random user into the database for test purposes.
// This helper function focuses on setting up data, not performing assertions.
func CreateRandomUser(t *testing.T) User {
	hashedPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)
	// Prepare the parameters to create a user
	arg := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}

	// Insert the user into the database
	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)   // Ensure no error occurred during user creation
	require.NotEmpty(t, user) // Assert that the created user is not empty

	// Perform assertions to verify that the created user matches the input
	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)

	require.True(t, user.PasswordChangedAt.IsZero()) // Ensure password_changed_at is zero (initial state)
	require.NotZero(t, user.CreatedAt)               // Ensure created_at is set

	return user // Return the created user for use in tests
}

// TestCreateUser tests the user creation functionality.
func TestCreateUser(t *testing.T) {
	// Use the helper to create a random user
	CreateRandomUser(t)
}

// TestGetUser tests retrieving a user from the database by username.
func TestGetUser(t *testing.T) {
	// Create a new random user to retrieve
	user1 := CreateRandomUser(t)

	// Retrieve the user from the database using its username
	user2, err := testQueries.GetUser(context.Background(), user1.Username)
	require.NoError(t, err)    // Ensure no error occurred during user retrieval
	require.NotEmpty(t, user2) // Assert that the retrieved user is not empty

	// Perform assertions to ensure that the retrieved user matches the created one
	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.Email, user2.Email)
	require.WithinDuration(t, user1.PasswordChangedAt, user2.PasswordChangedAt, time.Second)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
}
