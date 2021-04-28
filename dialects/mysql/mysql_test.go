package mysql

import (
	"fmt"
	"gorm.io/driver/mysql"
	"os"
	"testing"
)

func TestNew(t *testing.T) {
	b := New(Protocol("tcp"), AllowAllFiles(true), AllowCleartextPasswords(true),
		AllowNativePasswords(true), Charset("utf8"), Collation("utf"), ClientFoundRows(true),
		ColumnsWithAlias(true), InterpolateParams(true), Loc("Local"), MaxAllowedPacket(1000),
		MultiStatements(true), ParseTime(true), ReadTimeout("100"), RejectReadOnly(true),
		ServerPubKey(".key"), Timeout("1000"), Tls("200"), WriteTimeout("10"))
	actual := b.Build("user", "pass", "host", 8088, "test")
	expected := "user:pass@tcp(host:8088)/test?allowAllFiles=true&allowCleartextPasswords=true&allowNativePasswords=true&charset=utf8&collation=utf&clientFoundRows=true&columnsWithAlias=true&interpolateParams=true&loc=Local&maxAllowedPacket=1000&multiStatements=true&parseTime=true&readTimeout=100&rejectReadOnly=true&serverPubKey=.key&timeout=1000&tls=200&writeTimeout=10"
	dialector := actual.(*mysql.Dialector)
	if expected != dialector.DSN {
		t.Errorf("expected=%s, actual=%s", expected, actual)
	} else {
		fmt.Printf("%s\n", actual)
	}
}

func TestEnv(t *testing.T) {
	env := Environment{
		InstanceName:            "MYSQL_INSTANCE_NAME",
		Protocol:                "MYSQL_PROTOCOL",
		AllowAllFiles:           "MYSQL_ALLOW_ALL_FILES",
		AllowCleartextPasswords: "MYSQL_ALLOW_CLEARTEXT_PASSWORDS",
		AllowNativePasswords:    "MYSQL_ALLOW_NATIVE_PASSWORDS",
		AllowOldPasswords:       "MYSQL_ALLOW_OLD_PASSWORD",
		Charset:                 "MYSQL_CHARSET",
		Collation:               "MYSQL_COLLATION",
		ClientFoundRows:         "MYSQL_CLIENT_FOUND_ROWS",
		ColumnsWithAlias:        "MYSQL_COLUMNS_WITH_ALIAS",
		InterpolateParams:       "MYSQL_INTERPOLATE_PARAMS",
		Loc:                     "MYSQL_LOC",
		MaxAllowedPacket:        "MYSQL_MAX_ALLOWED_PACKET",
		MultiStatements:         "MYSQL_MULTI_STATEMENTS",
		ParseTime:               "MYSQL_PARSE_TIME",
		ReadTimeout:             "MYSQL_READ_TIMEOUT",
		RejectReadOnly:          "MYSQL_REJECT_READ_ONLY",
		ServerPubKey:            "MYSQL_SERVER_PUB_KEY",
		Timeout:                 "MYSQL_TIMEOUT",
		Tls:                     "MYSQL_TLS",
		WriteTimeout:            "MYSQL_WRITE_TIMEOUT",
	}
	_ = os.Setenv("MYSQL_PROTOCOL", "udp")
	_ = os.Setenv("MYSQL_CHARSET", "utf8mb4")
	_ = os.Setenv("MYSQL_COLLATION", "utf8_unicode_ci")

	b := New(Env(&env))
	actual := b.Build("user", "pass", "host", 8088, "test")
	expected := "user:pass@udp(host:8088)/test?charset=utf8mb4&collation=utf8_unicode_ci"
	dialector := actual.(*mysql.Dialector)
	if expected != dialector.DSN {
		t.Errorf("expected=%s, actual=%s", expected, actual)
	} else {
		fmt.Printf("%s\n", actual)
	}
}
