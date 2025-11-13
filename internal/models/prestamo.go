package models

import "time"

type Prestamo struct {
	IDPrestamo              int        `json:"id_prestamo" db:"IDPRESTAMO"`
	FechaPrestamo           time.Time  `json:"fecha_prestamo" db:"FECHAPRESTAMO"`
	FechaDevolucionPrevista time.Time  `json:"fecha_devolucion_prevista" db:"FECHADEVOLUCIONPREVISTA"`
	FechaDevolucionReal     *time.Time `json:"fecha_devolucion_real,omitempty" db:"FECHADEVOLUCIONREAL"`
	Estado                  string     `json:"estado" db:"ESTADO"`
	UsuarioID               int        `json:"usuario_id" db:"USUARIO_IDUSUARIO"`
	DevolucionID            int        `json:"devolucion_id" db:"DEVOLUCION_IDDEVOLUCION"`
}
