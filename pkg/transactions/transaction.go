package transactions

import (
	"context"
	"database/sql"
	"github.com/goccha/gormsource/pkg/foundations"

	"gorm.io/gorm"
)

type contextKey struct {
	key string
}

func (key contextKey) String() string {
	return key.key
}

var transactionSource = contextKey{key: "transactionSource"}

var defaultDB *gorm.DB
var defaultOptions []*sql.TxOptions

type transactionOption struct {
	db      *gorm.DB
	options []*sql.TxOptions
}

func Setup(conn func() (*gorm.DB, error), opt ...*sql.TxOptions) (*gorm.DB, error) {
	if db, err := conn(); err != nil {
		return nil, err
	} else {
		defaultDB = db
		defaultOptions = opt
		return defaultDB, nil
	}
}

func DB(ctx context.Context) *gorm.DB {
	if v := ctx.Value(foundations.WithTransaction()); v != nil {
		return v.(*gorm.DB)
	}
	return getConnection(ctx)
}

func getConnection(ctx context.Context) *gorm.DB {
	if v := ctx.Value(transactionSource); v != nil {
		return v.(*transactionOption).db.WithContext(ctx)
	}
	return defaultDB.WithContext(ctx)
}

func Begin(ctx context.Context, db *gorm.DB, opts ...*sql.TxOptions) context.Context {
	return context.WithValue(ctx, transactionSource, &transactionOption{
		db:      db,
		options: opts,
	})
}

func With[T any](ctx context.Context, f func(ctx context.Context, db *gorm.DB) (T, error), opts ...*sql.TxOptions) (T, error) {
	if v := ctx.Value(foundations.WithTransaction()); foundations.IsActive(v) {
		return f(ctx, v.(*foundations.TransactionContainer).DB)
	} else {
		return Run(ctx, f, opts...)
	}
}

func Run[T any](ctx context.Context, txFunc func(ctx context.Context, db *gorm.DB) (T, error), opts ...*sql.TxOptions) (res T, err error) {
	if v := ctx.Value(foundations.WithTransaction()); v != nil {
		ctx = context.WithValue(ctx, foundations.WithTransaction(), nil) // 新しいトランザクションをはじめる
	}
	return foundations.RunTransaction[T](ctx, begin, func(ctx context.Context, db *gorm.DB) (context.Context, T, error) {
		ctx = context.WithValue(ctx, foundations.WithTransaction(), &foundations.TransactionContainer{
			DB:              db,
			TransactionType: foundations.Transaction,
		})
		res, err = txFunc(ctx, db)
		return ctx, res, err
	}, foundations.WithTransaction(), opts...)
}

func begin(ctx context.Context, opts ...*sql.TxOptions) *gorm.DB {
	db := getConnection(ctx)
	if len(opts) > 0 {
		return db.Begin(opts...)
	}
	return db.Begin(defaultOptions...)
}

func HandleRollback(ctx context.Context, hook foundations.Hook) {
	foundations.RegisterRollback(ctx, foundations.WithTransaction(), hook)
}

func HandleCommit(ctx context.Context, hook foundations.Hook) {
	foundations.RegisterCommit(ctx, foundations.WithTransaction(), hook)
}
