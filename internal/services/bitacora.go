package services

import (
	"database/sql"
	"proyecto-bd-final/internal/config"
	"proyecto-bd-final/internal/models"
	"time"
)

type BitacoraService struct{}

func NewBitacoraService() *BitacoraService {
	return &BitacoraService{}
}

// RegistrarAccion registra una acción en la bitácora
func (s *BitacoraService) RegistrarAccion(usuarioID int, accion, entidad, detalle string) error {
	query := `INSERT INTO Bitacora (IDBITACORA, ACCION, FECHAHORA, DETALLE, ENTIDAD, USUARIO_IDUSUARIO) 
			  VALUES (BITACORA_SEQ.NEXTVAL, :1, :2, :3, :4, :5)`

	_, err := config.DB.Exec(query, accion, time.Now(), detalle, entidad, usuarioID)
	return err
}

// ObtenerBitacora obtiene registros de la bitácora con filtros opcionales
func (s *BitacoraService) ObtenerBitacora(limite int, entidad string) ([]*models.Bitacora, error) {
	var query string
	var rows *sql.Rows
	var err error

	if entidad != "" {
		query = `SELECT B.IDBITACORA, B.ACCION, B.FECHAHORA, B.DETALLE, B.ENTIDAD, B.USUARIO_IDUSUARIO 
				 FROM Bitacora B 
				 WHERE B.ENTIDAD = :1 
				 ORDER BY B.FECHAHORA DESC 
				 FETCH FIRST :2 ROWS ONLY`
		rows, err = config.DB.Query(query, entidad, limite)
	} else {
		query = `SELECT B.IDBITACORA, B.ACCION, B.FECHAHORA, B.DETALLE, B.ENTIDAD, B.USUARIO_IDUSUARIO 
				 FROM Bitacora B 
				 ORDER BY B.FECHAHORA DESC 
				 FETCH FIRST :1 ROWS ONLY`
		rows, err = config.DB.Query(query, limite)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var registros []*models.Bitacora
	for rows.Next() {
		var registro models.Bitacora
		if err := rows.Scan(
			&registro.IDBitacora,
			&registro.Accion,
			&registro.FechaHora,
			&registro.Detalle,
			&registro.Entidad,
			&registro.UsuarioID,
		); err != nil {
			return nil, err
		}
		registros = append(registros, &registro)
	}

	return registros, nil
}
