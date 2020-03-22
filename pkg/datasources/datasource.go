package datasources

import (
	"github.com/goccha/envar"
	"github.com/goccha/log"
	"github.com/pkg/errors"
	"time"

	"github.com/jinzhu/gorm"
)

func connect(dialect string, options interface{}) (*gorm.DB, error) {
	d := 1 * time.Second
	timeout := true
	var conn *gorm.DB
	count := 20
	for count > 0 {
		c, err := gorm.Open(dialect, options)
		if err != nil {
			log.Info("%v", err)
			time.Sleep(d)
		} else {
			timeout = false
			conn = c
			log.Info("connection established.")
			break
		}
		count--
	}
	if timeout {
		return nil, errors.New("ping timeout")
	}
	return conn, nil
}

func newDB(config *Config) *gorm.DB {
	dialect := config.dialect.Name()
	log.Info("newConnection(" + dialect + ")")
	log.Debug(config.String())
	conn, err := connect(dialect, config.String())
	if err != nil {
		panic(err)
	}
	conn.LogMode(envar.Get("GORM_LOG_MODE").Bool(false))
	conn.DB().SetMaxIdleConns(config.MaxIdleConns)
	conn.DB().SetMaxOpenConns(config.MaxOpenConns)
	conn.DB().SetConnMaxLifetime(config.ConnMaxLifetime)
	return conn
}

type DataSource struct {
	db *gorm.DB
}

func (ds *DataSource) GetConnection() *gorm.DB {
	return ds.db
}
func NewDataSource(c *Config) *DataSource {
	if c == nil {
		c = &Config{}
	}
	return &DataSource{newDB(c)}
}
