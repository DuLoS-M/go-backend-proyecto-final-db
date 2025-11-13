package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/godror/godror"
)

var DB *sql.DB

// InitDB inicializa la conexión a la base de datos Oracle
func InitDB() error {
	// Formato: user/password@host:port/service_name
	dsn := fmt.Sprintf("%s/%s@%s:%s/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SERVICE"),
	)

	var err error
	DB, err = sql.Open("godror", dsn)
	if err != nil {
		return fmt.Errorf("error al abrir conexión: %v", err)
	}

	// Verificar la conexión
	if err = DB.Ping(); err != nil {
		return fmt.Errorf("error al verificar conexión: %v", err)
	}

	// Configurar pool de conexiones
	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(5)

	log.Println("✅ Conexión a base de datos Oracle establecida")
	return nil
}

// CloseDB cierra la conexión a la base de datos
func CloseDB() {
	if DB != nil {
		DB.Close()
		log.Println("Conexión a base de datos cerrada")
	}
}
