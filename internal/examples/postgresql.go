package examples

import (
	"context"
	"errors"
	"github.com/goccha/gormsource/pkg/datasources"
	"github.com/goccha/gormsource/pkg/datasources/postgresql"
	"github.com/goccha/gormsource/pkg/gormsource"
	"gorm.io/gorm"
	"os"
)

func InitPosgres() (*gorm.DB, error) {
	env := datasources.Env{
		User:   "POSTGRES_USER",
		Pass:   "POSTGRES_PASSWORD",
		Schema: "POSTGRES_SCHEMA",
	}
	_ = os.Setenv("POSTGRES_USER", "test")
	_ = os.Setenv("POSTGRES_PASSWORD", "test")
	_ = os.Setenv("POSTGRES_SCHEMA", "testdb")

	c := env.Build(postgresql.New(postgresql.SSLMode(postgresql.SslDisable)))
	ds := datasources.NewDataSource(c)
	db := ds.GetConnection()
	if db.Error != nil {
		return nil, db.Error
	}
	return db, nil
}

func GetPostgresEntity(ctx context.Context, id string) (*ExampleTable, error) {
	db := gormsource.DB(ctx)
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
