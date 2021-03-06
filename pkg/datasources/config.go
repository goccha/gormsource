package datasources

import (
	"github.com/goccha/gormsource/pkg/dialects"
	"gorm.io/gorm"
	"time"
)

type Config struct {
	User             string
	Pass             string
	Host             string
	Port             int
	Schema           string
	dialect          dialects.Builder
	ConnectionString string
	PoolConfig
	Debug bool
	gorm.Config
}

func (c *Config) Dialect(builder dialects.Builder) *Config {
	c.dialect = builder
	return c
}

type PoolConfig struct {
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
}

func (c *Config) String() string {
	if len(c.ConnectionString) > 0 {
		return c.ConnectionString
	}
	return c.dialect.BuildString(c.User, c.Pass, c.Host, c.Port, c.Schema)
}

func (c *Config) Build() gorm.Dialector {
	if len(c.ConnectionString) > 0 {
		return c.dialect.BuildDialector(c.ConnectionString)
	}
	return c.dialect.Build(c.User, c.Pass, c.Host, c.Port, c.Schema)
}
