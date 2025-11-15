package controllers

import (
	"net/http"
	"proyecto-bd-final/internal/services"
	"proyecto-bd-final/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

var (
	prestamoService = services.NewPrestamoService()
	bitacoraService = services.NewBitacoraService()
)

// GetMyLoans obtiene los préstamos del usuario actual
func GetMyLoans(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Usuario no autenticado", nil)
		return
	}

	prestamos, err := prestamoService.GetPrestamosByUsuario(userID.(int))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error al obtener préstamos", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Préstamos obtenidos", prestamos)
}

// CreateLoan crea un nuevo préstamo
func CreateLoan(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Usuario no autenticado", nil)
		return
	}

	var loanData struct {
		ISBN string `json:"isbn" binding:"required"`
	}

	if err := c.ShouldBindJSON(&loanData); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Datos inválidos", err)
		return
	}

	prestamo, err := prestamoService.CrearPrestamo(userID.(int), loanData.ISBN)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Error al crear préstamo", err)
		return
	}

	// Registrar en bitácora
	bitacoraService.RegistrarAccion(userID.(int), "CREATE", "Prestamo", "Préstamo creado para libro ISBN: "+loanData.ISBN)

	utils.SuccessResponse(c, http.StatusCreated, "Préstamo creado exitosamente", prestamo)
}

// ReturnLoan registra la devolución de un préstamo
func ReturnLoan(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Usuario no autenticado", nil)
		return
	}

	prestamoID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID de préstamo inválido", err)
		return
	}

	err = prestamoService.DevolverPrestamo(prestamoID, userID.(int))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Error al devolver libro", err)
		return
	}

	// Registrar en bitácora
	bitacoraService.RegistrarAccion(userID.(int), "UPDATE", "Prestamo", "Libro devuelto - Préstamo ID: "+strconv.Itoa(prestamoID))

	utils.SuccessResponse(c, http.StatusOK, "Libro devuelto exitosamente", nil)
}

// GetAllLoans obtiene todos los préstamos (admin)
func GetAllLoans(c *gin.Context) {
	prestamos, err := prestamoService.GetTodosPrestamos()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error al obtener préstamos", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Préstamos obtenidos", prestamos)
}
