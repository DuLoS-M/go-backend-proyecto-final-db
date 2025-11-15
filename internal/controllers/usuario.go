package controllers

import (
	"net/http"
	"proyecto-bd-final/internal/repository"
	"proyecto-bd-final/internal/services"
	"proyecto-bd-final/pkg/utils"

	"github.com/gin-gonic/gin"
)

var (
	authService = services.NewAuthService()
	userRepo    = repository.NewUserRepository()
	bitacora    = services.NewBitacoraService()
)

// Login maneja el inicio de sesión
func Login(c *gin.Context) {
	var loginData struct {
		Correo      string `json:"correo" binding:"required,email"`
		Contrasenia string `json:"contrasenia" binding:"required"`
	}

	if err := c.ShouldBindJSON(&loginData); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Datos inválidos", err)
		return
	}

	// Autenticar usuario
	token, usuario, roles, err := authService.Login(loginData.Correo, loginData.Contrasenia)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Credenciales inválidas", err)
		return
	}

	// Registrar en bitácora
	bitacora.RegistrarAccion(usuario.IDUsuario, "LOGIN", "Usuario", "Inicio de sesión exitoso")

	// Respuesta exitosa
	utils.SuccessResponse(c, http.StatusOK, "Login exitoso", gin.H{
		"token":   token,
		"usuario": usuario,
		"roles":   roles,
	})
}

// Register maneja el registro de nuevos usuarios
func Register(c *gin.Context) {
	var registerData struct {
		Nombre      string `json:"nombre" binding:"required"`
		Apellido    string `json:"apellido" binding:"required"`
		Correo      string `json:"correo" binding:"required,email"`
		Contrasenia string `json:"contrasenia" binding:"required,min=6"`
		Telefono    int    `json:"telefono" binding:"required"`
	}

	if err := c.ShouldBindJSON(&registerData); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Datos inválidos", err)
		return
	}

	// Registrar usuario
	token, usuario, roles, err := authService.Register(
		registerData.Nombre,
		registerData.Apellido,
		registerData.Correo,
		registerData.Contrasenia,
		registerData.Telefono,
	)

	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Error al registrar usuario", err)
		return
	}

	// Registrar en bitácora
	bitacora.RegistrarAccion(usuario.IDUsuario, "REGISTRO", "Usuario", "Usuario registrado exitosamente")

	// Respuesta exitosa
	utils.SuccessResponse(c, http.StatusCreated, "Usuario registrado exitosamente", gin.H{
		"token":   token,
		"usuario": usuario,
		"roles":   roles,
	})
}

// GetProfile obtiene el perfil del usuario actual
func GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Usuario no autenticado", nil)
		return
	}

	// Convertir a int de forma segura
	var userIDInt int
	switch v := userID.(type) {
	case int:
		userIDInt = v
	case float64:
		userIDInt = int(v)
	default:
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error en el formato del ID de usuario", nil)
		return
	}

	usuario, err := userRepo.GetByID(userIDInt)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Usuario no encontrado", err)
		return
	}

	roles, err := userRepo.GetRoles(userIDInt)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error al obtener roles", err)
		return
	}

	// Asegurarse de que roles no sea nil
	if roles == nil {
		roles = []string{}
	}

	utils.SuccessResponse(c, http.StatusOK, "Perfil obtenido", gin.H{
		"usuario": usuario,
		"roles":   roles,
	})
}

// UpdateProfile actualiza el perfil del usuario
func UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Usuario no autenticado", nil)
		return
	}

	var updateData struct {
		Nombre   string `json:"nombre" binding:"required"`
		Apellido string `json:"apellido" binding:"required"`
		Telefono int    `json:"telefono" binding:"required"`
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Datos inválidos", err)
		return
	}

	// Obtener usuario actual
	usuario, err := userRepo.GetByID(userID.(int))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Usuario no encontrado", err)
		return
	}

	// Actualizar datos
	usuario.Nombre = updateData.Nombre
	usuario.Apellido = updateData.Apellido
	usuario.Telefono = updateData.Telefono

	err = userRepo.Update(usuario)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error al actualizar perfil", err)
		return
	}

	// Registrar en bitácora
	bitacora.RegistrarAccion(userID.(int), "UPDATE", "Usuario", "Perfil actualizado")

	utils.SuccessResponse(c, http.StatusOK, "Perfil actualizado exitosamente", usuario)
}

// GetAllUsers obtiene todos los usuarios (admin)
func GetAllUsers(c *gin.Context) {
	usuarios, err := userRepo.GetAll()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error al obtener usuarios", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Usuarios obtenidos", usuarios)
}
