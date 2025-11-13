package models

type Estudiante struct {
	Carnet    int    `json:"carnet" db:"CARNET"`
	Carrera   string `json:"carrera" db:"CARRERA"`
	Semestre  int    `json:"semestre" db:"SEMESTRE"`
	UsuarioID int    `json:"usuario_id" db:"USUARIO_IDUSUARIO"`
}
