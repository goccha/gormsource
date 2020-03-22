package sqlite3

import (
	"fmt"
	"os"
	"testing"
)

func TestNew(t *testing.T) {
	b := New(Path("./test.db"))
	actual := b.Build("", "", "", 0, "")
	expected := "./test.db"
	if expected != actual {
		t.Errorf("expected=%s, actual=%s", expected, actual)
	} else {
		fmt.Printf("%s\n", actual)
	}
}

func TestEnv(t *testing.T) {
	env := Environment{Path: "SQLITE3_PATH"}
	_ = os.Setenv("SQLITE3_PATH", "./testenv.db")

	b := New(Env(&env))
	actual := b.Build("", "", "", 0, "")
	expected := "./testenv.db"
	if expected != actual {
		t.Errorf("expected=%s, actual=%s", expected, actual)
	} else {
		fmt.Printf("%s\n", actual)
	}
}
