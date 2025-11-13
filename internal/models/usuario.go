package models

import "time"

type Usuario struct {
	IDUsuario     int       `json:"id_usuario" db:"IDUSUARIO"`
	Nombre        string    `json:"nombre" db:"NOMBRE"`
	Apellido      string    `json:"apellido" db:"APELLIDO"`
	Contrasenia   string    `json:"-" db:"CONTRASENIA"` // No exponer en JSON
	Correo        string    `json:"correo" db:"CORREO"`
	Telefono      int       `json:"telefono" db:"TELEFONO"`
	FechaRegistro time.Time `json:"fecha_registro" db:"FECHAREGISTRO"`
}

type UsuarioRol struct {
	IDUsuarioRol int `json:"id_usuario_rol" db:"IDUSUARIOROL"`
	UsuarioID    int `json:"usuario_id" db:"USUARIO_IDUSUARIO"`
	RolID        int `json:"rol_id" db:"ROLES_IDROL"`
}
