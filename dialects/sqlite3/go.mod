module github.com/goccha/gormsource/dialects/sqlite3

go 1.15

require (
	github.com/goccha/envar v0.1.0
	github.com/goccha/gormsource v1.2.0
	gorm.io/driver/sqlite v1.1.4
	gorm.io/gorm v1.20.8
)

replace github.com/goccha/gormsource => ../..
