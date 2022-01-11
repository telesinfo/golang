package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	stringConnection := "golang:Aspire@5740@/devbook?charset=utf8&parseTimeTrue&loc=Local"
	db, error := sql.Open("mysql", stringConnection)
	if error != nil {
		log.Fatal(error)
	}
	defer db.Close()

	if error = db.Ping(); error != nil {
		log.Fatal(error)
	}

	fmt.Println("Successful connection!")

	rows, error := db.Query("select * from users")
	if error != nil {
		log.Fatal(error)
	}
	defer rows.Close()

	fmt.Println(rows)

}
