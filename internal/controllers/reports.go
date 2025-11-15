package controllers

import (
	"net/http"
	"proyecto-bd-final/internal/services"
	"proyecto-bd-final/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

var reportsService = services.NewReportsService()

// GetReportePrestamosActivos genera reporte de préstamos activos (admin)
func GetReportePrestamosActivos(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Usuario no autenticado", nil)
		return
	}

	reporte, err := reportsService.GetReportePrestamosActivos(userID.(int))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error al generar reporte", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Reporte generado exitosamente", reporte)
}

// GetReporteUsuariosActivos genera reporte de usuarios más activos (admin)
func GetReporteUsuariosActivos(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Usuario no autenticado", nil)
		return
	}

	// Obtener límite opcional
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 10
	}

	reporte, err := reportsService.GetReporteUsuariosActivos(userID.(int), limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error al generar reporte", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Reporte generado exitosamente", reporte)
}

// GetReporteLibrosPopulares genera reporte de libros más populares (admin)
func GetReporteLibrosPopulares(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Usuario no autenticado", nil)
		return
	}

	// Obtener límite opcional
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 10
	}

	reporte, err := reportsService.GetReporteLibrosPopulares(userID.(int), limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error al generar reporte", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Reporte generado exitosamente", reporte)
}

// GetEstadisticasGenerales obtiene estadísticas generales del sistema (admin)
func GetEstadisticasGenerales(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Usuario no autenticado", nil)
		return
	}

	estadisticas, err := reportsService.GetEstadisticasGenerales(userID.(int))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error al obtener estadísticas", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Estadísticas obtenidas exitosamente", estadisticas)
}
