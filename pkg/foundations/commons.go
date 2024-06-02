package foundations

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

type Hook func(ctx context.Context)
type Hooks []Hook

func (h Hooks) Invoke(ctx context.Context) {
	for _, f := range h {
		f(ctx)
	}
}

func fromContext(ctx context.Context, key any) (*TransactionContainer, bool) {
	if v := ctx.Value(key); v != nil {
		v, ok := v.(*TransactionContainer)
		return v, ok
	}
	return &TransactionContainer{}, false
}

func RegisterRollback(ctx context.Context, key any, hook ...Hook) {
	if v := ctx.Value(key); v != nil {
		if v, ok := v.(*TransactionContainer); ok {
			v.addRollback(hook...)
		}
	}
}

func RegisterCommit(ctx context.Context, key any, hook ...Hook) {
	if v := ctx.Value(key); v != nil {
		if v, ok := v.(*TransactionContainer); ok {
			v.addCommit(hook...)
		}
	}
}

const (
	ReadOnly    = "readOnly"
	Transaction = "transaction"
)

type TransactionContainer struct {
	DB              *gorm.DB
	TransactionType string
	rollback        Hooks
	commit          Hooks
}

func (c *TransactionContainer) addRollback(h ...Hook) {
	if c.rollback == nil {
		c.rollback = make(Hooks, 0, len(h))
	}
	c.rollback = append(c.rollback, h...)
}
func (c *TransactionContainer) addCommit(h ...Hook) {
	if c.commit == nil {
		c.commit = make(Hooks, 0, len(h))
	}
	c.commit = append(c.commit, h...)
}
func (c *TransactionContainer) invokeRollback(ctx context.Context) {
	if c.rollback != nil {
		c.rollback.Invoke(ctx)
	}
}
func (c *TransactionContainer) invokeCommit(ctx context.Context) {
	if c.commit != nil {
		c.commit.Invoke(ctx)
	}
}

var withTransaction = contextKey{key: "transactionContext"}

func WithTransaction() interface{} {
	return withTransaction
}

type Begin func(ctx context.Context, opts ...*sql.TxOptions) *gorm.DB

func IsActive(v interface{}) bool {
	if container, ok := v.(*TransactionContainer); ok {
		if committer, ok := container.DB.Statement.ConnPool.(gorm.TxCommitter); ok &&
			committer != nil && !reflect.ValueOf(committer).IsNil() {
			return true
		}
	}
	return false
}

func RunTransaction[T any](ctx context.Context, begin Begin, txFunc func(ctx context.Context, db *gorm.DB) (context.Context, T, error), key any, opts ...*sql.TxOptions) (res T, err error) {
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
			if f, ok := fromContext(ctx, key); ok {
				f.invokeRollback(ctx)
			}
		} else {
			if db = db.Commit(); db.Error != nil {
				err = db.Error
				return
			}
			if f, ok := fromContext(ctx, key); ok {
				f.invokeCommit(ctx)
			}
		}
		if p != nil {
			panic(p) // re-throw panic after Rollback
		}
	}()
	ctx, res, err = txFunc(ctx, db)
	if err != nil {
		return
	}
	return res, nil
}
