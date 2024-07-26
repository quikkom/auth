package database

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
)

var DBConn *pgx.Conn

// Establish connection to the database
func CreateDBConnection() {
	databaseUrl, found := os.LookupEnv("DATABASE_URL")

	if !found {
		log.Panic("DATABASE_URL environment variable is not defined")
	}

	conn, err := pgx.Connect(context.Background(), databaseUrl)

	if err != nil {
		log.Panicf("An error occurred when trying to connect to the database: %v", err)
	}

	DBConn = conn
}
