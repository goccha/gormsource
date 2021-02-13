package pkg

import (
	"context"
	"database/sql"
	"github.com/goccha/errors"
	"gorm.io/gorm"
	"reflect"
)

const (
	WithTransaction = "transactionContext"
)

type Begin func(ctx context.Context, opts ...*sql.TxOptions) *gorm.DB

func IsActive(v interface{}) bool {
	if db, ok := v.(*gorm.DB); ok {
		if committer, ok := db.Statement.ConnPool.(gorm.TxCommitter); ok &&
			committer != nil && !reflect.ValueOf(committer).IsNil() {
			return true
		}
	}
	return false
}

func RunTransaction(ctx context.Context, begin Begin, txFunc func(ctx context.Context, db *gorm.DB) error, opts ...*sql.TxOptions) (err error) {
	db := begin(ctx, opts...)
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
