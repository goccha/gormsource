package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/goccha/gormsource/internal/examples"
	"github.com/goccha/gormsource/pkg/replicas"
	"github.com/goccha/gormsource/pkg/transactions"
	"gorm.io/gorm"
	"time"
)

func main() {
	if mydb, err := transactions.Setup(func() (*gorm.DB, error) {
		return examples.InitPrimaryMysql()
	}, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
	}); err != nil {
		panic(err)
	} else {
		if err := examples.Migrate(mydb); err != nil {
			panic(err)
		}
	}
	if _, err := replicas.Setup(func() (*gorm.DB, error) {
		return examples.InitReplicaMysql()
	}); err != nil {
		panic(err)
	}
	if pdb, err := examples.InitPrimaryPosgres(); err != nil {
		panic(err)
	} else {
		var rpdb *replicas.DB
		if rpdb, err = replicas.New(func() (*gorm.DB, error) {
			return examples.InitReplicaPosgres()
		}); err != nil {
			panic(err)
		}
		if err := examples.Migrate(pdb); err != nil {
			panic(err)
		}
		var ctx = context.Background()
		if err := transactions.Run(ctx, func(ctx context.Context, db *gorm.DB) error {
			entity := examples.ExampleTable{
				ID:   "key1",
				Desc: "test01",
			}
			db = db.Create(&entity)
			if db.Error != nil {
				return db.Error
			}
			if err := transactions.Run(transactions.Begin(ctx, pdb), func(ctx context.Context, db *gorm.DB) error {
				entity := examples.ExampleTable{
					ID:   "key2",
					Desc: "test02",
				}
				db = db.Create(&entity)
				if db.Error != nil {
					return db.Error
				}
				if entity2, err := examples.GetPrimaryEntity(ctx, entity.ID); err != nil {
					return err
				} else if entity2 == nil {
					return fmt.Errorf("%s not found", entity.ID)
				}
				if entity3, err := examples.GetPrimaryEntity(context.Background(), entity.ID); err != nil {
					return err
				} else if entity3 != nil {
					return fmt.Errorf("invalid transaction")
				}
				if err = replicas.WithTransaction(ctx, func(ctx context.Context, db *gorm.DB) error {
					if entity, err := examples.GetReplicaEntity(db, entity.ID); err != nil {
						return err
					} else if entity == nil {
						return fmt.Errorf("replica %s not found", entity.ID)
					}
					return nil
				}); err != nil {
					return err
				}
				if err = replicas.With(replicas.Begin(ctx, rpdb), func(ctx context.Context, db *gorm.DB) error {
					if entity, err := examples.GetReplicaEntity(db, ""); err != nil {
						return err
					} else if entity != nil {
						return fmt.Errorf("invalid transaction")
					}
					return nil
				}); err != nil {
					return err
				}
				return nil
			}); err != nil {
				return err
			}
			entity = examples.ExampleTable{
				ID:   "key3",
				Desc: "test03",
			}
			db = db.Create(&entity)
			if db.Error != nil {
				return db.Error
			}
			time.Sleep(1 * time.Second)
			if err = replicas.With(ctx, func(ctx context.Context, db *gorm.DB) error {
				if entity3, err := examples.GetReplicaEntity(db, entity.ID); err != nil {
					return err
				} else if entity3 != nil {
					return fmt.Errorf("invalid transaction")
				}
				return nil
			}); err != nil {
				panic(err)
			}
			return nil
		}); err != nil {
			panic(err)
		}
		time.Sleep(1 * time.Second)
		if err = replicas.Run(ctx, func(ctx context.Context, db *gorm.DB) error {
			if entity3, err := examples.GetReplicaEntity(db, "key1"); err != nil {
				return err
			} else if entity3 == nil {
				return fmt.Errorf("replica %s not found", "key1")
			}
			return nil
		}); err != nil {
			panic(err)
		}
	}
}
