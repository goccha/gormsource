package dialects

import (
	"gorm.io/gorm"
	"strings"
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
