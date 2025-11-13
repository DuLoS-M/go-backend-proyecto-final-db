package controllers

import (
	"net/http"
	"proyecto-bd-final/internal/config"
	"proyecto-bd-final/internal/models"
	"proyecto-bd-final/internal/services"
	"proyecto-bd-final/pkg/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

var bitacoraService = services.NewBitacoraService()

// GetMyLoans obtiene los préstamos del usuario actual
func GetMyLoans(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Usuario no autenticado", nil)
		return
	}

	query := `SELECT P.IDPRESTAMO, P.FECHAPRESTAMO, P.FECHADEVOLUCIONPREVISTA, 
			  P.FECHADEVOLUCIONREAL, P.ESTADO, P.USUARIO_IDUSUARIO
			  FROM Prestamo P
			  WHERE P.USUARIO_IDUSUARIO = :1
			  ORDER BY P.FECHAPRESTAMO DESC`

	rows, err := config.DB.Query(query, userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error al obtener préstamos", err)
		return
	}
	defer rows.Close()

	var prestamos []*models.Prestamo
	for rows.Next() {
		var prestamo models.Prestamo
		if err := rows.Scan(
			&prestamo.IDPrestamo,
			&prestamo.FechaPrestamo,
			&prestamo.FechaDevolucionPrevista,
			&prestamo.FechaDevolucionReal,
			&prestamo.Estado,
			&prestamo.UsuarioID,
		); err != nil {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Error al procesar préstamos", err)
			return
		}
		prestamos = append(prestamos, &prestamo)
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
		LibroISBN int `json:"libro_isbn" binding:"required"`
	}

	if err := c.ShouldBindJSON(&loanData); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Datos inválidos", err)
		return
	}

	// Verificar disponibilidad (simplificado)
	// TODO: Implementar verificación de ejemplares disponibles

	// Crear préstamo
	fechaPrestamo := time.Now()
	fechaDevolucion := fechaPrestamo.AddDate(0, 0, 15) // 15 días

	query := `INSERT INTO Prestamo (IDPRESTAMO, FECHAPRESTAMO, FECHADEVOLUCIONPREVISTA, 
			  ESTADO, USUARIO_IDUSUARIO, DEVOLUCION_IDDEVOLUCION) 
			  VALUES (PRESTAMO_SEQ.NEXTVAL, :1, :2, :3, :4, NULL)`

	_, err := config.DB.Exec(query, fechaPrestamo, fechaDevolucion, "ACTIVO", userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error al crear préstamo", err)
		return
	}

	// Registrar en bitácora
	bitacoraService.RegistrarAccion(userID.(int), "CREATE", "Prestamo", "Préstamo creado para libro ISBN: "+strconv.Itoa(loanData.LibroISBN))

	utils.SuccessResponse(c, http.StatusCreated, "Préstamo creado exitosamente", nil)
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

	// Actualizar préstamo
	query := `UPDATE Prestamo 
			  SET FECHADEVOLUCIONREAL = :1, ESTADO = :2 
			  WHERE IDPRESTAMO = :3 AND USUARIO_IDUSUARIO = :4`

	_, err = config.DB.Exec(query, time.Now(), "DEVUELTO", prestamoID, userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error al devolver libro", err)
		return
	}

	// Registrar en bitácora
	bitacoraService.RegistrarAccion(userID.(int), "UPDATE", "Prestamo", "Libro devuelto - Préstamo ID: "+strconv.Itoa(prestamoID))

	utils.SuccessResponse(c, http.StatusOK, "Libro devuelto exitosamente", nil)
}
