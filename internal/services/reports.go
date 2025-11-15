package services

import (
	"proyecto-bd-final/internal/config"
)

type ReportsService struct {
	bitacoraService *BitacoraService
}

func NewReportsService() *ReportsService {
	return &ReportsService{
		bitacoraService: NewBitacoraService(),
	}
}

// ReportePrestamosActivosResponse estructura para el reporte de préstamos activos
type ReportePrestamosActivosResponse struct {
	Total        int                   `json:"total"`
	Vencidos     int                   `json:"vencidos"`
	PorVencer    int                   `json:"por_vencer"`
	Prestamos    []PrestamoDetalleInfo `json:"prestamos"`
	FechaReporte string                `json:"fecha_reporte"`
}

// PrestamoDetalleInfo información detallada de préstamos para reportes
type PrestamoDetalleInfo struct {
	IDPrestamo              int    `json:"id_prestamo"`
	FechaPrestamo           string `json:"fecha_prestamo"`
	FechaDevolucionPrevista string `json:"fecha_devolucion_prevista"`
	DiasPrestamo            int    `json:"dias_prestamo"`
	Estado                  string `json:"estado"`
	UsuarioNombre           string `json:"usuario_nombre"`
	UsuarioApellido         string `json:"usuario_apellido"`
	LibroTitulo             string `json:"libro_titulo"`
	LibroISBN               string `json:"libro_isbn"`
}

// ReporteUsuariosActivosResponse estructura para el reporte de usuarios más activos
type ReporteUsuariosActivosResponse struct {
	UsuariosActivos []UsuarioActivoInfo `json:"usuarios_activos"`
	FechaReporte    string              `json:"fecha_reporte"`
}

// UsuarioActivoInfo información de usuarios para reportes
type UsuarioActivoInfo struct {
	UsuarioID          int    `json:"usuario_id"`
	NombreCompleto     string `json:"nombre_completo"`
	TotalPrestamos     int    `json:"total_prestamos"`
	PrestamosActivos   int    `json:"prestamos_activos"`
	PrestamosDevueltos int    `json:"prestamos_devueltos"`
}

// ReporteLibrosPopularesResponse estructura para el reporte de libros más solicitados
type ReporteLibrosPopularesResponse struct {
	LibrosPopulares []LibroPopularInfo `json:"libros_populares"`
	FechaReporte    string             `json:"fecha_reporte"`
}

// LibroPopularInfo información de libros para reportes
type LibroPopularInfo struct {
	ISBN             string `json:"isbn"`
	Titulo           string `json:"titulo"`
	TotalPrestamos   int    `json:"total_prestamos"`
	PrestamosActivos int    `json:"prestamos_activos"`
	Editorial        string `json:"editorial"`
}

// EstadisticasGeneralesResponse estructura para estadísticas generales del sistema
type EstadisticasGeneralesResponse struct {
	TotalUsuarios      int    `json:"total_usuarios"`
	TotalLibros        int    `json:"total_libros"`
	TotalEjemplares    int    `json:"total_ejemplares"`
	PrestamosActivos   int    `json:"prestamos_activos"`
	PrestamosDevueltos int    `json:"prestamos_devueltos"`
	PrestamosVencidos  int    `json:"prestamos_vencidos"`
	FechaReporte       string `json:"fecha_reporte"`
}

