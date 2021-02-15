package repo

import (
	"database/sql"
	"fmt"
	"os"
)

type Db struct{}

func (d *Db) getConnection() *sql.DB {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")
	name := os.Getenv("DB_NAME")

	cs := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, pass, name)

	// open database
	db, err := sql.Open("postgres", cs)
	if err != nil {
		panic(err)
	}

	return db
}
