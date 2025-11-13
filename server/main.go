package main

import (
	"log"
	"os"

	"proyecto-bd-final/internal/config"
	"proyecto-bd-final/internal/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Cargar variables de entorno
	if err := godotenv.Load(); err != nil {
		log.Println("No se encontró archivo .env, usando variables de entorno del sistema")
	}

	// Inicializar base de datos
	if err := config.InitDB(); err != nil {
		log.Fatalf("❌ Error al conectar a la base de datos: %v", err)
	}
	defer config.CloseDB()

	// Configurar Gin
	router := gin.Default()

	// Configurar CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // Vite dev server
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Configurar rutas
	routes.SetupRoutes(router)

	// Iniciar servidor
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Servidor iniciado en puerto %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("❌ Error al iniciar servidor: %v", err)
	}
}
