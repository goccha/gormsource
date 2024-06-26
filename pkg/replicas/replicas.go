package replicas

import (
	"context"
	"database/sql"
	"github.com/goccha/gormsource/pkg/foundations"
	"sync"
	"sync/atomic"

	"github.com/goccha/envar/pkg/log"
	"gorm.io/gorm"
)

type contextKey struct {
	key string
}

func (key contextKey) String() string {
	return key.key
}

var withReadOnly = contextKey{key: "readOnlyTransaction"}
var replicaSource = contextKey{key: "replicaSource"}

var replicaOption = &sql.TxOptions{
	ReadOnly: true,
}

type Connector func() (*gorm.DB, error)

var defaultReplica *DB

type DB struct {
	dbs     []*gorm.DB
	counter *CyclicCounter
}

func (db *DB) DB() *gorm.DB {
	return db.dbs[db.counter.next()]
}

func (db *DB) Use(plugin gorm.Plugin) error {
	for _, v := range db.dbs {
		if err := v.Use(plugin); err != nil {
			return err
		}
	}
	return nil
}

func (db *DB) Close() {
	for _, d := range db.dbs {
		if sqlDB, err := d.DB(); err != nil {
			log.Warn("%v", err)
		} else {
			if err := sqlDB.Close(); err != nil {
				log.Warn("%v", err)
			}
		}
	}
}

func Setup(connectors ...Connector) (*DB, error) {
	if db, err := New(connectors...); err != nil {
		return nil, err
	} else {
		defaultReplica = db
		return db, nil
	}
}

func New(connectors ...Connector) (*DB, error) {
	dbs := make([]*gorm.DB, 0)
	for _, c := range connectors {
		if db, err := c(); err != nil {
			return nil, err
		} else {
			dbs = append(dbs, db)
		}
	}
	return &DB{
		dbs: dbs,
		counter: &CyclicCounter{
			mu:  &sync.RWMutex{},
			max: int32(len(dbs)),
			cnt: -1,
		},
	}, nil
}

func getConnection(ctx context.Context) *gorm.DB {
	if v := ctx.Value(replicaSource); v != nil {
		return v.(*DB).DB().WithContext(ctx)
	}
	return defaultReplica.DB().WithContext(ctx)
}

func Begin(ctx context.Context, db *DB) context.Context {
	return context.WithValue(ctx, replicaSource, db)
}

func With[T any](ctx context.Context, f func(ctx context.Context, db *gorm.DB) (T, error)) (T, error) {
	if v := ctx.Value(withReadOnly); v != nil {
		return f(ctx, v.(*foundations.TransactionContainer).DB)
	} else {
		return Run(ctx, f)
	}
}

func WithTransaction[T any](ctx context.Context, f func(ctx context.Context, db *gorm.DB) (T, error)) (T, error) {
	if v := ctx.Value(withReadOnly); foundations.IsActive(v) {
		return f(ctx, v.(*foundations.TransactionContainer).DB)
	} else if v = ctx.Value(foundations.WithTransaction()); v != nil {
		return f(ctx, v.(*foundations.TransactionContainer).DB)
	} else {
		return Run(ctx, f)
	}
}

func Run[T any](ctx context.Context, f func(ctx context.Context, db *gorm.DB) (T, error), opts ...*sql.TxOptions) (res T, err error) {
	if v := ctx.Value(withReadOnly); v != nil {
		ctx = context.WithValue(ctx, withReadOnly, nil) // 新しいトランザクションをはじめる
	}
	return foundations.RunTransaction[T](ctx, begin, func(ctx context.Context, db *gorm.DB) (context.Context, T, error) {
		ctx = context.WithValue(ctx, withReadOnly, &foundations.TransactionContainer{
			DB:              db,
			TransactionType: foundations.ReadOnly,
		})
		res, err = f(ctx, db)
		return ctx, res, err
	}, withReadOnly, opts...)
}

func begin(ctx context.Context, _ ...*sql.TxOptions) *gorm.DB {
	db := getConnection(ctx)
	return db.Begin(replicaOption)
}

type CyclicCounter struct {
	mu  *sync.RWMutex
	max int32
	cnt int32
}

func (c *CyclicCounter) next() int32 {
	if c.max == 1 {
		return 0
	}
	c.mu.RLock()
	v := atomic.AddInt32(&c.cnt, 1)
	if v >= c.max {
		c.mu.RUnlock()
		c.mu.Lock()
		defer c.mu.Unlock()
		c.cnt = -1
		return 0
	} else {
		c.mu.RUnlock()
	}
	return c.cnt
}

func HandleRollback(ctx context.Context, hook foundations.Hook) {
	foundations.RegisterRollback(ctx, withReadOnly, hook)
}

func HandleCommit(ctx context.Context, hook foundations.Hook) {
	foundations.RegisterCommit(ctx, withReadOnly, hook)
}
