package sqliteust

import (
	"database/sql"
	"errors"

	"github.com/pargomx/gecko/gko"

	"monorepo/ust"
)

//  ================================================================  //
//  ========== INSERT ==============================================  //

func (s *Repositorio) InsertTramo(tra ust.Tramo) error {
	const op string = "InsertTramo"
	if tra.HistoriaID == 0 {
		return gko.ErrDatoIndef.Str("pk_indefinida").Op(op).Msg("HistoriaID sin especificar")
	}
	if tra.Posicion == 0 {
		return gko.ErrDatoIndef.Str("pk_indefinida").Op(op).Msg("Posicion sin especificar")
	}
	if tra.Texto == "" {
		return gko.ErrDatoIndef.Str("required_sin_valor").Op(op).Msg("Texto sin especificar")
	}
	_, err := s.db.Exec("INSERT INTO tramos "+
		"(historia_id, posicion, texto, imagen) "+
		"VALUES (?, ?, ?, ?) ",
		tra.HistoriaID, tra.Posicion, tra.Texto, tra.Imagen,
	)
	if err != nil {
		return gko.ErrAlEscribir.Err(err).Op(op)
	}
	return nil
}

//  ================================================================  //
//  ========== UPDATE ==============================================  //

// UpdateTramo valida y sobreescribe el registro en la base de datos.
func (s *Repositorio) UpdateTramo(tra ust.Tramo) error {
	const op string = "UpdateTramo"
	if tra.HistoriaID == 0 {
		return gko.ErrDatoIndef.Str("pk_indefinida").Op(op).Msg("HistoriaID sin especificar")
	}
	if tra.Posicion == 0 {
		return gko.ErrDatoIndef.Str("pk_indefinida").Op(op).Msg("Posicion sin especificar")
	}
	if tra.Texto == "" {
		return gko.ErrDatoIndef.Str("required_sin_valor").Op(op).Msg("Texto sin especificar")
	}
	_, err := s.db.Exec(
		"UPDATE tramos SET "+
			"historia_id=?, posicion=?, texto=?, imagen=? "+
			"WHERE historia_id = ? AND posicion = ?",
		tra.HistoriaID, tra.Posicion, tra.Texto, tra.Imagen,
		tra.HistoriaID, tra.Posicion,
	)
	if err != nil {
		return gko.ErrInesperado.Err(err).Op(op)
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
			return gko.ErrNoEncontrado.Msg("Tramo no encontrado").Op(op)
		}
		return gko.ErrInesperado.Err(err).Op(op)
	}
	if num > 1 {
		return gko.ErrInesperado.Err(nil).Op(op).Str("existen más de un registro para la pk").Ctx("registros", num)
	} else if num == 0 {
		return gko.ErrNoEncontrado.Msg("Tramo no encontrado").Op(op)
	}
	return nil
}

//  ================================================================  //
//  ========== DELETE ==============================================  //

func (s *Repositorio) DeleteTramo(HistoriaID int, Posicion int) error {
	const op string = "DeleteTramo"
	if HistoriaID == 0 {
		return gko.ErrDatoIndef.Str("pk_indefinida").Op(op).Msg("HistoriaID sin especificar")
	}
	if Posicion == 0 {
		return gko.ErrDatoIndef.Str("pk_indefinida").Op(op).Msg("Posicion sin especificar")
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
		return gko.ErrAlEscribir.Err(err).Op(op).Ctx("historia_id", HistoriaID).Ctx("Pos", Posicion)
	}
	_, err = s.db.Exec(
		"UPDATE tramos SET posicion = posicion - 1 WHERE historia_id = ? AND posicion > ?",
		HistoriaID, Posicion,
	)
	if err != nil {
		return gko.ErrAlEscribir.Err(err).Op(op).Ctx("historia_id", HistoriaID).Ctx("Pos", Posicion)
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

// Utilizar luego de un sql.QueryRow(). No es necesario hacer row.Close()
func (s *Repositorio) scanRowTramo(row *sql.Row, tra *ust.Tramo) error {
	err := row.Scan(
		&tra.HistoriaID, &tra.Posicion, &tra.Texto, &tra.Imagen,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return gko.ErrNoEncontrado.Msg("Tramo no encontrado")
		}
		return gko.ErrInesperado.Err(err)
	}
	return nil
}

//  ================================================================  //
//  ========== GET =================================================  //

// GetTramo devuelve un Tramo de la DB por su clave primaria.
// Error si no encuentra ninguno, o si encuentra más de uno.
func (s *Repositorio) GetTramo(HistoriaID int, Posicion int) (*ust.Tramo, error) {
	const op string = "GetTramo"
	if HistoriaID == 0 {
		return nil, gko.ErrDatoIndef.Str("pk_indefinida").Op(op).Msg("HistoriaID sin especificar")
	}
	if Posicion == 0 {
		return nil, gko.ErrDatoIndef.Str("pk_indefinida").Op(op).Msg("Posicion sin especificar")
	}
	row := s.db.QueryRow(
		"SELECT "+columnasTramo+" "+fromTramo+
			"WHERE historia_id = ? AND posicion = ?",
		HistoriaID, Posicion,
	)
	tra := &ust.Tramo{}
	err := s.scanRowTramo(row, tra)
	if err != nil {
		return nil, err
	}
	return tra, nil
}

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
			return nil, gko.ErrInesperado.Err(err).Op(op)
		}
		items = append(items, tra)
	}
	return items, nil
}

//  ================================================================  //
//  ========== LIST_BY HISTORIA_ID =================================  //

func (s *Repositorio) ListTramosByHistoriaID(HistoriaID int) ([]ust.Tramo, error) {
	const op string = "ListTramosByHistoriaID"
	if HistoriaID == 0 {
		return nil, gko.ErrDatoIndef.Str("param_indefinido").Op(op).Msg("HistoriaID sin especificar")
	}
	rows, err := s.db.Query(
		"SELECT "+columnasTramo+" "+fromTramo+
			"WHERE historia_id = ?",
		HistoriaID,
	)
	if err != nil {
		return nil, gko.ErrInesperado.Err(err).Op(op)
	}
	return s.scanRowsTramo(rows, op)
}
