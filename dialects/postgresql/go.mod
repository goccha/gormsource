module github.com/goccha/gormsource/dialects/postgresql

go 1.15

require (
	github.com/goccha/envar v0.1.1
	github.com/goccha/gormsource v1.3.1
	github.com/jackc/pgconn v1.8.1
	gorm.io/driver/postgres v1.1.0
	gorm.io/gorm v1.21.11
)

replace github.com/goccha/gormsource => ../..
