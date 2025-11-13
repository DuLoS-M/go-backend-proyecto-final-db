package controllers

import (
	"net/http"
	"proyecto-bd-final/internal/models"
	"proyecto-bd-final/internal/services"
	"proyecto-bd-final/pkg/utils"

	"github.com/gin-gonic/gin"
)

var bookService = services.NewBookService()

// GetBooks obtiene todos los libros
func GetBooks(c *gin.Context) {
	// Verificar si hay un parámetro de búsqueda
	searchTerm := c.Query("q")

	var libros []*models.Libro
	var err error

	if searchTerm != "" {
		libros, err = bookService.SearchBooks(searchTerm)
	} else {
		libros, err = bookService.GetAll()
	}

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error al obtener libros", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Libros obtenidos", libros)
}

// GetBookByISBN obtiene un libro por ISBN
func GetBookByISBN(c *gin.Context) {
	isbn := c.Param("isbn")

	libro, err := bookService.GetByISBN(isbn)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Libro no encontrado", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Libro obtenido", libro)
}

// CreateBook crea un nuevo libro (admin)
func CreateBook(c *gin.Context) {
	var libroData struct {
		ISBN            string `json:"isbn" binding:"required"`
		Titulo          string `json:"titulo" binding:"required"`
		AnioPublicacion int    `json:"anio_publicacion" binding:"required"`
		Cantidad        int    `json:"cantidad" binding:"required"`
		EditorialID     int    `json:"editorial_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&libroData); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Datos inválidos", err)
		return
	}

	libro := &models.Libro{
		ISBN:            libroData.ISBN,
		Titulo:          libroData.Titulo,
		AnioPublicacion: libroData.AnioPublicacion,
		Cantidad:        libroData.Cantidad,
		EditorialID:     libroData.EditorialID,
	}

	userID, _ := c.Get("user_id")
	err := bookService.Create(libro, userID.(int))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error al crear libro", err)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Libro creado exitosamente", libro)
}

// UpdateBook actualiza un libro (admin)
func UpdateBook(c *gin.Context) {
	isbn := c.Param("isbn")

	var libroData struct {
		Titulo          string `json:"titulo" binding:"required"`
		AnioPublicacion int    `json:"anio_publicacion" binding:"required"`
		Cantidad        int    `json:"cantidad" binding:"required"`
		EditorialID     int    `json:"editorial_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&libroData); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Datos inválidos", err)
		return
	}

	libro := &models.Libro{
		ISBN:            isbn,
		Titulo:          libroData.Titulo,
		AnioPublicacion: libroData.AnioPublicacion,
		Cantidad:        libroData.Cantidad,
		EditorialID:     libroData.EditorialID,
	}

	userID, _ := c.Get("user_id")
	err := bookService.Update(libro, userID.(int))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error al actualizar libro", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Libro actualizado exitosamente", libro)
}

// DeleteBook elimina un libro (admin)
func DeleteBook(c *gin.Context) {
	isbn := c.Param("isbn")

	userID, _ := c.Get("user_id")
	err := bookService.Delete(isbn, userID.(int))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error al eliminar libro", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Libro eliminado exitosamente", nil)
}