// GetReportePrestamosActivos genera un reporte de todos los préstamos activos
func (s *ReportsService) GetReportePrestamosActivos(userID int) (*ReportePrestamosActivosResponse, error) {
	query := `SELECT 
				P.IDPRESTAMO, 
				TO_CHAR(P.FECHAPRESTAMO, 'YYYY-MM-DD') as FECHA_PRESTAMO,
				TO_CHAR(P.FECHADEVOLUCIONPREVISTA, 'YYYY-MM-DD') as FECHA_DEVOLUCION_PREVISTA,
				TRUNC(SYSDATE - P.FECHAPRESTAMO) as DIAS_PRESTAMO,
				P.ESTADO,
				U.NOMBRE as USUARIO_NOMBRE,
				U.APELLIDO as USUARIO_APELLIDO,
				L.TITULO as LIBRO_TITULO,
				L.ISBN as LIBRO_ISBN
			  FROM Prestamo P
			  INNER JOIN Usuario U ON P.USUARIO_IDUSUARIO = U.IDUSUARIO
			  INNER JOIN Ejemplar E ON P.IDPRESTAMO = E.Prestamo_idPrestamo
			  INNER JOIN Libro L ON E.Libro_ISBN = L.ISBN
			  WHERE P.ESTADO = 'ACTIVO'
			  ORDER BY P.FECHADEVOLUCIONPREVISTA ASC`

	rows, err := config.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prestamos []PrestamoDetalleInfo
	totalActivos := 0
	vencidos := 0
	porVencer := 0

	for rows.Next() {
		var prestamo PrestamoDetalleInfo
		if err := rows.Scan(
			&prestamo.IDPrestamo,
			&prestamo.FechaPrestamo,
			&prestamo.FechaDevolucionPrevista,
			&prestamo.DiasPrestamo,
			&prestamo.Estado,
			&prestamo.UsuarioNombre,
			&prestamo.UsuarioApellido,
			&prestamo.LibroTitulo,
			&prestamo.LibroISBN,
		); err != nil {
			return nil, err
		}

		prestamos = append(prestamos, prestamo)
		totalActivos++

		// Verificar si está vencido (más de 15 días)
		if prestamo.DiasPrestamo > 15 {
			vencidos++
		} else if prestamo.DiasPrestamo > 12 { // Por vencer en 3 días o menos
			porVencer++
		}
	}

	// Registrar en bitácora
	s.bitacoraService.RegistrarAccion(userID, "READ", "Reporte", "Generación de reporte de préstamos activos")

	return &ReportePrestamosActivosResponse{
		Total:        totalActivos,
		Vencidos:     vencidos,
		PorVencer:    porVencer,
		Prestamos:    prestamos,
		FechaReporte: "SYSDATE",
	}, nil
}

// GetReporteUsuariosActivos genera un reporte de los usuarios más activos
func (s *ReportsService) GetReporteUsuariosActivos(userID int, limite int) (*ReporteUsuariosActivosResponse, error) {
	query := `SELECT 
				U.IDUSUARIO,
				U.NOMBRE || ' ' || U.APELLIDO as NOMBRE_COMPLETO,
				COUNT(*) as TOTAL_PRESTAMOS,
				COUNT(CASE WHEN P.ESTADO = 'ACTIVO' THEN 1 END) as PRESTAMOS_ACTIVOS,
				COUNT(CASE WHEN P.ESTADO = 'DEVUELTO' THEN 1 END) as PRESTAMOS_DEVUELTOS
			  FROM Usuario U
			  INNER JOIN Prestamo P ON U.IDUSUARIO = P.USUARIO_IDUSUARIO
			  GROUP BY U.IDUSUARIO, U.NOMBRE, U.APELLIDO
			  ORDER BY TOTAL_PRESTAMOS DESC
			  FETCH FIRST :1 ROWS ONLY`

	rows, err := config.DB.Query(query, limite)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var usuariosActivos []UsuarioActivoInfo
	for rows.Next() {
		var usuario UsuarioActivoInfo
		if err := rows.Scan(
			&usuario.UsuarioID,
			&usuario.NombreCompleto,
			&usuario.TotalPrestamos,
			&usuario.PrestamosActivos,
			&usuario.PrestamosDevueltos,
		); err != nil {
			return nil, err
		}
		usuariosActivos = append(usuariosActivos, usuario)
	}

	// Registrar en bitácora
	s.bitacoraService.RegistrarAccion(userID, "READ", "Reporte", "Generación de reporte de usuarios activos")

	return &ReporteUsuariosActivosResponse{
		UsuariosActivos: usuariosActivos,
		FechaReporte:    "SYSDATE",
	}, nil
}

