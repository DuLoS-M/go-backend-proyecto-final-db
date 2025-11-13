package controllers

import (
	"net/http"
	"proyecto-bd-final/internal/config"
	"proyecto-bd-final/internal/models"
	"proyecto-bd-final/internal/repository"
	"proyecto-bd-final/internal/services"
	"proyecto-bd-final/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

var (
	bitacoraAdminService = services.NewBitacoraService()
	adminUserRepo        = repository.NewUserRepository()
)

// GetStatistics obtiene estadísticas del sistema (admin)
func GetStatistics(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Usuario no autenticado", nil)
		return
	}

	// Estadísticas de usuarios
	var totalUsuarios, totalEstudiantes, totalProfesores, totalPersonal int
	config.DB.QueryRow("SELECT COUNT(*) FROM Usuario").Scan(&totalUsuarios)
	config.DB.QueryRow("SELECT COUNT(*) FROM Estudiante").Scan(&totalEstudiantes)
	config.DB.QueryRow("SELECT COUNT(*) FROM Profesor").Scan(&totalProfesores)
	config.DB.QueryRow("SELECT COUNT(*) FROM Personal").Scan(&totalPersonal)

	// Estadísticas de libros
	var totalLibros, totalEjemplares int
	config.DB.QueryRow("SELECT COUNT(*) FROM Libro").Scan(&totalLibros)
	config.DB.QueryRow("SELECT COUNT(*) FROM Ejemplar").Scan(&totalEjemplares)

	// Estadísticas de préstamos
	var prestamosActivos, prestamosDevueltos, prestamosVencidos int
	config.DB.QueryRow("SELECT COUNT(*) FROM Prestamo WHERE ESTADO = 'ACTIVO'").Scan(&prestamosActivos)
	config.DB.QueryRow("SELECT COUNT(*) FROM Prestamo WHERE ESTADO = 'DEVUELTO'").Scan(&prestamosDevueltos)
	config.DB.QueryRow(`SELECT COUNT(*) FROM Prestamo 
						WHERE ESTADO = 'ACTIVO' 
						AND FECHADEVOLUCIONPREVISTA < SYSDATE`).Scan(&prestamosVencidos)

	estadisticas := map[string]interface{}{
		"usuarios": map[string]int{
			"total":       totalUsuarios,
			"estudiantes": totalEstudiantes,
			"profesores":  totalProfesores,
			"personal":    totalPersonal,
		},
		"libros": map[string]int{
			"total":      totalLibros,
			"ejemplares": totalEjemplares,
		},
		"prestamos": map[string]int{
			"activos":   prestamosActivos,
			"devueltos": prestamosDevueltos,
			"vencidos":  prestamosVencidos,
		},
	}

	// Registrar en bitácora
	bitacoraAdminService.RegistrarAccion(userID.(int), "READ", "Estadisticas", "Consulta de estadísticas del sistema")

	utils.SuccessResponse(c, http.StatusOK, "Estadísticas obtenidas", estadisticas)
}

// GetBitacora obtiene el registro de bitácora (admin)
func GetBitacora(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Usuario no autenticado", nil)
		return
	}

	// Parámetros de filtrado
	entidad := c.Query("entidad")
	limitStr := c.DefaultQuery("limit", "100")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 100
	}

	// Llamar al servicio con los parámetros en el orden correcto
	bitacoras, err := bitacoraAdminService.ObtenerBitacora(limit, entidad)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error al obtener bitácora", err)
		return
	}

	// Registrar consulta
	bitacoraAdminService.RegistrarAccion(userID.(int), "READ", "Bitacora", "Consulta de bitácora")

	utils.SuccessResponse(c, http.StatusOK, "Bitácora obtenida", bitacoras)
}

// GetRoles obtiene todos los roles (admin)
func GetRoles(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Usuario no autenticado", nil)
		return
	}

	query := `SELECT IDROL, NOMBREROL FROM Roles`
	rows, err := config.DB.Query(query)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error al obtener roles", err)
		return
	}
	defer rows.Close()

	var roles []*models.Rol
	for rows.Next() {
		var rol models.Rol
		if err := rows.Scan(&rol.IDRol, &rol.NombreRol); err != nil {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Error al procesar roles", err)
			return
		}
		roles = append(roles, &rol)
	}

	// Registrar en bitácora
	bitacoraAdminService.RegistrarAccion(userID.(int), "READ", "Roles", "Consulta de roles del sistema")

	utils.SuccessResponse(c, http.StatusOK, "Roles obtenidos", roles)
}

// AssignRole asigna un rol a un usuario (admin)
func AssignRole(c *gin.Context) {
	adminID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Usuario no autenticado", nil)
		return
	}

	var assignData struct {
		UsuarioID int `json:"usuario_id" binding:"required"`
		RolID     int `json:"rol_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&assignData); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Datos inválidos", err)
		return
	}

	// Asignar rol
	if err := adminUserRepo.AssignRole(assignData.UsuarioID, assignData.RolID); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error al asignar rol", err)
		return
	}

	// Registrar en bitácora
	bitacoraAdminService.RegistrarAccion(
		adminID.(int),
		"UPDATE",
		"UsuarioRol",
		"Rol ID "+strconv.Itoa(assignData.RolID)+" asignado al usuario ID: "+strconv.Itoa(assignData.UsuarioID),
	)

	utils.SuccessResponse(c, http.StatusOK, "Rol asignado exitosamente", nil)
}
