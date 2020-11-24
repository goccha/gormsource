package examples

import (
	"github.com/goccha/gormsource/pkg/datasources"
	"github.com/goccha/gormsource/pkg/datasources/postgresql"
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
