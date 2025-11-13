package models

type Rol struct {
	IDRol     int    `json:"id_rol" db:"IDROL"`
	NombreRol string `json:"nombre_rol" db:"NOMBREROL"`
}

type Permiso struct {
	IDPermiso   int    `json:"id_permiso" db:"IDPERMISO"`
	Descripcion string `json:"descripcion" db:"DESCRIPCION"`
}

type RolPermiso struct {
	IDRolPermiso int `json:"id_rol_permiso" db:"IDROLPERMISO"`
	RolID        int `json:"rol_id" db:"ROLES_IDROL"`
	PermisoID    int `json:"permiso_id" db:"PERMISO_IDPERMISO"`
}
