package sqlite3

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"os"
	"testing"
)

func TestNew(t *testing.T) {
	b := New(Path("./test.db"))
	actual := b.Build("", "", "", 0, "")
	expected := "./test.db"
	dialector := actual.(*sqlite.Dialector)
	if expected != dialector.DSN {
		t.Errorf("expected=%s, actual=%s", expected, actual)
	} else {
		fmt.Printf("%s\n", actual)
	}
}

func TestEnv(t *testing.T) {
	_ = os.Setenv("SQLITE3_PATH", "./testenv.db")
	b := New(Env(Environment{Path: "SQLITE3_PATH"}))
	actual := b.Build("", "", "", 0, "")
	expected := "./testenv.db"
	dialector := actual.(*sqlite.Dialector)
	if expected != dialector.DSN {
		t.Errorf("expected=%s, actual=%s", expected, actual)
	} else {
		fmt.Printf("%s\n", actual)
	}
}
