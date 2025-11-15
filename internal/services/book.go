package services

import (
	"database/sql"
	"proyecto-bd-final/internal/config"
	"proyecto-bd-final/internal/models"
)

type BookService struct {
	bitacoraService *BitacoraService
}

func NewBookService() *BookService {
	return &BookService{
		bitacoraService: NewBitacoraService(),
	}
}

// GetAll obtiene todos los libros con sus autores y editorial
func (s *BookService) GetAll() ([]*models.Libro, error) {
	query := `SELECT DISTINCT L.ISBN, L.titulo, EXTRACT(YEAR FROM L.anioEdicion) as anio,
              L.Editorial_idEditorial, E.nombre AS EDITORIAL_NOMBRE
              FROM Libro L
              LEFT JOIN Editorial E ON L.Editorial_idEditorial = E.idEditorial
              ORDER BY L.titulo`

	rows, err := config.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var libros []*models.Libro
	for rows.Next() {
		var libro models.Libro
		var editorialNombre sql.NullString

		if err := rows.Scan(
			&libro.ISBN,
			&libro.Titulo,
			&libro.AnioPublicacion,
			&libro.EditorialID,
			&editorialNombre,
		); err != nil {
			return nil, err
		}

		if editorialNombre.Valid {
			libro.EditorialNombre = editorialNombre.String
		}

		// Obtener cantidad de ejemplares
		cantidad, _ := s.GetCantidadEjemplares(libro.ISBN)
		libro.Cantidad = cantidad

		// Obtener autores del libro
		autores, err := s.GetAutoresByISBN(libro.ISBN)
		if err == nil {
			libro.Autores = autores
		}

		// Verificar disponibilidad
		disponible, _ := s.VerificarDisponibilidad(libro.ISBN)
		libro.Disponible = disponible

		libros = append(libros, &libro)
	}

	return libros, nil
}

// GetByISBN obtiene un libro por ISBN
func (s *BookService) GetByISBN(isbn string) (*models.Libro, error) {
	query := `SELECT L.ISBN, L.titulo, EXTRACT(YEAR FROM L.anioEdicion) as anio,
              L.Editorial_idEditorial, E.nombre AS EDITORIAL_NOMBRE
              FROM Libro L
              LEFT JOIN Editorial E ON L.Editorial_idEditorial = E.idEditorial
              WHERE L.ISBN = :1`

	var libro models.Libro
	var editorialNombre sql.NullString

	err := config.DB.QueryRow(query, isbn).Scan(
		&libro.ISBN,
		&libro.Titulo,
		&libro.AnioPublicacion,
		&libro.EditorialID,
		&editorialNombre,
	)

	if err != nil {
		return nil, err
	}

	if editorialNombre.Valid {
		libro.EditorialNombre = editorialNombre.String
	}

	// Obtener cantidad de ejemplares
	cantidad, _ := s.GetCantidadEjemplares(isbn)
	libro.Cantidad = cantidad

	// Obtener autores
	autores, err := s.GetAutoresByISBN(isbn)
	if err == nil {
		libro.Autores = autores
	}

	// Verificar disponibilidad
	disponible, _ := s.VerificarDisponibilidad(isbn)
	libro.Disponible = disponible

	return &libro, nil
}

// GetAutoresByISBN obtiene los autores de un libro
func (s *BookService) GetAutoresByISBN(isbn string) ([]string, error) {
	query := `SELECT A.nombre || ' ' || A.apellido AS NombreCompleto
              FROM Autor A
              INNER JOIN LibroAutor LA ON A.idAutor = LA.Autor_idAutor
              WHERE LA.Libro_ISBN = :1`

	rows, err := config.DB.Query(query, isbn)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var autores []string
	for rows.Next() {
		var nombre string
		if err := rows.Scan(&nombre); err != nil {
			return nil, err
		}
		autores = append(autores, nombre)
	}

	return autores, nil
}

