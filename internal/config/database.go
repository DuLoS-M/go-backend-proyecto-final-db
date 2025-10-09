package config

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/godror/godror"
)

var DB *sql.DB

func InitDB() {
	var err error
	dsn := "user/password@host:port/service_name"
	DB, err = sql.Open("godror", dsn)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	fmt.Println("Database connection established")
}
