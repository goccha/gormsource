package postgresql

import (
	"fmt"
	"gorm.io/driver/postgres"
	"os"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	b := New(SSLMode(SslRequire), FallbackApplicationName(""), ConnectTimeout(3*time.Minute),
		SSLCert("cert"), SSLKey(".key"), SSLRootCert("root"))
	actual := b.Build("user", "pass", "host", 8080, "test")
	expected := "user=user password=pass host=host port=8080 dbname=test sslmode=require connect_timeout=180 sslcert=cert sslkey=.key sslrootcert=root"
	dialector := actual.(*postgres.Dialector)
	if expected != dialector.DSN {
		t.Errorf("expected=%s, actual=%s", expected, actual)
	} else {
		fmt.Printf("%s\n", actual)
	}
}

func TestEnv(t *testing.T) {
	env := Environment{
		SslMode:                 "POSTGRES_SSL_MODE",
		FallbackApplicationName: "POSTGRES_FALLBACK_APPLICATION_NAME",
		ConnectTimeout:          "POSTGRES_CONNECT_TIMEOUT",
		SslCert:                 "POSTGRES_SSL_CERT",
		SslKey:                  "POSTGRES_SSL_KEY",
		SslRootCert:             "POSTGRES_SSL_ROOT_CERT",
	}
	_ = os.Setenv("POSTGRES_SSL_MODE", string(SslRequire))
	_ = os.Setenv("POSTGRES_CONNECT_TIMEOUT", "30s")

	b := New(Env(&env))
	actual := b.Build("user", "pass", "host", 8088, "test")
	expected := "user=user password=pass host=host port=8088 dbname=test sslmode=require connect_timeout=30"
	dialector := actual.(*postgres.Dialector)
	if expected != dialector.DSN {
		t.Errorf("expected=%s, actual=%s", expected, actual)
	} else {
		fmt.Printf("%s\n", actual)
	}
}
