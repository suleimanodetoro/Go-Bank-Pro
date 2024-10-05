package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/suleimanodetoro/Go-Bank-Pro/db/util"
)

// CreateRandomAccount generates and inserts a random account into the database for test purposes.
// This helper function focuses on setting up data, not performing assertions.
func CreateRandomAccount(t *testing.T) Account {
	arg := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	// Create the account in the database
	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err) // Ensure no error occurred during account creation
	return account          // Return the created account for use in tests
}

// TestCreateAccount tests the account creation functionality.
func TestCreateAccount(t *testing.T) {
	// Use the helper to create a random account
	account := CreateRandomAccount(t)

	// Perform assertions on the created account
	require.NotEmpty(t, account)          // Assert that the account is not empty
	require.NotZero(t, account.ID)        // Ensure account ID is not zero
	require.NotZero(t, account.CreatedAt) // Ensure account creation time is set
	require.NotEmpty(t, account.Owner)    // Owner should not be empty
	require.NotEmpty(t, account.Currency) // Currency should not be empty
}

// TestGetAccount tests retrieving an account from the database by ID.
func TestGetAccount(t *testing.T) {
	// Create a new random account to retrieve
	account1 := CreateRandomAccount(t)

	// Retrieve the account from the database using its ID
	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)       // Ensure no error occurred during account retrieval
	require.NotEmpty(t, account2) // Assert that the retrieved account is not empty

	// Perform assertions to ensure that the retrieved account matches the created one
	require.Equal(t, account1.ID, account2.ID)             // IDs should match
	require.Equal(t, account1.Owner, account2.Owner)       // Owners should match
	require.Equal(t, account1.Balance, account2.Balance)   // Balances should match
	require.Equal(t, account1.Currency, account2.Currency) // Currencies should match
}

// TestCreateAccountWithInvalidData tests creating an account with invalid data.
func TestCreateAccountWithInvalidData(t *testing.T) {
	// Attempt to create an account with an invalid (empty) owner name
	arg := CreateAccountParams{
		Owner:    "",                    // Invalid empty owner
		Balance:  util.RandomMoney(),    // Valid random balance
		Currency: util.RandomCurrency(), // Valid random currency
	}

	// Attempt to create the account and expect an error
	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.Error(t, err)     // Expect an error due to invalid data
	require.Empty(t, account) // The account should not have been created
}
