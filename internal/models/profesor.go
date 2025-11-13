package models

type Profesor struct {
	CodigoDocencia int    `json:"codigo_docencia" db:"CODIGODOCENCIA"`
	Facultad       string `json:"facultad" db:"FACULTAD"`
	UsuarioID      int    `json:"usuario_id" db:"USUARIO_IDUSUARIO"`
}
