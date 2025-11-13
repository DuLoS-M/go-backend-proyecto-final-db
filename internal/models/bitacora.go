package models

import "time"

type Bitacora struct {
	IDBitacora int       `json:"id_bitacora" db:"IDBITACORA"`
	Accion     string    `json:"accion" db:"ACCION"`
	FechaHora  time.Time `json:"fecha_hora" db:"FECHAHORA"`
	Detalle    string    `json:"detalle" db:"DETALLE"`
	Entidad    string    `json:"entidad" db:"ENTIDAD"`
	UsuarioID  int       `json:"usuario_id" db:"USUARIO_IDUSUARIO"`
}
