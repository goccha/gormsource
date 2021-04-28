module github.com/goccha/gormsource/dialects/postgresql

go 1.15

require (
	github.com/goccha/envar v0.1.0
	github.com/goccha/gormsource v1.3.0
	github.com/jackc/pgconn v1.8.0
	gorm.io/driver/postgres v1.0.8
	gorm.io/gorm v1.20.12
)

replace github.com/goccha/gormsource => ../..
