module github.com/goccha/gormsource/dialects/mysql

go 1.15

require (
	github.com/go-sql-driver/mysql v1.5.0
	github.com/goccha/envar v0.1.0
	github.com/goccha/gormsource v1.2.0
	gorm.io/driver/mysql v1.0.5
	gorm.io/gorm v1.21.3
)

replace github.com/goccha/gormsource => ../..
