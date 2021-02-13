package examples

import (
	"context"
	"github.com/goccha/errors"
	"github.com/goccha/gormsource/pkg/transactions"
	"gorm.io/gorm"
)

func GetPrimaryEntity(ctx context.Context, id string) (*ExampleTable, error) {
	db := transactions.DB(ctx)
	entity := &ExampleTable{
		ID: id,
	}
	db = db.First(entity)
	if errors.Is(db.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if db.Error != nil {
		return nil, db.Error
	}
	return entity, nil
}

func GetReplicaEntity(db *gorm.DB, id string) (*ExampleTable, error) {
	entity := &ExampleTable{
		ID: id,
	}
	db = db.First(entity)
	if errors.Is(db.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if db.Error != nil {
		return nil, db.Error
	}
	return entity, nil
}
