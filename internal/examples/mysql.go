package examples

import (
	"github.com/goccha/gormsource/pkg/datasources"
	"github.com/goccha/gormsource/pkg/dialects/mysql"
	"gorm.io/gorm"
	"os"
	"time"
)

func InitPrimaryMysql() (*gorm.DB, error) {
	c := datasources.Config{
		User:   "test",
		Pass:   "test",
		Schema: "testdb",
		PoolConfig: datasources.PoolConfig{
			MaxIdleConns:    10,
			MaxOpenConns:    50,
			ConnMaxLifetime: 5 * time.Minute,
		},
		Debug: true,
	}
	env := mysql.Environment{
		ParseTime:            "MYSQL_PARSE_TIME",
		Loc:                  "MYSQL_LOC",
		AllowNativePasswords: "MYSQL_ALLOW_NATIVE_PASSWORDS",
	}
	_ = os.Setenv("MYSQL_PARSE_TIME", "true")
	_ = os.Setenv("MYSQL_LOC", "Local")
	_ = os.Setenv("MYSQL_ALLOW_NATIVE_PASSWORDS", "true")

	//	ds := datasources.NewDataSource(c.Dialect(mysql.New(mysql.Env(&env), mysql.Charset("utf8mb4"), mysql.Collation("utf8mb4_bin"))))
	ds := datasources.NewDataSource(c.Dialect(mysql.New(mysql.Charset("utf8mb4"), mysql.Collation("utf8mb4_bin"), mysql.Env(&env))))
	db := ds.GetConnection()
	if db.Error != nil {
		return nil, db.Error
	}
	return db, nil
}

func InitReplicaMysql() (*gorm.DB, error) {
	c := datasources.Config{
		User:   "test",
		Pass:   "test",
		Schema: "testdb",
		Port:   3406,
		PoolConfig: datasources.PoolConfig{
			MaxIdleConns:    10,
			MaxOpenConns:    50,
			ConnMaxLifetime: 5 * time.Minute,
		},
		Debug: true,
	}
	env := mysql.Environment{
		ParseTime:            "MYSQL_PARSE_TIME",
		Loc:                  "MYSQL_LOC",
		AllowNativePasswords: "MYSQL_ALLOW_NATIVE_PASSWORDS",
	}
	_ = os.Setenv("MYSQL_PARSE_TIME", "true")
	_ = os.Setenv("MYSQL_LOC", "Local")
	_ = os.Setenv("MYSQL_ALLOW_NATIVE_PASSWORDS", "true")

	ds := datasources.NewDataSource(c.Dialect(mysql.New(mysql.Env(&env), mysql.Charset("utf8mb4"), mysql.Collation("utf8mb4_bin"))))
	db := ds.GetConnection()
	if db.Error != nil {
		return nil, db.Error
	}
	return db, nil
}
