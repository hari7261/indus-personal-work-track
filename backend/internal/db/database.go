package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/glebarez/sqlite"
	"github.com/jmoiron/sqlx"
)

type Config struct {
	DSN string
}

type DB struct {
	*sqlx.DB
}

func NewDatabase(config Config) (*DB, error) {
	db, err := sqlx.Connect("sqlite", config.DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := migrate(db); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return &DB{db}, nil
}

func (db *DB) BeginTx(ctx context.Context) (*sqlx.Tx, error) {
	return db.BeginTxx(ctx, nil)
}

func (db *DB) BeginTxx(ctx context.Context, opts *sql.TxOptions) (*sqlx.Tx, error) {
	return db.DB.BeginTxx(ctx, opts)
}

func (db *DB) WithTransaction(fn func(*sqlx.Tx) error) error {
	tx, err := db.BeginTxx(context.Background(), nil)
	if err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("rollback failed: %v (original error: %w)", rbErr, err)
		}
		return err
	}

	return tx.Commit()
}
