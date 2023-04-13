package database

import (
	"database/sql"
	"fmt"
)

const (
	host     = "localhost"
	port     = "5432"
	user     = "postgres"
	password = "pG2r4hack"
	dbName   = "tusergbot"
)

var psqlsconn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
	host, port, user, password, dbName)

func Connect() {
	db, err := sql.Open("postgres", psqlsconn)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println(db)
}
