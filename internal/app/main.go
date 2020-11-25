package main

import (
	"context"
	"fmt"
	"github.com/goccha/gormsource/internal/examples"
	"github.com/goccha/gormsource/pkg/gormsource"
	"gorm.io/gorm"
)

func main() {
	if mydb, err := examples.InitMysql(); err != nil {
		panic(err)
	} else {
		gormsource.SetDefaultConnector(func() *gorm.DB {
			return mydb
		})
		if err := examples.Migrate(mydb); err != nil {
			panic(err)
		}
	}
	if pdb, err := examples.InitPosgres(); err != nil {
		panic(err)
	} else {
		if err := examples.Migrate(pdb); err != nil {
			panic(err)
		}
		var ctx = context.Background()
		if err := gormsource.RunTransaction(ctx, func(ctx context.Context, db *gorm.DB) error {
			entity := examples.ExampleTable{
				ID:   "key1",
				Desc: "test01",
			}
			db = db.Create(&entity)
			if db.Error != nil {
				return db.Error
			}
			if err := gormsource.RunTransaction(gormsource.Begin(ctx, func() *gorm.DB {
				return pdb
			}), func(ctx context.Context, db *gorm.DB) error {
				entity := examples.ExampleTable{
					ID:   "key2",
					Desc: "test02",
				}
				db = db.Create(&entity)
				if db.Error != nil {
					return db.Error
				}
				if entity2, err := examples.GetPostgresEntity(ctx, entity.ID); err != nil {
					return err
				} else if entity2 == nil {
					return fmt.Errorf("%s not found", entity.ID)
				}
				if entity3, err := examples.GetPostgresEntity(context.Background(), entity.ID); err != nil {
					return err
				} else if entity3 != nil {
					return fmt.Errorf("invalid transaction")
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
			return nil
		}); err != nil {
			panic(err)
		}
	}
}