// GetReporteLibrosPopulares genera un reporte de los libros más solicitados
func (s *ReportsService) GetReporteLibrosPopulares(userID int, limite int) (*ReporteLibrosPopularesResponse, error) {
	query := `SELECT 
				L.ISBN,
				L.TITULO,
				COUNT(*) as TOTAL_PRESTAMOS,
				COUNT(CASE WHEN P.ESTADO = 'ACTIVO' THEN 1 END) as PRESTAMOS_ACTIVOS,
				E.NOMBRE as EDITORIAL
			  FROM Libro L
			  INNER JOIN Ejemplar EJ ON L.ISBN = EJ.Libro_ISBN
			  INNER JOIN Prestamo P ON EJ.Prestamo_idPrestamo = P.IDPRESTAMO
			  INNER JOIN Editorial E ON L.Editorial_idEditorial = E.IDEDITORIAL
			  GROUP BY L.ISBN, L.TITULO, E.NOMBRE
			  ORDER BY TOTAL_PRESTAMOS DESC
			  FETCH FIRST :1 ROWS ONLY`

	rows, err := config.DB.Query(query, limite)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var librosPopulares []LibroPopularInfo
	for rows.Next() {
		var libro LibroPopularInfo
		if err := rows.Scan(
			&libro.ISBN,
			&libro.Titulo,
			&libro.TotalPrestamos,
			&libro.PrestamosActivos,
			&libro.Editorial,
		); err != nil {
			return nil, err
		}
		librosPopulares = append(librosPopulares, libro)
	}

	// Registrar en bitácora
	s.bitacoraService.RegistrarAccion(userID, "READ", "Reporte", "Generación de reporte de libros populares")

	return &ReporteLibrosPopularesResponse{
		LibrosPopulares: librosPopulares,
		FechaReporte:    "SYSDATE",
	}, nil
}

// GetEstadisticasGenerales genera estadísticas generales del sistema
func (s *ReportsService) GetEstadisticasGenerales(userID int) (*EstadisticasGeneralesResponse, error) {
	var stats EstadisticasGeneralesResponse

	// Total de usuarios
	err := config.DB.QueryRow("SELECT COUNT(*) FROM Usuario").Scan(&stats.TotalUsuarios)
	if err != nil {
		return nil, err
	}

	// Total de libros
	err = config.DB.QueryRow("SELECT COUNT(*) FROM Libro").Scan(&stats.TotalLibros)
	if err != nil {
		return nil, err
	}

	// Total de ejemplares
	err = config.DB.QueryRow("SELECT COUNT(*) FROM Ejemplar").Scan(&stats.TotalEjemplares)
	if err != nil {
		return nil, err
	}

	// Préstamos activos
	err = config.DB.QueryRow("SELECT COUNT(*) FROM Prestamo WHERE ESTADO = 'ACTIVO'").Scan(&stats.PrestamosActivos)
	if err != nil {
		return nil, err
	}

	// Préstamos devueltos
	err = config.DB.QueryRow("SELECT COUNT(*) FROM Prestamo WHERE ESTADO = 'DEVUELTO'").Scan(&stats.PrestamosDevueltos)
	if err != nil {
		return nil, err
	}

	// Préstamos vencidos
	err = config.DB.QueryRow(`SELECT COUNT(*) FROM Prestamo 
								WHERE ESTADO = 'ACTIVO' 
								AND FECHADEVOLUCIONPREVISTA < SYSDATE`).Scan(&stats.PrestamosVencidos)
	if err != nil {
		return nil, err
	}

	stats.FechaReporte = "SYSDATE"

	// Registrar en bitácora
	s.bitacoraService.RegistrarAccion(userID, "READ", "Reporte", "Generación de estadísticas generales")

	return &stats, nil
}
