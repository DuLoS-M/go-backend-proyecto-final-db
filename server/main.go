package main

import (
	"proyecto-bd-final/internal/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()

	routes.SetupRoutes(server)

	if err := server.Run(":8080"); err != nil {
		panic(err)
	}
}
