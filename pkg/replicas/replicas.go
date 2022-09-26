package replicas

import (
	"context"
	"database/sql"
	"sync"
	"sync/atomic"

	"github.com/goccha/gormsource/pkg"
	"github.com/goccha/log"
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

func With(ctx context.Context, f func(ctx context.Context, db *gorm.DB) error) error {
	if v := ctx.Value(withReadOnly); v != nil {
		return f(ctx, v.(*gorm.DB))
	} else {
		return Run(ctx, f)
	}
}

func WithTransaction(ctx context.Context, f func(ctx context.Context, db *gorm.DB) error) error {
	if v := ctx.Value(withReadOnly); pkg.IsActive(v) {
		return f(ctx, v.(*gorm.DB))
	} else if v = ctx.Value(pkg.WithTransaction()); v != nil {
		return f(ctx, v.(*gorm.DB))
	} else {
		return Run(ctx, f)
	}
}

func Run(ctx context.Context, f func(ctx context.Context, db *gorm.DB) error) error {
	if v := ctx.Value(withReadOnly); v != nil {
		ctx = context.WithValue(ctx, withReadOnly, nil) // 新しいトランザクションをはじめる
	}
	return pkg.RunTransaction(ctx, begin, func(ctx context.Context, db *gorm.DB) error {
		ctx = context.WithValue(ctx, withReadOnly, db)
		return f(ctx, db)
	})
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
