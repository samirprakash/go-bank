package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Store provides all functions to execute db queries and transactions.
// We need to extend on the exisiting *Queries struct that sqlc provides as it only supports executing queries on one table at a time.
// In order to execute transactions, we will use store to create a set of quesries to be executed in sequence
type Store interface {
	Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
}

// SQLStore provides all functions to execute SQL queries and transactions.
// We need to extend on the exisiting *Queries struct that sqlc provides as it only supports executing queries on one table at a time.
// In order to execute transactions, we will use store to create a set of quesries to be executed in sequence
type SQLStore struct {
	db *sql.DB
	*Queries
}

// NewStore creates a new store
func NewStore(db *sql.DB) Store {
	return &SQLStore{
		Queries: New(db),
		db:      db,
	}
}

// execTx executes a function within a database transaction.
// This function is added on to the store and it takes context and a callback function as parameters.
// When we start a new transaction, we will create a set of queries and then call the callback function with those queries to execute the transaction.
// If there is a problem, then there is an option to rollback or else the transaction can be committed.
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		// rollback if there is an error
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("transaction error : %v, rollback error: %v", err, rbErr)
		}
		return err
	}
	// commit if everything works
	return tx.Commit()
}

// TransferTxParams represents the arguments required to execute the transfer transaction
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id,omitempty"`
	ToAccountID   int64 `json:"to_account_id,omitempty"`
	Amount        int64 `json:"amount,omitempty"`
}

// TransferTxResult represents the result of the transfer transaction
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer,omitempty"`
	FromAccount Account  `json:"from_account,omitempty"`
	ToAccount   Account  `json:"to_account,omitempty"`
	FromEntry   Entry    `json:"from_entry,omitempty"`
	ToEntry     Entry    `json:"to_entry,omitempty"`
}

// TransferTx performs a money transfer from one account to another.
// It executes a database transaction to create a transfer record, update account entries and update the account balance.
func (store *SQLStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		// create a transfer query and execute it
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		// create an entry for the account from which money has been transferred
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		// create an entry for the acoount to which the money has been transferred
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		//Update account balance
		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, result.ToAccount, err = updateBalancesToAccounts(ctx, q, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)
		} else {
			result.ToAccount, result.FromAccount, err = updateBalancesToAccounts(ctx, q, arg.ToAccountID, arg.Amount, arg.FromAccountID, -arg.Amount)
		}
		return nil
	})

	return result, err
}

func updateBalancesToAccounts(ctx context.Context, q *Queries, accountID1 int64, amount1 int64, accountID2 int64, amount2 int64) (account1 Account, account2 Account, err error) {
	account1, err = q.UpdateAccountBalance(ctx, UpdateAccountBalanceParams{
		Amount: amount1,
		ID:     accountID1,
	})
	if err != nil {
		return
	}

	account2, err = q.UpdateAccountBalance(ctx, UpdateAccountBalanceParams{
		Amount: amount2,
		ID:     accountID2,
	})
	if err != nil {
		return
	}

	return
}
