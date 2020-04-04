package datasources

import (
	"github.com/goccha/gormsource/pkg/datasources/dialects"
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
	return c.dialect.Build(c.User, c.Pass, c.Host, c.Port, c.Schema)
}
