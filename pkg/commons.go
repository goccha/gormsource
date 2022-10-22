package pkg

import (
	"context"
	"database/sql"
	"reflect"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type contextKey struct {
	key string
}

func (key contextKey) String() string {
	return key.key
}

var withHook = contextKey{key: "hookContext"}

type Hook func(ctx context.Context)
type Hooks []Hook

func (h Hooks) Invoke(ctx context.Context) {
	for _, f := range h {
		f(ctx)
	}
}

type hookContainer struct {
	rollback Hooks
	commit   Hooks
}

func (c *hookContainer) addRollback(h ...Hook) {
	if c.rollback == nil {
		c.rollback = make(Hooks, 0, len(h))
	}
	c.rollback = append(c.rollback, h...)
}
func (c *hookContainer) addCommit(h ...Hook) {
	if c.commit == nil {
		c.commit = make(Hooks, 0, len(h))
	}
	c.commit = append(c.commit, h...)
}
func (c *hookContainer) invokeRollback(ctx context.Context) {
	if c.rollback != nil {
		c.rollback.Invoke(ctx)
	}
}
func (c *hookContainer) invokeCommit(ctx context.Context) {
	if c.commit != nil {
		c.commit.Invoke(ctx)
	}
}

func fromContext(ctx context.Context) (*hookContainer, bool) {
	if v := ctx.Value(withHook); v != nil {
		v, ok := v.(*hookContainer)
		return v, ok
	}
	return nil, false
}

func EnableHook(ctx context.Context) context.Context {
	if v := ctx.Value(withHook); v == nil {
		c := &hookContainer{}
		return context.WithValue(ctx, withHook, c)
	} else {
		return ctx
	}
}

func RegisterRollback(ctx context.Context, hook ...Hook) {
	if v := ctx.Value(withHook); v != nil {
		if v, ok := v.(*hookContainer); ok {
			v.addRollback(hook...)
		}
	}
}

func RegisterCommit(ctx context.Context, hook ...Hook) {
	if v := ctx.Value(withHook); v != nil {
		if v, ok := v.(*hookContainer); ok {
			v.addCommit(hook...)
		}
	}
}

var withTransaction = contextKey{key: "transactionContext"}

func WithTransaction() interface{} {
	return withTransaction
}

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
			switch p := p.(type) {
			case error:
				err = p
			default:
				err = errors.New("panic")
			}
		}
		if err != nil {
			db.Rollback()
			if f, ok := fromContext(ctx); ok {
				f.invokeRollback(ctx)
			}
		} else {
			if db = db.Commit(); db.Error != nil {
				err = db.Error
				return
			}
			if f, ok := fromContext(ctx); ok {
				f.invokeCommit(ctx)
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
