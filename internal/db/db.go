package db

import (
	"context"

	"github.com/SecureParadise/go_attendence/internal/db/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Store defines all functions to execute db queries and transactions
type Store interface {
	sqlc.Querier
	WithTx(ctx context.Context, fn func(*sqlc.Queries) error) error
}

// SQLStore provides all functions to execute SQL queries and transactions
type SQLStore struct {
	connPool *pgxpool.Pool
	*sqlc.Queries
}

// NewStore creates a new store
func NewStore(connPool *pgxpool.Pool) Store {
	return &SQLStore{
		connPool: connPool,
		Queries:  sqlc.New(connPool),
	}
}

// WithTx executes a transaction
func (store *SQLStore) WithTx(ctx context.Context, fn func(*sqlc.Queries) error) error {
	tx, err := store.connPool.Begin(ctx)
	if err != nil {
		return err
	}

	q := sqlc.New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return rbErr
		}
		return err
	}

	return tx.Commit(ctx)
}

// Implement the Querier interface
func (store *SQLStore) CreateUser(ctx context.Context, arg sqlc.CreateUserParams) (sqlc.User, error) {
	return store.Queries.CreateUser(ctx, arg)
}
