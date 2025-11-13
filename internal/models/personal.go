package models

type Personal struct {
	CodigoEmpleado int    `json:"codigo_empleado" db:"CODIGOEMPLEADO"`
	Puesto         string `json:"puesto" db:"PUESTO"`
	UsuarioID      int    `json:"usuario_id" db:"USUARIO_IDUSUARIO"`
}
