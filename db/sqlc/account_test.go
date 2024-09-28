package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/suleimanodetoro/Go-Bank-Pro/db/util"
)

// CreateRandomAccount creates a random account for testing purposes and returns it.
func CreateRandomAccount(t *testing.T) Account {
	// Prepare the arguments to create an account
	arg := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	// Create the account and check for errors
	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)                          // Assert that there is no error
	require.NotEmpty(t, account)                     // Assert that the account is not empty
	require.Equal(t, arg.Owner, account.Owner)       // Check if the account owner is the same as the argument
	require.Equal(t, arg.Balance, account.Balance)   // Check if the account balance is the same as the argument
	require.Equal(t, arg.Currency, account.Currency) // Check if the currency is the same as the argument

	require.NotZero(t, account.ID)        // Assert that the account ID is not zero
	require.NotZero(t, account.CreatedAt) // Assert that the account creation time is not zero

	// Return the created account for use in other tests
	return account
}

// TestCreateAccount tests the account creation functionality
func TestCreateAccount(t *testing.T) {
	// Call CreateRandomAccount to test account creation
	CreateRandomAccount(t)
}

// TestGetAccount tests retrieving an account from the database
func TestGetAccount(t *testing.T) {
	// Create a new random account
	account1 := CreateRandomAccount(t)

	// Retrieve the account using its ID
	account2, err := testQueries.GetAccount(context.Background(), account1.ID)

	require.NoError(t, err)       // Assert no error
	require.NotEmpty(t, account1) // Assert that the create account is not empty
	require.NotEmpty(t, account2) // Assert that the retrieved account is not empty

	// Check if the retrieved account matches the created account
	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
}
