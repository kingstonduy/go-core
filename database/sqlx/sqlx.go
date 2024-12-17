package sqlx

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/kingstonduy/go-core/database"
	"github.com/uptrace/opentelemetry-go-extra/otelsqlx"
)

// SqlxDBx is the sqlx.DB based implementation of GDBC
type SqlxDB struct {
	db *sqlx.DB
}

func NewSqlxGdbc(driverName string, dsn string, opts ...database.DatabaseOption) (*database.Gdbc, error) {

	db := otelsqlx.MustConnect(driverName, dsn)

	err := db.Ping()
	if err != nil {
		if db != nil {
			err = db.Close()
		}
		return nil, err
	}

	options := database.NewDatabaseOptions(opts...)
	if options.MaxIdleCount > 0 {
		db.SetMaxIdleConns(options.MaxIdleCount)
	}

	if options.MaxOpen > 0 {
		db.SetMaxOpenConns(options.MaxOpen)
	}

	if options.MaxLifetime > 0 {
		db.SetConnMaxLifetime(options.MaxLifetime)
	}

	if options.MaxIdleTime > 0 {
		db.SetConnMaxIdleTime(options.MaxIdleTime)
	}

	return &database.Gdbc{
		Executor: &SqlxDB{
			db: db,
		},
		Options: options,
	}, nil
}

// Get implements database.SqlGdbc.
func (s *SqlxDB) Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return s.db.GetContext(ctx, dest, query, args...)
}

// Select implements database.SqlGdbc.
func (s *SqlxDB) Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return s.db.SelectContext(ctx, dest, query, args...)
}

// Exec implements SqlGdbc.
func (s *SqlxDB) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return s.db.ExecContext(ctx, query, args...)
}

// Prepare implements SqlGdbc.
func (s *SqlxDB) Prepare(ctx context.Context, query string) (*sql.Stmt, error) {
	return s.db.PrepareContext(ctx, query)
}

// Query implements SqlGdbc.
func (s *SqlxDB) Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return s.db.QueryContext(ctx, query, args...)
}

// QueryRow implements SqlGdbc.
func (s *SqlxDB) QueryRow(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return s.db.QueryRowContext(ctx, query, args...)
}

func (sdt *SqlxDB) WithinTransaction(ctx context.Context, txFunc func(ctx context.Context) error, opts ...database.TransactionOption) error {
	var err error
	var tx *sqlx.Tx

	if len(opts) != 0 {
		txOptions := database.NewTransactionOptions(opts...)
		sqlTxOptions := &sql.TxOptions{
			Isolation: txOptions.Isolation,
			ReadOnly:  txOptions.ReadOnly,
		}

		tx, err = sdt.db.BeginTxx(ctx, sqlTxOptions)
	} else {
		tx, err = sdt.db.Beginx()
	}

	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	sct := &SqlxTx{
		db:        tx,
		wrapperDB: sdt.db,
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	err = txFunc(database.InjectTx(ctx, sct))
	return err
}

func (s *SqlxDB) Stats(ctx context.Context) sql.DBStats {
	return s.db.Stats()
}

func (s *SqlxDB) Close(ctx context.Context) {
	s.db.Close()
}

// SqlxConnx is the sqlx.Tx based implementation of GDBC
type SqlxTx struct {
	db        *sqlx.Tx
	wrapperDB *sqlx.DB
}

// Exec implements SqlGdbc.
func (s *SqlxTx) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return s.db.ExecContext(ctx, query, args...)
}

// Prepare implements SqlGdbc.
func (s *SqlxTx) Prepare(ctx context.Context, query string) (*sql.Stmt, error) {
	return s.db.PrepareContext(ctx, query)
}

// Query implements SqlGdbc.
func (s *SqlxTx) Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return s.db.QueryContext(ctx, query, args...)
}

// QueryRow implements SqlGdbc.
func (s *SqlxTx) QueryRow(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return s.db.QueryRowContext(ctx, query, args...)
}

// Get implements database.SqlGdbc.
func (s *SqlxTx) Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return s.db.GetContext(ctx, dest, query, args...)
}

// Select implements database.SqlGdbc.
func (s *SqlxTx) Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return s.db.SelectContext(ctx, dest, query, args...)
}

func (sct *SqlxTx) WithinTransaction(ctx context.Context, txFunc func(ctx context.Context) error, opts ...database.TransactionOption) error {
	var err error
	tx := sct.db
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	err = txFunc(database.InjectTx(ctx, sct))
	return err
}

func (s *SqlxTx) Stats(ctx context.Context) sql.DBStats {
	return s.wrapperDB.Stats()
}

func (s *SqlxTx) Close(ctx context.Context) {
	// do nothing
}
