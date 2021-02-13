package transactions

import (
	"context"
	"database/sql"
	"github.com/goccha/gormsource/pkg"
	"gorm.io/gorm"
)

const (
	transactionSource = "transactionSource"
)

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
	if v := ctx.Value(pkg.WithTransaction); v != nil {
		return v.(*gorm.DB)
	}
	return getConnection(ctx)
}

func getConnection(ctx context.Context) *gorm.DB {
	if v := ctx.Value(transactionSource); v != nil {
		return v.(*transactionOption).db
	}
	return defaultDB
}

func Begin(ctx context.Context, db *gorm.DB, opts ...*sql.TxOptions) context.Context {
	return context.WithValue(ctx, transactionSource, &transactionOption{
		db:      db,
		options: opts,
	})
}

func With(ctx context.Context, f func(ctx context.Context, db *gorm.DB) error, opts ...*sql.TxOptions) error {
	if v := ctx.Value(pkg.WithTransaction); pkg.IsActive(v) {
		return f(ctx, v.(*gorm.DB))
	} else {
		return Run(ctx, f, opts...)
	}
}

func Run(ctx context.Context, txFunc func(ctx context.Context, db *gorm.DB) error, opts ...*sql.TxOptions) (err error) {
	if v := ctx.Value(pkg.WithTransaction); v != nil {
		ctx = context.WithValue(ctx, pkg.WithTransaction, nil) // 新しいトランザクションをはじめる
	}
	return pkg.RunTransaction(ctx, begin, func(ctx context.Context, db *gorm.DB) error {
		ctx = context.WithValue(ctx, pkg.WithTransaction, db)
		return txFunc(ctx, db)
	}, opts...)
}

func begin(ctx context.Context, opts ...*sql.TxOptions) *gorm.DB {
	db := getConnection(ctx)
	if opts != nil && len(opts) > 0 {
		return db.Begin(opts...)
	}
	return db.Begin(defaultOptions...)
}
