package transactions

import (
	"context"
	"database/sql"

	"github.com/goccha/gormsource/pkg"
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
	if v := ctx.Value(pkg.WithTransaction()); v != nil {
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
	if v := ctx.Value(pkg.WithTransaction()); pkg.IsActive(v) {
		return f(ctx, v.(*pkg.TransactionContainer).DB)
	} else {
		return Run(ctx, f, opts...)
	}
}

func Run[T any](ctx context.Context, txFunc func(ctx context.Context, db *gorm.DB) (T, error), opts ...*sql.TxOptions) (res T, err error) {
	if v := ctx.Value(pkg.WithTransaction()); v != nil {
		ctx = context.WithValue(ctx, pkg.WithTransaction(), nil) // 新しいトランザクションをはじめる
	}
	return pkg.RunTransaction[T](ctx, begin, func(ctx context.Context, db *gorm.DB) (context.Context, T, error) {
		ctx = context.WithValue(ctx, pkg.WithTransaction(), &pkg.TransactionContainer{
			DB:              db,
			TransactionType: pkg.Transaction,
		})
		res, err = txFunc(ctx, db)
		return ctx, res, err
	}, pkg.WithTransaction(), opts...)
}

func begin(ctx context.Context, opts ...*sql.TxOptions) *gorm.DB {
	db := getConnection(ctx)
	if len(opts) > 0 {
		return db.Begin(opts...)
	}
	return db.Begin(defaultOptions...)
}

func HandleRollback(ctx context.Context, hook pkg.Hook) {
	pkg.RegisterRollback(ctx, pkg.WithTransaction(), hook)
}

func HandleCommit(ctx context.Context, hook pkg.Hook) {
	pkg.RegisterCommit(ctx, pkg.WithTransaction(), hook)
}
