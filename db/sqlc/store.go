package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Store provides all functions to execute db queries and transactions.
// We need to extend on the exisiting *Queries struct that sqlc provides as it only supports executing queries on one table at a time.
// In order to execute transactions, we will use store to create a set of quesries to be executed in sequence
type Store struct {
	*Queries
	db *sql.DB
}

// NewStore creates a new store
func NewStore(db *sql.DB) *Store {
	return &Store{
		Queries: New(db),
		db:      db,
	}
}

// execTx executes a function within a database transaction.
// This function is added on to the store and it takes context and a callback function as parameters.
// When we start a new transaction, we will create a set of queries and then call the callback function with those queries to execute the transaction.
// If there is a problem, then there is an option to rollback or else the transaction can be committed.
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
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
