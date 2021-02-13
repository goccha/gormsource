package datasources

import (
	"github.com/goccha/envar"
	"github.com/goccha/gormsource/pkg/dialects"
	"time"
)

type Env struct {
	User             string
	Pass             string
	Host             string
	Port             string
	Schema           string
	ConnectionString string
	MaxIdleConns     string
	MaxOpenConns     string
	ConnMaxLifetime  string
	Debug            string
}

func (e *Env) Build(builder dialects.Builder) *Config {
	config := &Config{}
	config.ConnectionString = envar.String(e.ConnectionString, "DB_CONNECT_URL")
	if len(config.ConnectionString) == 0 {
		config.User = envar.String(e.User, "DB_USER")
		config.Pass = envar.String(e.Pass, "DB_PASSWORD")
		config.Host = envar.Get(e.Host, "DB_HOST").String("127.0.0.1")
		config.Port = envar.Get(e.Port, "DB_PORT").Int(0)
		config.Schema = envar.String(e.Schema, "DB_SCHEMA")
	}
	config.Dialect(builder)
	config.MaxIdleConns = envar.Get(e.MaxIdleConns, "DB_MAX_IDLE_CONNECTIONS").Int(10)
	config.MaxOpenConns = envar.Get(e.MaxOpenConns, "DB_MAX_OPEN_CONNECTIONS").Int(50)
	config.ConnMaxLifetime = envar.Get(e.ConnMaxLifetime, "DB_CONNECTION_MAX_LIFETIME").Duration(time.Hour)
	config.Debug = envar.Get(e.Debug, "GORM_LOG_MODE").Bool(false)
	return config
}
