package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/suleimanodetoro/Go-Bank-Pro/db/util"
)

// CreateRandomAccount generates and inserts a random account into the database for test purposes.
// This helper function focuses on setting up data, not performing assertions.
func CreateRandomAccount(t *testing.T) Account {
	// Prepare the parameters to create an account
	arg := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	// Insert the account into the database
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

// TestUpdateAccount tests updating the balance of an existing account.
func TestUpdateAccount(t *testing.T) {
	// Create a random account
	account1 := CreateRandomAccount(t)

	// Set a new balance different from the original one
	newBalance := util.RandomMoney()

	// Prepare the update parameters (update the balance)
	arg := UpdateAccountParams{
		ID:      account1.ID, // The account ID to update
		Balance: newBalance,  // The new balance to set
	}

	// Perform the update
	account2, err := testQueries.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)                                                        // Expect no error after update
	require.NotEmpty(t, account2)                                                  // Assert that the updated account is not empty
	require.Equal(t, newBalance, account2.Balance)                                 // Ensure the balance was updated
	require.Equal(t, account1.ID, account2.ID)                                     // The ID should remain the same
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second) // Ensure CreatedAt wasn't modified
}

// TestUpdateAccountWithInvalidID tests updating an account with an invalid ID.
func TestUpdateAccountWithInvalidID(t *testing.T) {
	// Prepare the update parameters with an invalid ID (0)
	arg := UpdateAccountParams{
		ID:      0,                  // Invalid account ID
		Balance: util.RandomMoney(), // Valid random balance
	}

	// Attempt to update and expect an error
	account, err := testQueries.UpdateAccount(context.Background(), arg)
	require.Error(t, err)     // Expect an error due to invalid ID
	require.Empty(t, account) // The update should not succeed
}

// TestDeleteAccount tests deleting an account and verifying it no longer exists.
func TestDeleteAccount(t *testing.T) {
	// Create a random account to delete
	account1 := CreateRandomAccount(t)

	// Delete the account by its ID
	err := testQueries.DeleteAccount(context.Background(), account1.ID)
	require.NoError(t, err) // Expect no error during deletion

	// Attempt to retrieve the deleted account
	account2, err2 := testQueries.GetAccount(context.Background(), account1.ID)
	require.Error(t, err2)     // Expect an error since the account should no longer exist
	require.Empty(t, account2) // The account object should be empty
}

func TestListAccounts(t *testing.T) {
	// Create 5 random accounts first
	for i := 0; i < 5; i++ {
		CreateRandomAccount(t)
	}

	// Prepare the parameters to list accounts (limiting to 5 results with an offset of 5)
	arg := ListAccountsParams{
		Limit:  5,
		Offset: 0, // Offset should typically start from 0 to list the first 5 results
	}

	// List accounts from the database
	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)     // Expect no error after listing accounts
	require.Len(t, accounts, 5) // Expect exactly 5 accounts to be listed

	// Ensure all listed accounts are not empty
	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}
