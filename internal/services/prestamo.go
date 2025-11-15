package services

import (
	"database/sql"
	"errors"
	"proyecto-bd-final/internal/config"
	"proyecto-bd-final/internal/models"
	"time"
)

type PrestamoService struct {
	bitacoraService *BitacoraService
}

func NewPrestamoService() *PrestamoService {
	return &PrestamoService{
		bitacoraService: NewBitacoraService(),
	}
}

// VerificarDisponibilidad verifica si hay ejemplares disponibles de un libro
func (s *PrestamoService) VerificarDisponibilidad(isbn string) (bool, int, error) {
	query := `SELECT codigo 
			  FROM Ejemplar 
			  WHERE Libro_ISBN = :1 AND estado = 'DISPONIBLE'
			  FETCH FIRST 1 ROW ONLY`

	var codigoEjemplar int
	err := config.DB.QueryRow(query, isbn).Scan(&codigoEjemplar)

	if err == sql.ErrNoRows {
		return false, 0, nil
	}

	if err != nil {
		return false, 0, err
	}

	return true, codigoEjemplar, nil
}

// CrearPrestamo crea un nuevo préstamo y actualiza el estado del ejemplar
func (s *PrestamoService) CrearPrestamo(usuarioID int, libroISBN string) (*models.Prestamo, error) {
	// Verificar disponibilidad
	disponible, codigoEjemplar, err := s.VerificarDisponibilidad(libroISBN)
	if err != nil {
		return nil, err
	}

	if !disponible {
		return nil, errors.New("no hay ejemplares disponibles para este libro")
	}

	// Iniciar transacción
	tx, err := config.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Crear préstamo
	fechaPrestamo := time.Now()
	fechaDevolucion := fechaPrestamo.AddDate(0, 0, 15) // 15 días

	// Obtener el próximo ID de la secuencia
	var idPrestamo int
	err = tx.QueryRow("SELECT PRESTAMO_SEQ.NEXTVAL FROM DUAL").Scan(&idPrestamo)
	if err != nil {
		return nil, err
	}

	// Insertar el préstamo
	queryPrestamo := `INSERT INTO Prestamo 
					  (IDPRESTAMO, FECHAPRESTAMO, FECHADEVOLUCIONPREVISTA, ESTADO, USUARIO_IDUSUARIO, DEVOLUCION_IDDEVOLUCION) 
					  VALUES (:1, :2, :3, :4, :5, NULL)`

	_, err = tx.Exec(queryPrestamo, idPrestamo, fechaPrestamo, fechaDevolucion, "ACTIVO", usuarioID)
	if err != nil {
		return nil, err
	}

	// Actualizar estado del ejemplar
	queryEjemplar := `UPDATE Ejemplar 
					  SET estado = 'PRESTADO', Prestamo_idPrestamo = :1 
					  WHERE codigo = :2`

	_, err = tx.Exec(queryEjemplar, idPrestamo, codigoEjemplar)
	if err != nil {
		return nil, err
	}

	// Commit transacción
	if err = tx.Commit(); err != nil {
		return nil, err
	}

	prestamo := &models.Prestamo{
		IDPrestamo:              idPrestamo,
		FechaPrestamo:           fechaPrestamo,
		FechaDevolucionPrevista: fechaDevolucion,
		Estado:                  "ACTIVO",
		UsuarioID:               usuarioID,
	}

	return prestamo, nil
}

// DevolverPrestamo registra la devolución de un libro
func (s *PrestamoService) DevolverPrestamo(prestamoID, usuarioID int) error {
	// Verificar que el préstamo pertenece al usuario
	var usuarioIDPrestamo int
	queryVerificar := `SELECT USUARIO_IDUSUARIO FROM Prestamo WHERE IDPRESTAMO = :1`
	err := config.DB.QueryRow(queryVerificar, prestamoID).Scan(&usuarioIDPrestamo)
	if err != nil {
		return errors.New("préstamo no encontrado")
	}

	if usuarioIDPrestamo != usuarioID {
		return errors.New("no tienes permiso para devolver este préstamo")
	}

	// Iniciar transacción
	tx, err := config.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Actualizar préstamo
	queryPrestamo := `UPDATE Prestamo 
					  SET FECHADEVOLUCIONREAL = :1, ESTADO = :2 
					  WHERE IDPRESTAMO = :3`

	_, err = tx.Exec(queryPrestamo, time.Now(), "DEVUELTO", prestamoID)
	if err != nil {
		return err
	}

	// Actualizar estado del ejemplar
	queryEjemplar := `UPDATE Ejemplar 
					  SET estado = 'DISPONIBLE', Prestamo_idPrestamo = NULL
					  WHERE Prestamo_idPrestamo = :1`

	_, err = tx.Exec(queryEjemplar, prestamoID)
	if err != nil {
		return err
	}

	// Commit transacción
	return tx.Commit()
}

// GetPrestamosByUsuario obtiene todos los préstamos de un usuario
func (s *PrestamoService) GetPrestamosByUsuario(usuarioID int) ([]*models.Prestamo, error) {
	query := `SELECT P.IDPRESTAMO, P.FECHAPRESTAMO, P.FECHADEVOLUCIONPREVISTA, 
			  P.FECHADEVOLUCIONREAL, P.ESTADO, P.USUARIO_IDUSUARIO
			  FROM Prestamo P
			  WHERE P.USUARIO_IDUSUARIO = :1
			  ORDER BY P.FECHAPRESTAMO DESC`

	rows, err := config.DB.Query(query, usuarioID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prestamos []*models.Prestamo
	for rows.Next() {
		var prestamo models.Prestamo
		var fechaDevolucionReal sql.NullTime

		if err := rows.Scan(
			&prestamo.IDPrestamo,
			&prestamo.FechaPrestamo,
			&prestamo.FechaDevolucionPrevista,
			&fechaDevolucionReal,
			&prestamo.Estado,
			&prestamo.UsuarioID,
		); err != nil {
			return nil, err
		}

		if fechaDevolucionReal.Valid {
			t := fechaDevolucionReal.Time
			prestamo.FechaDevolucionReal = &t
		}

		prestamos = append(prestamos, &prestamo)
	}

	return prestamos, nil
}

// GetTodosPrestamos obtiene todos los préstamos (admin)
func (s *PrestamoService) GetTodosPrestamos() ([]*models.Prestamo, error) {
	query := `SELECT P.IDPRESTAMO, P.FECHAPRESTAMO, P.FECHADEVOLUCIONPREVISTA, 
			  P.FECHADEVOLUCIONREAL, P.ESTADO, P.USUARIO_IDUSUARIO
			  FROM Prestamo P
			  ORDER BY P.FECHAPRESTAMO DESC`

	rows, err := config.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prestamos []*models.Prestamo
	for rows.Next() {
		var prestamo models.Prestamo
		var fechaDevolucionReal sql.NullTime

		if err := rows.Scan(
			&prestamo.IDPrestamo,
			&prestamo.FechaPrestamo,
			&prestamo.FechaDevolucionPrevista,
			&fechaDevolucionReal,
			&prestamo.Estado,
			&prestamo.UsuarioID,
		); err != nil {
			return nil, err
		}

		if fechaDevolucionReal.Valid {
			t := fechaDevolucionReal.Time
			prestamo.FechaDevolucionReal = &t
		}

		prestamos = append(prestamos, &prestamo)
	}

	return prestamos, nil
}
