package sqliteust

import (
	"database/sql"

	"github.com/pargomx/gecko/gko"

	"monorepo/ust"
)

//  ================================================================  //
//  ========== INSERT ==============================================  //

func (s *Repositorio) InsertTramo(tra ust.Tramo) error {
	const op string = "InsertTramo"
	if tra.HistoriaID == 0 {
		return gko.ErrDatoIndef().Op(op).Msg("HistoriaID sin especificar").Str("pk_indefinida")
	}
	if tra.Posicion == 0 {
		return gko.ErrDatoIndef().Op(op).Msg("Posicion sin especificar").Str("pk_indefinida")
	}
	if tra.Texto == "" {
		return gko.ErrDatoIndef().Op(op).Msg("Texto sin especificar").Str("required_sin_valor")
	}
	_, err := s.db.Exec("INSERT INTO tramos "+
		"(historia_id, posicion, texto, imagen) "+
		"VALUES (?, ?, ?, ?) ",
		tra.HistoriaID, tra.Posicion, tra.Texto, tra.Imagen,
	)
	if err != nil {
		return gko.ErrAlEscribir().Err(err).Op(op)
	}
	return nil
}

//  ================================================================  //
//  ========== UPDATE ==============================================  //

// UpdateTramo valida y sobreescribe el registro en la base de datos.
func (s *Repositorio) UpdateTramo(tra ust.Tramo) error {
	const op string = "UpdateTramo"
	if tra.HistoriaID == 0 {
		return gko.ErrDatoIndef().Op(op).Msg("HistoriaID sin especificar").Str("pk_indefinida")
	}
	if tra.Posicion == 0 {
		return gko.ErrDatoIndef().Op(op).Msg("Posicion sin especificar").Str("pk_indefinida")
	}
	if tra.Texto == "" {
		return gko.ErrDatoIndef().Op(op).Msg("Texto sin especificar").Str("required_sin_valor")
	}
	_, err := s.db.Exec(
		"UPDATE tramos SET "+
			"historia_id=?, posicion=?, texto=?, imagen=? "+
			"WHERE historia_id = ? AND posicion = ?",
		tra.HistoriaID, tra.Posicion, tra.Texto, tra.Imagen,
		tra.HistoriaID, tra.Posicion,
	)
	if err != nil {
		return gko.ErrInesperado().Err(err).Op(op)
	}
	return nil
}

//  ================================================================  //
//  ========== EXISTE ==============================================  //

// Retorna error nil si existe solo un registro con esta clave primaria.
func (s *Repositorio) ExisteTramo(HistoriaID int, Posicion int) error {
	const op string = "ExisteTramo"
	var num int
	err := s.db.QueryRow("SELECT COUNT(historia_id) FROM tramos WHERE historia_id = ? AND posicion = ?",
		HistoriaID, Posicion,
	).Scan(&num)
	if err != nil {
		if err == sql.ErrNoRows {
			return gko.ErrNoEncontrado().Err(ust.ErrTramoNotFound).Op(op)
		}
		return gko.ErrInesperado().Err(err).Op(op)
	}
	if num > 1 {
		return gko.ErrInesperado().Err(nil).Op(op).Str("existen más de un registro para la pk").Ctx("registros", num)
	} else if num == 0 {
		return gko.ErrNoEncontrado().Err(ust.ErrTramoNotFound).Op(op)
	}
	return nil
}

//  ================================================================  //
//  ========== DELETE ==============================================  //

func (s *Repositorio) DeleteTramo(HistoriaID int, Posicion int) error {
	const op string = "DeleteTramo"
	if HistoriaID == 0 {
		return gko.ErrDatoIndef().Op(op).Msg("HistoriaID sin especificar").Str("pk_indefinida")
	}
	if Posicion == 0 {
		return gko.ErrDatoIndef().Op(op).Msg("Posicion sin especificar").Str("pk_indefinida")
	}
	err := s.ExisteTramo(HistoriaID, Posicion)
	if err != nil {
		return gko.Err(err).Op(op)
	}
	_, err = s.db.Exec(
		"DELETE FROM tramos WHERE historia_id = ? AND posicion = ?",
		HistoriaID, Posicion,
	)
	if err != nil {
		return gko.ErrAlEscribir().Err(err).Op(op)
	}
	_, err = s.db.Exec(
		"UPDATE tramos SET posicion = posicion - 1 "+
			"WHERE historia_id = ? AND posicion > ?",
		HistoriaID, Posicion,
	)
	if err != nil {
		return gko.ErrAlEscribir().Err(err).Op(op)
	}
	return nil
}

//  ================================================================  //
//  ========== CONSTANTES ==========================================  //

// Lista de columnas separadas por coma para usar en consulta SELECT
// en conjunto con scanRow o scanRows, ya que las columnas coinciden
// con los campos escaneados.
//
//	historia_id,
//	posicion,
//	texto,
//	imagen
const columnasTramo string = "historia_id, posicion, texto, imagen"

// Origen de los datos de ust.Tramo
//
//	FROM tramos
const fromTramo string = "FROM tramos "

//  ================================================================  //
//  ========== SCAN ================================================  //

// scanRowsTramo escanea cada row en la struct Tramo
// y devuelve un slice con todos los items.
// Siempre se encarga de llamar rows.Close()
func (s *Repositorio) scanRowsTramo(rows *sql.Rows, op string) ([]ust.Tramo, error) {
	defer rows.Close()
	items := []ust.Tramo{}
	for rows.Next() {
		tra := ust.Tramo{}
		err := rows.Scan(
			&tra.HistoriaID, &tra.Posicion, &tra.Texto, &tra.Imagen,
		)
		if err != nil {
			return nil, gko.ErrInesperado().Err(err).Op(op)
		}
		items = append(items, tra)
	}
	return items, nil
}

//  ================================================================  //
//  ========== LIST_BY =============================================  //

func (s *Repositorio) ListTramosByHistoriaID(HistoriaID int) ([]ust.Tramo, error) {
	const op string = "ListTramosByHistoriaID"
	if HistoriaID == 0 {
		return nil, gko.ErrDatoIndef().Op(op).Msg("HistoriaID sin especificar").Str("param_indefinido")
	}
	rows, err := s.db.Query(
		"SELECT "+columnasTramo+" "+fromTramo+
			"WHERE historia_id = ?",
		HistoriaID,
	)
	if err != nil {
		return nil, gko.ErrInesperado().Err(err).Op(op)
	}
	return s.scanRowsTramo(rows, op)
}
