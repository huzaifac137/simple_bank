package db

import (
	"context"
	"database/sql"
	"fmt"
)

// store provides all functions to execute database queries and transactions
type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		Queries: New(db),
		db:      db,
	}
}

// execute function within database transaction

func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)

	err = fn(q)
	if err != nil {
		if txRBErr := tx.Rollback(); txRBErr != nil {
			return fmt.Errorf("tx error: %v, tx rollback error: %v", err, txRBErr)
		}
		return err
	}
	return tx.Commit()
}

// TransferTx performs money transfer transactions for queries

type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferTxResult struct {
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	Transfer    Transfer `json:"transfer"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult
	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		result.Transfer, err = q.createTransfer(ctx, createTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})

		if err != nil {
			return err
		}

		result.FromEntry, err = q.createEntry(ctx, createEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})

		if err != nil {
			return err
		}

		result.ToEntry, err = q.createEntry(ctx, createEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})

		if err != nil {
			return err
		}

		// Update account balances

		return nil
	})

	return result, err
}
