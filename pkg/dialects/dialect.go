package dialects

import (
	"database/sql"
	"github.com/goccha/errors"
	"github.com/goccha/log"
	"gorm.io/gorm"
	"strings"
	"time"
)

type Builder interface {
	Name() string
	Build(user, password, host string, port int, dbname string) gorm.Dialector
	BuildString(user, password, host string, port int, dbname string) string
	BuildDialector(url string) gorm.Dialector
}

type Option func(b Builder)

func WriteString(buf *strings.Builder, key, value, sep string) {
	if len(sep) > 0 {
		buf.WriteString(sep)
	}
	buf.WriteString(key)
	buf.WriteString("=")
	buf.WriteString(value)
}

type Errors interface {
	IsNotAvailableLock(err error) bool
}

type Extension func(dialect, dsn string) (*sql.DB, error)

func Connect(dialect, dsn string, f Extension) (db *sql.DB, err error) {
	d := 1 * time.Second
	timeout := true
	count := 20
	for count > 0 {
		if db, err = f(dialect, dsn); err != nil {
			log.Info("%v", err)
			time.Sleep(d)
		} else {
			timeout = false
			log.Info("connection established.")
			break
		}
		count--
	}
	if timeout {
		return nil, errors.New("ping timeout")
	}
	return
}
