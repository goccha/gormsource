package postgresql

import (
	"errors"
	"github.com/goccha/envar"
	"github.com/goccha/gormsource/pkg/dialects"
	"github.com/jackc/pgconn"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"strconv"
	"strings"
	"time"
)

const (
	DefaultPort      = "5432"
	NotAvailableLock = "55P03"
)

func New(options ...dialects.Option) *Builder {
	b := &Builder{}
	for _, opt := range options {
		opt(b)
	}
	return b
}

type Builder struct {
	SslMode                 string
	FallbackApplicationName string
	ConnectTimeout          time.Duration
	SslCert                 string
	SslKey                  string
	SslRootCert             string
	PreferSimpleProtocol    bool
	Extension               dialects.Extension
}

func (b *Builder) Name() string {
	return "pgx"
}

func (b *Builder) BuildDialector(url string) gorm.Dialector {
	return postgres.Open(url)
}

func (b *Builder) BuildString(user, password, host string, port int, dbname string) string {
	buf := &strings.Builder{}
	dialects.WriteString(buf, "user", user, "")
	dialects.WriteString(buf, "password", password, " ")
	if host != "" {
		dialects.WriteString(buf, "host", host, " ")
	}
	if port > 0 {
		dialects.WriteString(buf, "port", strconv.Itoa(port), " ")
	} else {
		dialects.WriteString(buf, "port", DefaultPort, " ")
	}
	dialects.WriteString(buf, "dbname", dbname, " ")
	if b.SslMode != "" {
		dialects.WriteString(buf, "sslmode", b.SslMode, " ")
	}
	if b.FallbackApplicationName != "" {
		dialects.WriteString(buf, "fallback_application_name", b.FallbackApplicationName, " ")
	}
	if b.ConnectTimeout > 0 {
		sec := strconv.FormatFloat(b.ConnectTimeout.Seconds(), 'f', 0, 64)
		dialects.WriteString(buf, "connect_timeout", sec, " ")
	}
	if b.SslCert != "" {
		dialects.WriteString(buf, "sslcert", b.SslCert, " ")
	}
	if b.SslKey != "" {
		dialects.WriteString(buf, "sslkey", b.SslKey, " ")
	}
	if b.SslRootCert != "" {
		dialects.WriteString(buf, "sslrootcert", b.SslRootCert, " ")
	}
	return buf.String()
}

func (b *Builder) Build(user, password, host string, port int, dbname string) gorm.Dialector {
	dsn := b.BuildString(user, password, host, port, dbname)
	if b.Extension != nil {
		if db, err := dialects.Connect(b.Name(), dsn, b.Extension); err != nil {
			return nil
		} else {
			return postgres.New(postgres.Config{
				DSN:                  dsn,
				Conn:                 db,
				PreferSimpleProtocol: b.PreferSimpleProtocol,
			})
		}
	}
	return postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: b.PreferSimpleProtocol,
	})
}

func (b *Builder) IsNotAvailableLock(err error) bool {
	var v *pgconn.PgError
	if errors.As(err, &v) {
		return v.Code == NotAvailableLock
	}
	return false
}

type SSLOption string

const (
	// disable - No SSL
	SslDisable SSLOption = "disable"
	// require - Always SSL (skip verification)
	SslRequire SSLOption = "require"
	// verify-ca - Always SSL (verify that the certificate presented by the
	// server was signed by a trusted CA)
	SslVerifyCa SSLOption = "verify-ca"
	// verify-full - Always SSL (verify that the certification presented by the
	// server was signed by a trusted CA and the server host name matches the one in the certificate)
	SslVerifyFull SSLOption = "verify-full"
)

type Environment struct {
	SslMode                 string
	FallbackApplicationName string
	ConnectTimeout          string
	SslCert                 string
	SslKey                  string
	SslRootCert             string
	PreferSimpleProtocol    string
}

func (env *Environment) Build(b *Builder) {
	if ev := envar.Get(env.SslMode); ev.Has() {
		SSLMode(SSLOption(ev.String("disable")))(b)
	}
	FallbackApplicationName(envar.String(env.FallbackApplicationName))(b)
	d := envar.Duration(env.ConnectTimeout)
	if d > 0 {
		ConnectTimeout(d)(b)
	}
	SSLCert(envar.String(env.SslCert))(b)
	SSLKey(envar.String(env.SslKey))(b)
	SSLRootCert(envar.String(env.SslRootCert))(b)
	if envar.Has(env.PreferSimpleProtocol) {
		v := envar.Bool(env.PreferSimpleProtocol)
		PreferSimpleProtocol(&v)(b)
	}
}

func Env(env *Environment) dialects.Option {
	return func(b dialects.Builder) {
		env.Build(b.(*Builder))
	}
}

func SSLMode(value SSLOption) dialects.Option {
	return func(b dialects.Builder) {
		b.(*Builder).SslMode = string(value)
	}
}

func FallbackApplicationName(name string) dialects.Option {
	return func(b dialects.Builder) {
		if name != "" {
			b.(*Builder).FallbackApplicationName = name
		}
	}
}

func ConnectTimeout(t time.Duration) dialects.Option {
	return func(b dialects.Builder) {
		b.(*Builder).ConnectTimeout = t
	}
}

func SSLCert(location string) dialects.Option {
	return func(b dialects.Builder) {
		if location != "" {
			b.(*Builder).SslCert = location
		}
	}
}

func SSLKey(location string) dialects.Option {
	return func(b dialects.Builder) {
		if location != "" {
			b.(*Builder).SslKey = location
		}
	}
}

func SSLRootCert(location string) dialects.Option {
	return func(b dialects.Builder) {
		if location != "" {
			b.(*Builder).SslRootCert = location
		}
	}
}

func PreferSimpleProtocol(simple *bool) dialects.Option {
	return func(b dialects.Builder) {
		if simple != nil {
			b.(*Builder).PreferSimpleProtocol = *simple
		}
	}
}
func Extension(f dialects.Extension) dialects.Option {
	return func(b dialects.Builder) {
		if f != nil {
			b.(*Builder).Extension = f
		}
	}
}
