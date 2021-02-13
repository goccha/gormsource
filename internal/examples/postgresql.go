package examples

import (
	"github.com/goccha/gormsource/pkg/datasources"
	"github.com/goccha/gormsource/pkg/dialects/postgresql"
	"gorm.io/gorm"
	"os"
)

func InitPrimaryPosgres() (*gorm.DB, error) {
	env := datasources.Env{
		User:   "POSTGRES_USER",
		Pass:   "POSTGRES_PASSWORD",
		Schema: "POSTGRES_SCHEMA",
		Debug:  "POSTGRES_DEBUG",
	}
	_ = os.Setenv("POSTGRES_USER", "test")
	_ = os.Setenv("POSTGRES_PASSWORD", "test")
	_ = os.Setenv("POSTGRES_SCHEMA", "testdb")
	_ = os.Setenv("POSTGRES_DEBUG", "true")

	c := env.Build(postgresql.New(postgresql.SSLMode(postgresql.SslDisable)))
	ds := datasources.NewDataSource(c)
	db := ds.GetConnection()
	if db.Error != nil {
		return nil, db.Error
	}
	return db, nil
}

func InitReplicaPosgres() (*gorm.DB, error) {
	env := datasources.Env{
		User:   "POSTGRES_USER",
		Pass:   "POSTGRES_PASSWORD",
		Port:   "POSTGRES_PORT",
		Schema: "POSTGRES_SCHEMA",
		Debug:  "POSTGRES_DEBUG",
	}
	_ = os.Setenv("POSTGRES_USER", "test")
	_ = os.Setenv("POSTGRES_PASSWORD", "test")
	_ = os.Setenv("POSTGRES_SCHEMA", "testdb")
	_ = os.Setenv("POSTGRES_PORT", "5532")
	_ = os.Setenv("POSTGRES_DEBUG", "true")

	c := env.Build(postgresql.New(postgresql.SSLMode(postgresql.SslDisable)))
	ds := datasources.NewDataSource(c)
	db := ds.GetConnection()
	if db.Error != nil {
		return nil, db.Error
	}
	return db, nil
}
