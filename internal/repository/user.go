package repository

import (
	"database/sql"
	"errors"
	"proyecto-bd-final/internal/config"
	"proyecto-bd-final/internal/models"
)

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

// GetByEmail busca un usuario por correo electrónico
func (r *UserRepository) GetByEmail(email string) (*models.Usuario, error) {
	var user models.Usuario
	query := `SELECT IDUSUARIO, NOMBRE, APELLIDO, CONTRASENIA, CORREO, TELEFONO, FECHAREGISTRO 
			  FROM Usuario WHERE CORREO = :1`

	err := config.DB.QueryRow(query, email).Scan(
		&user.IDUsuario,
		&user.Nombre,
		&user.Apellido,
		&user.Contrasenia,
		&user.Correo,
		&user.Telefono,
		&user.FechaRegistro,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("usuario no encontrado")
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetByID busca un usuario por ID
func (r *UserRepository) GetByID(id int) (*models.Usuario, error) {
	var user models.Usuario
	query := `SELECT IDUSUARIO, NOMBRE, APELLIDO, CONTRASENIA, CORREO, TELEFONO, FECHAREGISTRO 
			  FROM Usuario WHERE IDUSUARIO = :1`

	err := config.DB.QueryRow(query, id).Scan(
		&user.IDUsuario,
		&user.Nombre,
		&user.Apellido,
		&user.Contrasenia,
		&user.Correo,
		&user.Telefono,
		&user.FechaRegistro,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("usuario no encontrado")
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Create crea un nuevo usuario
func (r *UserRepository) Create(user *models.Usuario) error {
	query := `INSERT INTO Usuario (IDUSUARIO, NOMBRE, APELLIDO, CONTRASENIA, CORREO, TELEFONO, FECHAREGISTRO) 
			  VALUES (USUARIO_SEQ.NEXTVAL, :1, :2, :3, :4, :5, SYSDATE) 
			  RETURNING IDUSUARIO INTO :6`

	_, err := config.DB.Exec(query,
		user.Nombre,
		user.Apellido,
		user.Contrasenia,
		user.Correo,
		user.Telefono,
		sql.Out{Dest: &user.IDUsuario},
	)

	return err
}

// Update actualiza un usuario existente
func (r *UserRepository) Update(user *models.Usuario) error {
	query := `UPDATE Usuario 
			  SET NOMBRE = :1, APELLIDO = :2, CORREO = :3, TELEFONO = :4 
			  WHERE IDUSUARIO = :5`

	_, err := config.DB.Exec(query,
		user.Nombre,
		user.Apellido,
		user.Correo,
		user.Telefono,
		user.IDUsuario,
	)

	return err
}

// GetRoles obtiene los roles de un usuario
func (r *UserRepository) GetRoles(userID int) ([]string, error) {
	query := `SELECT R.NOMBREROL 
			  FROM Roles R 
			  INNER JOIN UsuarioRol UR ON R.IDROL = UR.ROLES_IDROL 
			  WHERE UR.USUARIO_IDUSUARIO = :1`

	rows, err := config.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []string
	for rows.Next() {
		var role string
		if err := rows.Scan(&role); err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}

	return roles, nil
}

// AssignRole asigna un rol a un usuario
func (r *UserRepository) AssignRole(userID, roleID int) error {
	query := `INSERT INTO UsuarioRol (IDUSUARIOROL, USUARIO_IDUSUARIO, ROLES_IDROL) 
			  VALUES (USUARIOROL_SEQ.NEXTVAL, :1, :2)`

	_, err := config.DB.Exec(query, userID, roleID)
	return err
}

// GetAll obtiene todos los usuarios (para admin)
func (r *UserRepository) GetAll() ([]*models.Usuario, error) {
	query := `SELECT IDUSUARIO, NOMBRE, APELLIDO, CORREO, TELEFONO, FECHAREGISTRO 
			  FROM Usuario ORDER BY FECHAREGISTRO DESC`

	rows, err := config.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.Usuario
	for rows.Next() {
		var user models.Usuario
		if err := rows.Scan(
			&user.IDUsuario,
			&user.Nombre,
			&user.Apellido,
			&user.Correo,
			&user.Telefono,
			&user.FechaRegistro,
		); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	return users, nil
}
