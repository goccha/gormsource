package mysql

import (
	"github.com/goccha/envar"
	"github.com/goccha/gormsource/pkg/datasources/dialects"
	"strconv"
	"strings"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func New(options ...dialects.Option) *Builder {
	b := &Builder{}
	for _, opt := range options {
		opt(b)
	}
	return b
}

type Builder struct {
	InstanceName            string
	Protocol                string
	AllowAllFiles           bool
	AllowCleartextPasswords bool
	AllowNativePasswords    *bool
	AllowOldPasswords       bool
	Charset                 string
	Collation               string
	ClientFoundRows         bool
	ColumnsWithAlias        bool
	InterpolateParams       bool
	Loc                     string
	MaxAllowedPacket        int
	MultiStatements         bool
	ParseTime               bool
	ReadTimeout             string
	RejectReadOnly          bool
	ServerPubKey            string
	Timeout                 string
	Tls                     string
	WriteTimeout            string
	SystemVariables         map[string]string
}

func (b *Builder) Name() string {
	return "mysql"
}
func (b *Builder) Put(k string, v string) *Builder {
	b.SystemVariables[k] = v
	return b
}
func and(builder *strings.Builder) *strings.Builder {
	if builder.Len() > 0 {
		builder.WriteString("&")
	}
	return builder
}
func (b *Builder) Build(user, password, host string, port int, dbname string) string {
	buf := &strings.Builder{}
	buf.WriteString(user)
	buf.WriteString(":")
	buf.WriteString(password)
	buf.WriteString("@")
	if b.Protocol != "" {
		buf.WriteString(b.Protocol)
	} else {
		buf.WriteString("tcp")
	}
	buf.WriteString("(")

	if len(b.InstanceName) > 0 {
		buf.WriteString(b.InstanceName)
	} else {
		if host != "" {
			buf.WriteString(host)
		} else {
			buf.WriteString("127.0.0.1")
		}
		buf.WriteString(":")
		if port > 0 {
			buf.WriteString(strconv.Itoa(port))
		} else {
			buf.WriteString("3306")
		}
	}
	buf.WriteString(")/")
	buf.WriteString(dbname)

	options := &strings.Builder{}
	if b.AllowAllFiles {
		options.WriteString("allowAllFiles=true")
	}
	if b.AllowCleartextPasswords {
		and(options).WriteString("allowCleartextPasswords=true")
	}
	if b.AllowNativePasswords != nil {
		dialects.WriteString(options, "allowNativePasswords", strconv.FormatBool(*b.AllowNativePasswords), "&")
	}
	if b.AllowOldPasswords {
		and(options).WriteString("allowOldPasswords=true")
	}
	if len(b.Charset) > 0 {
		and(options).WriteString("charset=")
		options.WriteString(b.Charset)
	}
	if len(b.Collation) > 0 {
		dialects.WriteString(options, "collation", b.Collation, "&")
	}
	if b.ClientFoundRows {
		and(options).WriteString("clientFoundRows=true")
	}
	if b.ColumnsWithAlias {
		and(options).WriteString("columnsWithAlias=true")
	}
	if b.InterpolateParams {
		and(options).WriteString("interpolateParams=true")
	}
	if len(b.Loc) > 0 {
		dialects.WriteString(options, "loc", b.Loc, "&")
	}
	if b.MaxAllowedPacket > 0 {
		and(options).WriteString("maxAllowedPacket=")
		options.WriteString(strconv.Itoa(b.MaxAllowedPacket))
	}
	if b.MultiStatements {
		and(options).WriteString("multiStatements=true")
	}
	if b.ParseTime {
		and(options).WriteString("parseTime=true")
	}
	if len(b.ReadTimeout) > 0 {
		and(options).WriteString("readTimeout=")
		options.WriteString(b.ReadTimeout)
	}
	if b.RejectReadOnly {
		and(options).WriteString("rejectReadOnly=true")
	}
	if len(b.ServerPubKey) > 0 {
		and(options).WriteString("serverPubKey=")
		options.WriteString(b.ServerPubKey)
	}
	if len(b.Timeout) > 0 {
		and(options).WriteString("timeout=")
		options.WriteString(b.Timeout)
	}
	if len(b.Tls) > 0 {
		and(options).WriteString("tls=")
		options.WriteString(b.Tls)
	}
	if len(b.WriteTimeout) > 0 {
		and(options).WriteString("writeTimeout=")
		options.WriteString(b.WriteTimeout)
	}
	if len(b.SystemVariables) > 0 {
		for k := range b.SystemVariables {
			v := b.SystemVariables[k]
			and(options).WriteString(k)
			options.WriteString("=")
			options.WriteString(v)
		}
	}
	if options.Len() > 0 {
		buf.WriteString("?")
		buf.WriteString(options.String())
	}
	return buf.String()
}

//func DefaultOptions() *Option {
//	return &Option{Charset: "utf8mb4", ParseTime: true, Loc: "Local", AllowNativePasswords: true}
//}

type Environment struct {
	InstanceName            string
	Protocol                string
	AllowAllFiles           string
	AllowCleartextPasswords string
	AllowNativePasswords    string
	AllowOldPasswords       string
	Charset                 string
	Collation               string
	ClientFoundRows         string
	ColumnsWithAlias        string
	InterpolateParams       string
	Loc                     string
	MaxAllowedPacket        string
	MultiStatements         string
	ParseTime               string
	ReadTimeout             string
	RejectReadOnly          string
	ServerPubKey            string
	Timeout                 string
	Tls                     string
	WriteTimeout            string
}

func (env *Environment) Build(b *Builder) {
	InstanceName(envar.String(env.InstanceName))(b)
	Protocol(envar.String(env.Protocol))(b)
	AllowAllFiles(envar.Bool(env.AllowAllFiles))(b)
	AllowCleartextPasswords(envar.Bool(env.AllowCleartextPasswords))(b)
	if envar.Has(env.AllowNativePasswords) {
		AllowNativePasswords(envar.Bool(env.AllowNativePasswords))(b)
	}
	AllowOldPasswords(envar.Bool(env.AllowOldPasswords))(b)
	Charset(envar.String(env.Charset))(b)
	Collation(envar.String(env.Collation))(b)
	ClientFoundRows(envar.Bool(env.ClientFoundRows))(b)
	ColumnsWithAlias(envar.Bool(env.ColumnsWithAlias))(b)
	InterpolateParams(envar.Bool(env.InterpolateParams))(b)
	Loc(envar.String(env.Loc))(b)
	MaxAllowedPacket(envar.Int(env.MaxAllowedPacket))(b)
	MultiStatements(envar.Bool(env.MultiStatements))(b)
	ParseTime(envar.Bool(env.ParseTime))(b)
	ReadTimeout(envar.String(env.ReadTimeout))(b)
	RejectReadOnly(envar.Bool(env.RejectReadOnly))(b)
	ServerPubKey(envar.String(env.ServerPubKey))(b)
	Timeout(envar.String(env.Timeout))(b)
	Tls(envar.String(env.Tls))(b)
	WriteTimeout(envar.String(env.WriteTimeout))(b)
}

func Env(env *Environment) dialects.Option {
	return func(b dialects.Builder) {
		env.Build(b.(*Builder))
	}
}

func InstanceName(value string) dialects.Option {
	return func(b dialects.Builder) {
		b.(*Builder).InstanceName = value
	}
}

func Protocol(value string) dialects.Option {
	return func(b dialects.Builder) {
		b.(*Builder).Protocol = value
	}
}
func AllowAllFiles(value bool) dialects.Option {
	return func(b dialects.Builder) {
		b.(*Builder).AllowAllFiles = value
	}
}
func AllowCleartextPasswords(value bool) dialects.Option {
	return func(b dialects.Builder) {
		b.(*Builder).AllowCleartextPasswords = value
	}
}
func AllowNativePasswords(value bool) dialects.Option {
	return func(b dialects.Builder) {
		b.(*Builder).AllowNativePasswords = &value
	}
}
func AllowOldPasswords(value bool) dialects.Option {
	return func(b dialects.Builder) {
		b.(*Builder).AllowOldPasswords = value
	}
}
func Charset(value string) dialects.Option {
	return func(b dialects.Builder) {
		b.(*Builder).Charset = value
	}
}
func Collation(value string) dialects.Option {
	return func(b dialects.Builder) {
		b.(*Builder).Collation = value
	}
}
func ClientFoundRows(value bool) dialects.Option {
	return func(b dialects.Builder) {
		b.(*Builder).ClientFoundRows = value
	}
}
func ColumnsWithAlias(value bool) dialects.Option {
	return func(b dialects.Builder) {
		b.(*Builder).ColumnsWithAlias = value
	}
}
func InterpolateParams(value bool) dialects.Option {
	return func(b dialects.Builder) {
		b.(*Builder).InterpolateParams = value
	}
}
func Loc(value string) dialects.Option {
	return func(b dialects.Builder) {
		b.(*Builder).Loc = value
	}
}
func MaxAllowedPacket(value int) dialects.Option {
	return func(b dialects.Builder) {
		b.(*Builder).MaxAllowedPacket = value
	}
}
func MultiStatements(value bool) dialects.Option {
	return func(b dialects.Builder) {
		b.(*Builder).MultiStatements = value
	}
}
func ParseTime(value bool) dialects.Option {
	return func(b dialects.Builder) {
		b.(*Builder).ParseTime = value
	}
}
func ReadTimeout(value string) dialects.Option {
	return func(b dialects.Builder) {
		b.(*Builder).ReadTimeout = value
	}
}
func RejectReadOnly(value bool) dialects.Option {
	return func(b dialects.Builder) {
		b.(*Builder).RejectReadOnly = value
	}
}
func ServerPubKey(value string) dialects.Option {
	return func(b dialects.Builder) {
		b.(*Builder).ServerPubKey = value
	}
}
func Timeout(value string) dialects.Option {
	return func(b dialects.Builder) {
		b.(*Builder).Timeout = value
	}
}
func Tls(value string) dialects.Option {
	return func(b dialects.Builder) {
		b.(*Builder).Tls = value
	}
}
func WriteTimeout(value string) dialects.Option {
	return func(b dialects.Builder) {
		b.(*Builder).WriteTimeout = value
	}
}
