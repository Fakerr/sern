package persist

import (
	"database/sql"
	"log"

	"github.com/Fakerr/sern/config"

	_ "github.com/lib/pq"
)

// Postgres connection
var Conn *sql.DB

func InitConnection() {
	Conn = newConnection()
}

func newConnection() *sql.DB {
	var err error
	conn, err := sql.Open("postgres", config.DatabaseURL)
	if err != nil {
		panic(err)
	}

	if err = conn.Ping(); err != nil {
		panic(err)
	}

	log.Println("INFO: You connected to your database.")
	return conn
}

// Create tables if doesn't exist
func InitTables() {
	_, err := Conn.Exec("CREATE TABLE IF NOT EXISTS repositories (" +
		"ID  SERIAL PRIMARY KEY, " +
		"INSTALLATION            INT        NOT NULL, " +
		"FULLNAME                TEXT       NOT NULL, " +
		"OWNER                   TEXT       NOT NULL, " +
		"ORG                     TEXT       NOT NULL, " +
		"PRIVATE                 BOOL       NOT NULL, " +
		"CREATED_AT              TIMESTAMP  NOT NULL DEFAULT NOW() );")
	if err != nil {
		log.Printf("ERRO: creating database table repositories: %s", err)
		panic(err)
	}

	_, err = Conn.Exec("CREATE TABLE IF NOT EXISTS users (" +
		"ID  SERIAL PRIMARY KEY, " +
		"LOGIN                TEXT       NOT NULL, " +
		"EMAIL                TEXT       NOT NULL, " +
		"CREATED_AT           TIMESTAMP  NOT NULL DEFAULT NOW() );")
	if err != nil {
		log.Printf("ERRO: creating database table users: %s", err)
		panic(err)
	}
}
