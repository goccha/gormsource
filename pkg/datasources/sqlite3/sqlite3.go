package sqlite3

import (
	"github.com/goccha/envar"
	"github.com/goccha/gormsource/pkg/datasources/dialects"
	"strings"
	//_ "github.com/jinzhu/gorm/dialects/sqlite"
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
func (b *Builder) Build(user, password, host string, port int, dbname string) string {
	buf := &strings.Builder{}
	buf.WriteString(b.Path)
	return buf.String()
}

type Environment struct {
	Path string
}

func (env *Environment) Build(b *Builder) {
	Path(envar.String(env.Path))(b)
}

func Env(env *Environment) dialects.Option {
	return func(b dialects.Builder) {
		env.Build(b.(*Builder))
	}
}
func Path(value string) dialects.Option {
	return func(b dialects.Builder) {
		b.(*Builder).Path = value
	}
}
