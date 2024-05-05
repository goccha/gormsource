package sqlite3

import (
	"github.com/goccha/envar"
	"github.com/goccha/gormsource/pkg/dialects"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"strings"
)

func New(options ...dialects.Option) *Builder {
	b := &Builder{}
	for _, opt := range options {
		opt(b)
	}
	return b
}

type Builder struct {
	Path string
}

func (b *Builder) Name() string {
	return "sqlite3"
}

func (b *Builder) BuildDialector(url string) gorm.Dialector {
	return sqlite.Open(url)
}

func (b *Builder) BuildString(user, password, host string, port int, dbname string) string {
	buf := &strings.Builder{}
	buf.WriteString(b.Path)
	return buf.String()
}

func (b *Builder) Build(user, password, host string, port int, dbname string) gorm.Dialector {
	return sqlite.Open(b.BuildString(user, password, host, port, dbname))
}

type Environment struct {
	Path string
}

func (env Environment) Build(b *Builder) {
	Path(envar.String(env.Path, "SQLITE_PATH"))(b)
}

func Env(env Environment) dialects.Option {
	return func(b dialects.Builder) {
		env.Build(b.(*Builder))
	}
}
func Path(value string) dialects.Option {
	return func(b dialects.Builder) {
		if value != "" {
			b.(*Builder).Path = value
		}
	}
}
