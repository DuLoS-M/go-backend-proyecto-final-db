package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/godror/godror"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func main() {
	// Configuración directa del archivo .env
	envPath := "E:/S6/BD/go-backend-proyecto-final-db/.env"

	// Cargar variables de entorno
	if err := godotenv.Load(envPath); err != nil {
		log.Printf("⚠️ No se encontró archivo .env en %s, usando variables de entorno del sistema", envPath)
	}

	// Verificar que las variables existen
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbService := os.Getenv("DB_SERVICE")

	if dbUser == "" || dbPassword == "" || dbHost == "" || dbPort == "" || dbService == "" {
		log.Fatal("❌ Variables de entorno de base de datos no configuradas. Verifica tu archivo .env")
	}

	log.Printf("📝 Conectando a: %s@%s:%s/%s", dbUser, dbHost, dbPort, dbService)

	// Conectar a la base de datos
	dsn := fmt.Sprintf("%s/%s@%s:%s/%s",
		dbUser,
		dbPassword,
		dbHost,
		dbPort,
		dbService,
	)

	db, err := sql.Open("godror", dsn)
	if err != nil {
		log.Fatalf("❌ Error al conectar: %v", err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatalf("❌ Error al verificar conexión: %v", err)
	}

	log.Println("✅ Conectado a Oracle")

	log.Println("🔧 Actualizando contraseñas de usuarios existentes...")

	// Usuarios y contraseñas a actualizar
	usersToUpdate := []struct {
		Email    string
		Password string
	}{
		{"admin@biblioteca.edu", "admin123"},
		{"juan.perez@estudiante.edu", "estudiante123"},
		{"maria.lopez@profesor.edu", "profesor123"},
		{"carlos.garcia@biblioteca.edu", "personal123"},
		{"ana.martinez@estudiante.edu", "estudiante123"},
		{"pedro.rodriguez@estudiante.edu", "estudiante123"},
		// Usuarios de demostración adicionales
		{"carlos.mendoza@estudiante.com", "demo123"},
		{"ana.ramirez@estudiante.com", "demo123"},
		{"patricia.lopez@profesor.com", "demo123"},
		{"estudiante1@demo.com", "demo123"},
		{"estudiante2@demo.com", "demo123"},
		{"estudiante3@demo.com", "demo123"},
		{"estudiante4@demo.com", "demo123"},
		{"estudiante5@demo.com", "demo123"},
	}

	log.Println("� Actualizando contraseñas...")

	var updatedCount int
	var notFoundCount int

	for _, user := range usersToUpdate {
		// Hash de contraseña
		hashedPassword, err := hashPassword(user.Password)
		if err != nil {
			log.Printf("❌ Error al hashear contraseña para %s: %v", user.Email, err)
			continue
		}

		// Verificar si el usuario existe
		var userID int
		checkQuery := `SELECT IDUSUARIO FROM Usuario WHERE CORREO = :1`
		err = db.QueryRow(checkQuery, user.Email).Scan(&userID)
		if err != nil {
			if err == sql.ErrNoRows {
				log.Printf("⚠️ Usuario no encontrado: %s", user.Email)
				notFoundCount++
			} else {
				log.Printf("❌ Error al verificar usuario %s: %v", user.Email, err)
			}
			continue
		}

		// Actualizar contraseña
		updateQuery := `UPDATE Usuario SET CONTRASENIA = :1 WHERE CORREO = :2`
		result, err := db.Exec(updateQuery, hashedPassword, user.Email)
		if err != nil {
			log.Printf("❌ Error al actualizar contraseña para %s: %v", user.Email, err)
			continue
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected > 0 {
			log.Printf("✅ Contraseña actualizada: %s -> %s", user.Email, user.Password)
			updatedCount++
		} else {
			log.Printf("⚠️ No se actualizó: %s", user.Email)
		}
	}

	log.Printf("\n🎉 Actualización completada!")
	log.Printf("📊 Usuarios actualizados: %d", updatedCount)
	log.Printf("📊 Usuarios no encontrados: %d", notFoundCount)

	log.Println("\n� CREDENCIALES ACTUALIZADAS:")
	log.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	log.Println("USUARIOS PRINCIPALES:")
	log.Println("Admin:           admin@biblioteca.edu / admin123")
	log.Println("Estudiante:      juan.perez@estudiante.edu / estudiante123")
	log.Println("Profesor:        maria.lopez@profesor.edu / profesor123")
	log.Println("Personal:        carlos.garcia@biblioteca.edu / personal123")
	log.Println("")
	log.Println("USUARIOS DE DEMOSTRACIÓN:")
	log.Println("Estudiante demo: carlos.mendoza@estudiante.com / demo123")
	log.Println("Estudiante demo: ana.ramirez@estudiante.com / demo123")
	log.Println("Profesor demo:   patricia.lopez@profesor.com / demo123")
	log.Println("Estudiantes:     estudiante1-5@demo.com / demo123")
	log.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
}
