package transactions

import (
	"context"
	"errors"
	"gorm.io/gorm"
)

const (
	withTransaction   = "transactionContext"
	transactionSource = "transactionSource"
)

type Connector func() *gorm.DB

var defaultConnector Connector

func SetDefaultConnector(b Connector) {
	defaultConnector = b
}

func DB(ctx context.Context) *gorm.DB {
	if v := ctx.Value(withTransaction); v == nil {
		return getConnection(ctx)
	} else {
		return v.(*gorm.DB)
	}
}

func getConnection(ctx context.Context) *gorm.DB {
	if v := ctx.Value(transactionSource); v == nil {
		return defaultConnector()
	} else {
		return v.(Connector)()
	}
}

func Begin(ctx context.Context, f Connector) context.Context {
	return context.WithValue(ctx, transactionSource, f)
}

func With(ctx context.Context, f func(ctx context.Context, db *gorm.DB) error) error {
	if v := ctx.Value(withTransaction); v == nil {
		return Run(ctx, f)
	} else {
		return f(ctx, v.(*gorm.DB))
	}
}

func Run(ctx context.Context, txFunc func(ctx context.Context, db *gorm.DB) error) (err error) {
	if v := ctx.Value(withTransaction); v != nil {
		ctx = context.WithValue(ctx, withTransaction, nil) // 新しいトランザクションをはじめる
		return Run(ctx, txFunc)
	} else {
		return runTransaction(ctx, func(ctx context.Context, db *gorm.DB) error {
			ctx = context.WithValue(ctx, withTransaction, db)
			return txFunc(ctx, db)
		})
	}
}

func runTransaction(ctx context.Context, txFunc func(ctx context.Context, db *gorm.DB) error) (err error) {
	db := getConnection(ctx).Begin()
	if db.Error != nil {
		err = db.Error
		return
	}
	defer func() {
		var p interface{}
		if p = recover(); p != nil {
			switch p.(type) {
			case error:
				err = p.(error)
			default:
				err = errors.New("panic")
			}
		}
		if err != nil {
			db.Rollback()
		} else {
			if db = db.Commit(); db.Error != nil {
				err = db.Error
				return
			}
		}
		if p != nil {
			panic(p) // re-throw panic after Rollback
		}
	}()
	err = txFunc(ctx, db)
	if err != nil {
		return
	}
	return nil
}
