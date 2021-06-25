module github.com/goccha/gormsource/dialects/mysql

go 1.15

require (
	github.com/go-sql-driver/mysql v1.6.0
	github.com/goccha/envar v0.1.1
	github.com/goccha/gormsource v1.3.1
	gorm.io/driver/mysql v1.1.1
	gorm.io/gorm v1.21.11
)

replace github.com/goccha/gormsource => ../..
