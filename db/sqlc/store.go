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

var txKey = struct{}{}

func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult
	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		txName := ctx.Value(txKey)

		fmt.Println(txName, "create transfer")
		result.Transfer, err = q.createTransfer(ctx, createTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})

		if err != nil {
			return err
		}
		fmt.Println(txName, "create entry")
		result.FromEntry, err = q.createEntry(ctx, createEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})

		if err != nil {
			return err
		}

		fmt.Println(txName, "create entry")
		result.ToEntry, err = q.createEntry(ctx, createEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})

		if err != nil {
			return err
		}

		// Update account balances

		// fmt.Println(txName, "get account 1")
		// account1, err := q.getAccountForUpdate(ctx, arg.FromAccountID)
		// if err != nil {
		// 	return err
		// }

		fmt.Println(txName, "update account 2")
		result.FromAccount, err = q.updateAccountBalance(ctx, updateAccountBalanceParams{
			ID:      arg.FromAccountID,
			Balance: -arg.Amount,
		})

		if err != nil {
			return err
		}

		// fmt.Println(txName, "get account 2")
		// account2, err := q.getAccountForUpdate(ctx, arg.ToAccountID)
		// if err != nil {
		// 	return err
		// }

		fmt.Println(txName, "update account 2")
		result.ToAccount, err = q.updateAccountBalance(ctx, updateAccountBalanceParams{
			ID:      arg.ToAccountID,
			Balance: arg.Amount,
		})

		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}
