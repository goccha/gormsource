package datasources

import (
	"time"

	"github.com/goccha/log"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func connect(dialect gorm.Dialector, config *gorm.Config) (*gorm.DB, error) {
	d := 1 * time.Second
	timeout := true
	var conn *gorm.DB
	count := 20
	for count > 0 {
		c, err := gorm.Open(dialect, config)
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
	db, err := connect(config.Build(), &config.Config)
	if err != nil {
		panic(err)
	}
	sqlDb, err := db.DB()
	if err != nil {
		panic(err)
	}
	if config.Debug {
		db.Logger = db.Logger.LogMode(logger.Info)
	}
	sqlDb.SetMaxIdleConns(config.MaxIdleConns)
	sqlDb.SetMaxOpenConns(config.MaxOpenConns)
	sqlDb.SetConnMaxLifetime(config.ConnMaxLifetime)
	return db
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

func Close(db *gorm.DB) {
	if db != nil {
		if sqlDB, err := db.DB(); err != nil {
			log.Warn("%v", err)
		} else {
			if err := sqlDB.Close(); err != nil {
				log.Warn("%v", err)
			}
		}
	}
}
