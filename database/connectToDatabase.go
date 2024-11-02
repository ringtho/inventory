package database

import (
	"database/sql"
	"log"
	"os"
)

func ConnectToDatabase() *sql.DB {
	dbURL := os.Getenv("DB_URL")

	if dbURL == "" {
		log.Fatal("DB_URL not found in the environment")
	}
	conn, err := sql.Open("postgres", dbURL)

	if err != nil {
		log.Fatal("Cant connect to database ", err)
	}
	log.Printf("Connected to database!")
	return conn
}