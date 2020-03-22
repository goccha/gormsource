package datasources

import (
	"github.com/goccha/envar"
	"github.com/goccha/gormsource/pkg/datasources/dialects"
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
}

func (e *Env) getConnectionString() string {
	return envar.String(e.ConnectionString, "DB_CONNECT_URL")
}
func (e *Env) getMaxIdleConns() int {
	return envar.Get(e.MaxIdleConns, "DB_MAX_IDLE_CONNECTIONS").Int(10)
}
func (e *Env) getMaxOpenConns() int {
	return envar.Get(e.MaxOpenConns, "DB_MAX_OPEN_CONNECTIONS").Int(50)
}
func (e *Env) getConnMaxLifetime() time.Duration {
	return envar.Get(e.ConnMaxLifetime, "DB_CONNECTION_MAX_LIFETIME").Duration(time.Hour)
}

func (e *Env) Build(builder dialects.Builder) *Config {
	config := &Config{}
	config.ConnectionString = e.getConnectionString()
	if len(config.ConnectionString) == 0 {
		config.User = envar.String(e.User, "DB_USER")
		config.Pass = envar.String(e.Pass, "DB_PASSWORD")
		config.Host = envar.Get(e.Host, "DB_HOST").String("127.0.0.1")
		config.Port = envar.Get(e.Port, "DB_PORT").Int(3306)
		config.Schema = envar.String(e.Schema, "DB_SCHEMA")
	}
	config.Dialect(builder)
	config.MaxIdleConns = e.getMaxIdleConns()
	config.MaxOpenConns = e.getMaxOpenConns()
	config.ConnMaxLifetime = e.getConnMaxLifetime()
	return config
}
