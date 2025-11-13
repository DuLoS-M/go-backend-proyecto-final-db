package services

import (
	"errors"
	"proyecto-bd-final/internal/models"
	"proyecto-bd-final/internal/repository"
	"proyecto-bd-final/pkg/utils"
)

type AuthService struct {
	userRepo *repository.UserRepository
}

func NewAuthService() *AuthService {
	return &AuthService{
		userRepo: repository.NewUserRepository(),
	}
}

// Login autentica un usuario y genera un token JWT
func (s *AuthService) Login(email, password string) (string, *models.Usuario, []string, error) {
	// Buscar usuario por email
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return "", nil, nil, errors.New("credenciales inválidas")
	}

	// Verificar contraseña
	if !utils.CheckPasswordHash(password, user.Contrasenia) {
		return "", nil, nil, errors.New("credenciales inválidas")
	}

	// Obtener roles del usuario
	roles, err := s.userRepo.GetRoles(user.IDUsuario)
	if err != nil {
		return "", nil, nil, err
	}

	// Si no tiene roles, asignar rol por defecto
	if len(roles) == 0 {
		roles = []string{"usuario"}
	}

	// Generar token JWT
	token, err := utils.GenerateToken(user.IDUsuario, user.Correo, roles)
	if err != nil {
		return "", nil, nil, err
	}

	return token, user, roles, nil
}

// Register registra un nuevo usuario
func (s *AuthService) Register(nombre, apellido, email, password string, telefono int) (string, *models.Usuario, []string, error) {
	// Verificar si el email ya existe
	existingUser, _ := s.userRepo.GetByEmail(email)
	if existingUser != nil {
		return "", nil, nil, errors.New("el correo ya está registrado")
	}

	// Hash de la contraseña
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return "", nil, nil, err
	}

	// Crear usuario
	user := &models.Usuario{
		Nombre:      nombre,
		Apellido:    apellido,
		Correo:      email,
		Contrasenia: hashedPassword,
		Telefono:    telefono,
	}

	err = s.userRepo.Create(user)
	if err != nil {
		return "", nil, nil, err
	}

	// Asignar rol por defecto (estudiante)
	// Asumiendo que el rol de estudiante tiene ID 2
	roles := []string{"estudiante"}

	// Generar token
	token, err := utils.GenerateToken(user.IDUsuario, user.Correo, roles)
	if err != nil {
		return "", nil, nil, err
	}

	return token, user, roles, nil
}
