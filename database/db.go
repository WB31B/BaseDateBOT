package database

import (
	"TGbot/config"
	"database/sql"
	"fmt"
)

const (
	host     = "localhost"
	port     = "5432"
	user     = "postgres"
	password = config.DBKEY
	dbName   = "tusergbot"
)

func Connect() (*sql.DB, error) {
	var psqlsconn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbName)

	db, err := sql.Open("postgres", psqlsconn)
	if err != nil {
		return nil, err
	}

	return db, nil
}
