package routes

import (
	"proyecto-bd-final/internal/controllers"
	"proyecto-bd-final/internal/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configura todas las rutas de la aplicación
func SetupRoutes(router *gin.Engine) {
	// Rutas públicas
	public := router.Group("/api")
	{
		public.POST("/auth/login", controllers.Login)
		public.POST("/auth/register", controllers.Register)
	}

	// Rutas protegidas (requieren autenticación)
	protected := router.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		// Rutas de usuario
		protected.GET("/profile", controllers.GetProfile)
		protected.PUT("/profile", controllers.UpdateProfile)

		// Rutas de libros
		protected.GET("/books", controllers.GetBooks)
		protected.GET("/books/:isbn", controllers.GetBookByISBN)

		// Rutas de préstamos
		protected.GET("/loans/my-loans", controllers.GetMyLoans)
		protected.POST("/loans", controllers.CreateLoan)
		protected.PUT("/loans/:id/return", controllers.ReturnLoan)

		// Rutas de admin
		admin := protected.Group("/admin")
		admin.Use(middleware.RequireRole("admin"))
		{
			admin.GET("/users", controllers.GetAllUsers)
			admin.GET("/statistics", controllers.GetStatistics)
			admin.GET("/bitacora", controllers.GetBitacora)
			admin.GET("/loans", controllers.GetAllLoans) // Nuevo endpoint para admin

			// Gestión de libros (admin)
			admin.POST("/books", controllers.CreateBook)
			admin.PUT("/books/:isbn", controllers.UpdateBook)
			admin.DELETE("/books/:isbn", controllers.DeleteBook)

			// Gestión de roles
			admin.GET("/roles", controllers.GetRoles)
			admin.POST("/users/:id/roles", controllers.AssignRole)

			// Reportes (admin)
			admin.GET("/reports/prestamos-activos", controllers.GetReportePrestamosActivos)
			admin.GET("/reports/usuarios-activos", controllers.GetReporteUsuariosActivos)
			admin.GET("/reports/libros-populares", controllers.GetReporteLibrosPopulares)
			admin.GET("/reports/estadisticas", controllers.GetEstadisticasGenerales)
		}
	}
}
