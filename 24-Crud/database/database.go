package database

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func Connect() (*sql.DB, error) {
	stringConnection := "golang:Aspire@5740@/devbook?charset=utf8&parseTimeTrue&loc=Local"
	db, ex := sql.Open("mysql", stringConnection)
	if ex != nil {
		return nil, ex
	}

	if ex = db.Ping(); ex != nil {
		return nil, ex
	}

	return db, nil

}
