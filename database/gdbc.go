package database

import (
	"context"
	"database/sql"
)

// SqlGdbc (SQL Go database connection) is a wrapper for SQL database handler ( can be *sql.DB or *sql.Tx)
// It should be able to work with all SQL data that follows SQL standard.
type SqlGdbc interface {
	Transactor
	Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	Prepare(ctx context.Context, query string) (*sql.Stmt, error)
	Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(ctx context.Context, query string, args ...interface{}) *sql.Row
	Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	Stats(ctx context.Context) sql.DBStats
	Close(ctx context.Context)
}

type DBStats struct {
}

// Used this in repositories
type Gdbc struct {
	Executor SqlGdbc
	Options  DatabaseOptions
}

// Exec implements SqlGdbc.
func (g *Gdbc) Exec(ctx context.Context, query string, args ...interface{}) (res sql.Result, err error) {
	ctx, _ = g.beforeHook(ctx, query, args...)
	defer func() {
		g.afterHook(ctx, err, query, args...) //nolint
	}()

	return g.getConnection(ctx).Exec(ctx, query, args...)
}

// Get implements SqlGdbc.
func (g *Gdbc) Get(ctx context.Context, dest interface{}, query string, args ...interface{}) (err error) {
	ctx, _ = g.beforeHook(ctx, query, args...)
	defer func() {
		g.afterHook(ctx, err, query, args...) //nolint
	}()

	return g.getConnection(ctx).Get(ctx, dest, query, args...)
}

// Prepare implements SqlGdbc.
func (g *Gdbc) Prepare(ctx context.Context, query string) (*sql.Stmt, error) {
	return g.getConnection(ctx).Prepare(ctx, query)
}

// Query implements SqlGdbc.
func (g *Gdbc) Query(ctx context.Context, query string, args ...interface{}) (rows *sql.Rows, err error) {
	ctx, _ = g.beforeHook(ctx, query, args...)
	defer func() {
		g.afterHook(ctx, err, query, args...) //nolint
	}()

	return g.getConnection(ctx).Query(ctx, query, args...)
}

// QueryRow implements SqlGdbc.
func (g *Gdbc) QueryRow(ctx context.Context, query string, args ...interface{}) (row *sql.Row) {
	ctx, _ = g.beforeHook(ctx, query, args...)
	defer func() {
		g.afterHook(ctx, nil, query, args...) //nolint
	}()

	return g.getConnection(ctx).QueryRow(ctx, query, args...)
}

// Select implements SqlGdbc.
func (g *Gdbc) Select(ctx context.Context, dest interface{}, query string, args ...interface{}) (err error) {
	ctx, _ = g.beforeHook(ctx, query, args...)
	defer func() {
		g.afterHook(ctx, err, query, args...) //nolint
	}()

	return g.getConnection(ctx).Select(ctx, dest, query, args...)
}

func (g *Gdbc) WithinTransaction(ctx context.Context, txFunc func(ctx context.Context) error, options ...TransactionOption) error {
	return g.getConnection(ctx).WithinTransaction(ctx, txFunc, options...)
}

func (g *Gdbc) Stats(ctx context.Context) sql.DBStats {
	return g.getConnection(ctx).Stats(ctx)
}

func (g *Gdbc) Close(ctx context.Context) {
	// close the executor
	if g.Executor != nil {
		g.Executor.Close(ctx)
	}

	// close the connection injected in the context
	if cnn := g.getConnection(ctx); cnn != nil {
		g.getConnection(ctx)
	}
}

func (g *Gdbc) beforeHook(ctx context.Context, query string, args ...interface{}) (context.Context, error) {
	if hooks := g.Options.Hooks; hooks != nil {
		return hooks.Before(ctx, query, args...)
	}
	return ctx, nil
}

func (g *Gdbc) afterHook(ctx context.Context, err error, query string, args ...interface{}) (context.Context, error) {
	if hooks := g.Options.Hooks; hooks != nil {
		return hooks.After(ctx, err, query, args...)
	}
	return ctx, nil
}

func (g *Gdbc) getConnection(ctx context.Context) SqlGdbc {
	s := ExtractTx(ctx)
	if s != nil {
		return s
	}
	return g.Executor
}
