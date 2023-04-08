package pg

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type DB struct {
	db  *sqlx.DB
	DSN string
	Now func() time.Time
}

func NewDB(dsn string) *DB {
	return &DB{DSN: dsn, Now: time.Now}
}

func (db *DB) Open() error {
	if db.DSN == "" {
		return fmt.Errorf("dsn required")
	}
	var err error
	db.db, err = sqlx.Open("postgres", db.DSN)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) Close() error {
	if db.db != nil {
		return db.db.Close()
	}
	return nil
}

type Tx struct {
	*sqlx.Tx
	db  *DB
	now time.Time
}

func (db *DB) BeginTx(ctx context.Context) (*Tx, error) {
	tx, err := db.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &Tx{
		Tx:  tx,
		db:  db,
		now: db.Now().UTC().Truncate(time.Second),
	}, nil
}
