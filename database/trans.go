package database

import (
	"context"
)

type Transactor interface {
	WithinTransaction(ctx context.Context, txFunc func(ctx context.Context) error, options ...TransactionOption) error
}

type TxKey struct{}

func InjectTx(ctx context.Context, s SqlGdbc) context.Context {
	return context.WithValue(ctx, TxKey{}, s)
}

func ExtractTx(ctx context.Context) SqlGdbc {
	if tx, ok := ctx.Value(TxKey{}).(SqlGdbc); ok {
		return tx
	}
	return nil
}

type EnableTransactor interface {
	WithinTransaction(ctx context.Context, txFunc func(ctx context.Context) error, options ...TransactionOption) error
}
