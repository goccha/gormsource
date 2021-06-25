module github.com/goccha/gormsource/dialects/sqlite3

go 1.15

require (
	github.com/goccha/envar v0.1.1
	github.com/goccha/gormsource v1.3.1
	gorm.io/driver/sqlite v1.1.4
	gorm.io/gorm v1.21.11
)

replace github.com/goccha/gormsource => ../..
