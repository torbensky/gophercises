package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var user, password, dbName string

func init() {
	dbName = os.Getenv("PG_DB_NAME")
	user = os.Getenv("PG_USER")
	password = os.Getenv("PG_PASSWORD")
}

func main() {
	connStr := fmt.Sprintf("postgres://%s:%s@localhost/%s?sslmode=disable", user, password, dbName)
	db, err := sql.Open("postgres", connStr)
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query("SELECT number FROM phone_numbers")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Phone numbers:")
	for rows.Next() {
		var number string
		if err := rows.Scan(&number); err != nil {
			log.Fatal(err)
		}
		fmt.Println(number)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}
