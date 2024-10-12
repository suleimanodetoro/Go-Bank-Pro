package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Store interface combines the SQLC-generated Querier interface and
// custom transaction methods. This interface allows mocking of the store
// for testing and separates the logic from specific implementations.
type Store interface {
	Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
}

// SQLStore implements the Store interface, providing methods to interact
// with the database and handle transactions.
type SQLStore struct {
	*Queries
	db *sql.DB
}

// NewStore initializes a new SQLStore and returns it as a Store interface.
// This allows the returned store to satisfy the Store interface, enabling flexibility
// for testing, mocking, and easier future modifications.
func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

// execTx executes a function within a database transaction.
// It begins a transaction, passes a Queries object tied to the transaction
// to the provided function, and commits or rolls back based on the function's success.
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil) // Start a new transaction
	if err != nil {
		return err
	}

	q := New(tx) // Create a new Queries object tied to the transaction
	err = fn(q)  // Execute the provided function, passing in the transaction-bound Queries object

	if err != nil {
		// Roll back if an error occurs
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx error: %v, rollback error: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit() // Commit the transaction on success
}

// TransferTxParams contains all input parameters to transfer money between two accounts.
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

// TransferTxResult contains the result of a successful transfer transaction,
// including the transfer record, the updated account balances, and entries for both accounts.
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

// txKey is used to track the context of the current transaction.
// This key can be used to identify and label transactions in logs or tracking systems.
var txKey = struct{}{}

// TransferTx performs a money transfer between two accounts, ensuring that the operation is atomic and safe.
// It creates the necessary transfer and entry records and updates the accounts' balances.
func (store *SQLStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		txName := ctx.Value(txKey)

		fmt.Println(txName, "Create transfer")
		// Create the transfer record
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		// Add entries for both accounts
		fmt.Println(txName, "Create entry 1")
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		fmt.Println(txName, "Create entry 2")
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		// Update the account balances, ensuring that the account with the smaller ID is updated first to avoid deadlocks.
		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, result.ToAccount, err = addMoney(ctx, q, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)
			if err != nil {
				return err
			}
		} else {
			result.FromAccount, result.ToAccount, err = addMoney(ctx, q, arg.ToAccountID, arg.Amount, arg.FromAccountID, -arg.Amount)
			if err != nil {
				return err
			}
		}

		return nil
	})

	return result, err
}

// addMoney updates the balances of two accounts as part of the transfer transaction.
// It ensures that each account's balance is modified correctly based on the transfer amount.
func addMoney(
	ctx context.Context,
	q *Queries,
	accountID1 int64,
	amount1 int64,
	accountID2 int64,
	amount2 int64,
) (account1 Account, account2 Account, err error) {
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID1,
		Amount: amount1,
	})
	if err != nil {
		return account1, account2, err
	}

	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID2,
		Amount: amount2,
	})
	if err != nil {
		return account1, account2, err
	}

	return account1, account2, nil
}
