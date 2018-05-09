package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var user, password, dbName string

type phone struct {
	id     int
	number string
}

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

	rows, err := db.Query("SELECT id, number FROM phone_numbers")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Step 1.
	// First we iterate over the data and ensure the numbers all share the normalized form
	// this can create many duplicates, but we will deal with those in the second step
	for rows.Next() {
		var p phone
		if err := rows.Scan(&p.id, &p.number); err != nil {
			log.Fatal(err)
		}
		normNum := normalize(p.number)
		if normNum != p.number {
			err = updatePhoneNumber(db, p.id, normNum)
			if err != nil {
				fmt.Println("error updating to normalized form: ", p.id, normNum)
			}
		}
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	// Step 2.
	// Second, we execute a query that deletes all duplicates
	err = deleteDuplicates(db)
	if err != nil {
		fmt.Println(err)
	}
}

func updatePhoneNumber(db *sql.DB, id int, newNum string) error {
	stmt := `UPDATE phone_numbers SET number = $2 where id = $1`
	_, err := db.Exec(stmt, id, newNum)
	return err
}

func deleteDuplicates(db *sql.DB) error {
	stmt := `DELETE FROM phone_numbers a USING phone_numbers b where a.id > b.id AND a.number = b.number`
	_, err := db.Exec(stmt)
	return err
}
