package main

import (
	"github.com/goccha/gormsource/pkg/datasources"
	"github.com/goccha/gormsource/pkg/datasources/mysql"
	"github.com/rs/zerolog/log"
	"os"
	"time"
)

func main() {
	c := datasources.Config{
		User:   "user",
		Pass:   "password",
		Host:   "localhost",
		Schema: "testdb",
		PoolConfig: datasources.PoolConfig{
			MaxIdleConns:    10,
			MaxOpenConns:    50,
			ConnMaxLifetime: 5 * time.Minute,
		},
	}
	env := mysql.Environment{
		Charset:              "MYSQL_CHARSET",
		ParseTime:            "MYSQL_PARSE_TIME",
		Loc:                  "MYSQL_LOC",
		AllowNativePasswords: "MYSQL_ALLOW_NATIVE_PASSWORDS",
	}
	_ = os.Setenv("MYSQL_CHARSET", "utf8mb4")
	_ = os.Setenv("MYSQL_PARSE_TIME", "true")
	_ = os.Setenv("MYSQL_LOC", "Local")
	_ = os.Setenv("MYSQL_ALLOW_NATIVE_PASSWORDS", "true")

	ds := datasources.NewDataSource(c.Dialect(mysql.New(mysql.Env(&env))))
	db := ds.GetConnection()
	if db.Error != nil {
		log.Err(db.Error)
	}
}
