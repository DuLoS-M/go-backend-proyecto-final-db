package models

type Libro struct {
	ISBN            string   `json:"isbn" db:"ISBN"`
	Titulo          string   `json:"titulo" db:"TITULO"`
	AnioPublicacion int      `json:"anioPublicacion" db:"ANIOEDICION"`
	EditorialID     int      `json:"editorialId" db:"EDITORIAL_IDEDITORIAL"`
	EditorialNombre string   `json:"editorialNombre,omitempty"`
	Autores         []string `json:"autores,omitempty"`
	Cantidad        int      `json:"cantidad,omitempty"`
	Disponible      bool     `json:"disponible,omitempty"`
}

type Editorial struct {
	IDEditorial int    `json:"idEditorial" db:"IDEDITORIAL"`
	Nombre      string `json:"nombre" db:"NOMBRE"`
	Pais        string `json:"pais" db:"PAIS"`
}

type Autor struct {
	IDAutor      int    `json:"idAutor" db:"IDAUTOR"`
	Nombre       string `json:"nombre" db:"NOMBRE"`
	Apellido     string `json:"apellido" db:"APELLIDO"`
	Nacionalidad string `json:"nacionalidad" db:"NACIONALIDAD"`
}

type LibroAutor struct {
	IDLibroAutor int    `json:"idLibroAutor" db:"IDLIBROAUTOR"`
	TipoAutor    string `json:"tipoAutor" db:"TIPOAUTOR"`
	AutorID      int    `json:"autorId" db:"AUTOR_IDAUTOR"`
	LibroISBN    int    `json:"libroIsbn" db:"LIBRO_ISBN"`
}
type Ejemplar struct {
	IDEjemplar   int    `json:"id_ejemplar" db:"IDEJEMPLAR"`
	Localizacion string `json:"localizacion" db:"LOCALIZACION"`
	Estado       string `json:"estado" db:"ESTADO"`
	LibroISBN    string `json:"libro_isbn" db:"Libro_ISBN"`
}