// VerificarDisponibilidad verifica si hay ejemplares disponibles
func (s *BookService) VerificarDisponibilidad(isbn string) (bool, error) {
	query := `SELECT COUNT(*) 
              FROM Ejemplar 
              WHERE Libro_ISBN = :1 AND estado = 'DISPONIBLE'`

	var count int
	err := config.DB.QueryRow(query, isbn).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// GetCantidadEjemplares obtiene el total de ejemplares de un libro
func (s *BookService) GetCantidadEjemplares(isbn string) (int, error) {
	query := `SELECT COUNT(*) 
              FROM Ejemplar 
              WHERE Libro_ISBN = :1`

	var count int
	err := config.DB.QueryRow(query, isbn).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// SearchBooks busca libros por título o autor
func (s *BookService) SearchBooks(searchTerm string) ([]*models.Libro, error) {
	query := `SELECT DISTINCT L.ISBN, L.titulo, EXTRACT(YEAR FROM L.anioEdicion) as anio,
              L.Editorial_idEditorial, E.nombre AS EDITORIAL_NOMBRE
              FROM Libro L
              LEFT JOIN Editorial E ON L.Editorial_idEditorial = E.idEditorial
              LEFT JOIN LibroAutor LA ON L.ISBN = LA.Libro_ISBN
              LEFT JOIN Autor A ON LA.Autor_idAutor = A.idAutor
              WHERE LOWER(L.titulo) LIKE '%' || LOWER(:1) || '%'
              OR LOWER(A.nombre) LIKE '%' || LOWER(:2) || '%'
              OR LOWER(A.apellido) LIKE '%' || LOWER(:3) || '%'
              ORDER BY L.titulo`

	rows, err := config.DB.Query(query, searchTerm, searchTerm, searchTerm)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var libros []*models.Libro
	for rows.Next() {
		var libro models.Libro
		var editorialNombre sql.NullString

		if err := rows.Scan(
			&libro.ISBN,
			&libro.Titulo,
			&libro.AnioPublicacion,
			&libro.EditorialID,
			&editorialNombre,
		); err != nil {
			return nil, err
		}

		if editorialNombre.Valid {
			libro.EditorialNombre = editorialNombre.String
		}

		// Obtener cantidad de ejemplares
		cantidad, _ := s.GetCantidadEjemplares(libro.ISBN)
		libro.Cantidad = cantidad

		autores, err := s.GetAutoresByISBN(libro.ISBN)
		if err == nil {
			libro.Autores = autores
		}

		disponible, _ := s.VerificarDisponibilidad(libro.ISBN)
		libro.Disponible = disponible

		libros = append(libros, &libro)
	}

	return libros, nil
}

// Create crea un nuevo libro
func (s *BookService) Create(libro *models.Libro, userID int) error {
	tx, err := config.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insertar libro usando TO_DATE para convertir el año
	query := `INSERT INTO Libro (ISBN, titulo, anioEdicion, Editorial_idEditorial) 
              VALUES (:1, :2, TO_DATE(:3, 'YYYY'), :4)`

	_, err = tx.Exec(query, libro.ISBN, libro.Titulo, libro.AnioPublicacion, libro.EditorialID)
	if err != nil {
		return err
	}

	// Crear ejemplares automáticamente según la cantidad especificada
	for i := 0; i < libro.Cantidad; i++ {
		ejemplarQuery := `INSERT INTO Ejemplar (codigo, estado, Libro_ISBN, Prestamo_idPrestamo) 
                          VALUES (EJEMPLAR_SEQ.NEXTVAL, 'DISPONIBLE', :1, NULL)`
		_, err = tx.Exec(ejemplarQuery, libro.ISBN)
		if err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	// Registrar en bitácora
	s.bitacoraService.RegistrarAccion(userID, "CREATE", "LIBRO", "Libro creado: "+libro.Titulo)

	return nil
}

// Update actualiza un libro existente
func (s *BookService) Update(libro *models.Libro, userID int) error {
	query := `UPDATE Libro 
              SET titulo = :1, anioEdicion = TO_DATE(:2, 'YYYY'), Editorial_idEditorial = :3
              WHERE ISBN = :4`

	_, err := config.DB.Exec(query, libro.Titulo, libro.AnioPublicacion, libro.EditorialID, libro.ISBN)
	if err != nil {
		return err
	}

	// Registrar en bitácora
	s.bitacoraService.RegistrarAccion(userID, "UPDATE", "LIBRO", "Libro actualizado: "+libro.Titulo)

	return nil
}

// Delete elimina un libro (soft delete marcando ejemplares como no disponibles)
func (s *BookService) Delete(isbn string, userID int) error {
	// Verificar que no haya préstamos activos
	checkQuery := `SELECT COUNT(*) FROM Prestamo P 
                   INNER JOIN Ejemplar E ON P.idPrestamo = E.Prestamo_idPrestamo
                   WHERE E.Libro_ISBN = :1 AND P.estado = 'ACTIVO'`
	var count int
	err := config.DB.QueryRow(checkQuery, isbn).Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		return sql.ErrNoRows // Retornar error si hay préstamos activos
	}

	// Marcar ejemplares como no disponibles
	query := `UPDATE Ejemplar SET estado = 'NO_DISPONIBLE' WHERE Libro_ISBN = :1`
	_, err = config.DB.Exec(query, isbn)
	if err != nil {
		return err
	}

	// Registrar en bitácora
	s.bitacoraService.RegistrarAccion(userID, "DELETE", "LIBRO", "Libro marcado como no disponible: "+isbn)

	return nil
}
