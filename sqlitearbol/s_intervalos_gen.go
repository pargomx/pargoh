package sqlitearbol

import (
	"database/sql"
	"errors"

	"github.com/pargomx/gecko/gko"

	"monorepo/arbol"
)

//  ================================================================  //
//  ========== INSERT ==============================================  //

func (s *Repositorio) InsertIntervalo(itv arbol.Intervalo) error {
	const op string = "InsertIntervalo"
	if itv.NodoID == 0 {
		return gko.ErrDatoIndef.Str("pk_indefinida").Op(op).Msg("NodoID sin especificar")
	}
	if itv.TsIni == "" {
		return gko.ErrDatoIndef.Str("pk_indefinida").Op(op).Msg("TsIni sin especificar")
	}
	_, err := s.db.Exec("INSERT INTO intervalos "+
		"(nodo_id, ts_ini, ts_fin) "+
		"VALUES (?, ?, ?) ",
		itv.NodoID, itv.TsIni, itv.TsFin,
	)
	if err != nil {
		return gko.ErrAlEscribir.Err(err).Op(op)
	}
	return nil
}

//  ================================================================  //
//  ========== UPDATE ==============================================  //

// UpdateIntervalo valida y sobreescribe el registro en la base de datos.
func (s *Repositorio) UpdateIntervalo(NodoID int, TsIni string, itv arbol.Intervalo) error {
	const op string = "UpdateIntervalo"
	if itv.NodoID == 0 {
		return gko.ErrDatoIndef.Str("pk_indefinida").Op(op).Msg("NodoID sin especificar")
	}
	if itv.TsIni == "" {
		return gko.ErrDatoIndef.Str("pk_indefinida").Op(op).Msg("TsIni sin especificar")
	}
	_, err := s.db.Exec(
		"UPDATE intervalos SET "+
			"nodo_id=?, ts_ini=?, ts_fin=? "+
			"WHERE nodo_id = ? AND ts_ini = ?",
		itv.NodoID, itv.TsIni, itv.TsFin,
		NodoID, TsIni,
	)
	if err != nil {
		return gko.ErrInesperado.Err(err).Op(op)
	}
	return nil
}

//  ================================================================  //
//  ========== EXISTE ==============================================  //

// Retorna error nil si existe solo un registro con esta clave primaria.
func (s *Repositorio) ExisteIntervalo(NodoID int, TsIni string) error {
	const op string = "ExisteIntervalo"
	var num int
	err := s.db.QueryRow("SELECT COUNT(nodo_id) FROM intervalos WHERE nodo_id = ? AND ts_ini = ?",
		NodoID, TsIni,
	).Scan(&num)
	if err != nil {
		if err == sql.ErrNoRows {
			return gko.ErrNoEncontrado.Msg("Intervalo no encontrado").Op(op)
		}
		return gko.ErrInesperado.Err(err).Op(op)
	}
	if num > 1 {
		return gko.ErrInesperado.Err(nil).Op(op).Str("existen más de un registro para la pk").Ctx("registros", num)
	} else if num == 0 {
		return gko.ErrNoEncontrado.Msg("Intervalo no encontrado").Op(op)
	}
	return nil
}

//  ================================================================  //
//  ========== DELETE ==============================================  //

func (s *Repositorio) DeleteIntervalo(NodoID int, TsIni string) error {
	const op string = "DeleteIntervalo"
	if NodoID == 0 {
		return gko.ErrDatoIndef.Str("pk_indefinida").Op(op).Msg("NodoID sin especificar")
	}
	if TsIni == "" {
		return gko.ErrDatoIndef.Str("pk_indefinida").Op(op).Msg("TsIni sin especificar")
	}
	err := s.ExisteIntervalo(NodoID, TsIni)
	if err != nil {
		return gko.Err(err).Op(op)
	}
	_, err = s.db.Exec(
		"DELETE FROM intervalos WHERE nodo_id = ? AND ts_ini = ?",
		NodoID, TsIni,
	)
	if err != nil {
		return gko.ErrAlEscribir.Err(err).Op(op)
	}
	return nil
}

//  ================================================================  //
//  ========== CONSTANTES ==========================================  //

// Lista de columnas separadas por coma para usar en consulta SELECT
// en conjunto con scanRow o scanRows, ya que las columnas coinciden
// con los campos escaneados.
//
//	nodo_id,
//	ts_ini,
//	ts_fin
const columnasIntervalo string = "nodo_id, ts_ini, ts_fin"

// Origen de los datos de arbol.Intervalo
//
//	FROM intervalos
const fromIntervalo string = "FROM intervalos "

//  ================================================================  //
//  ========== SCAN ================================================  //

// Utilizar luego de un sql.QueryRow(). No es necesario hacer row.Close()
func (s *Repositorio) scanRowIntervalo(row *sql.Row, itv *arbol.Intervalo) error {
	err := row.Scan(
		&itv.NodoID, &itv.TsIni, &itv.TsFin,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return gko.ErrNoEncontrado.Msg("Intervalo no encontrado")
		}
		return gko.ErrInesperado.Err(err)
	}
	return nil
}

//  ================================================================  //
//  ========== GET =================================================  //

// GetIntervalo devuelve un Intervalo de la DB por su clave primaria.
// Error si no encuentra ninguno, o si encuentra más de uno.
func (s *Repositorio) GetIntervalo(NodoID int, TsIni string) (*arbol.Intervalo, error) {
	const op string = "GetIntervalo"
	if NodoID == 0 {
		return nil, gko.ErrDatoIndef.Str("pk_indefinida").Op(op).Msg("NodoID sin especificar")
	}
	if TsIni == "" {
		return nil, gko.ErrDatoIndef.Str("pk_indefinida").Op(op).Msg("TsIni sin especificar")
	}
	row := s.db.QueryRow(
		"SELECT "+columnasIntervalo+" "+fromIntervalo+
			"WHERE nodo_id = ? AND ts_ini = ?",
		NodoID, TsIni,
	)
	itv := &arbol.Intervalo{}
	err := s.scanRowIntervalo(row, itv)
	if err != nil {
		return nil, err
	}
	return itv, nil
}

//  ================================================================  //
//  ========== SCAN ================================================  //

// scanRowsIntervalo escanea cada row en la struct Intervalo
// y devuelve un slice con todos los items.
// Siempre se encarga de llamar rows.Close()
func (s *Repositorio) scanRowsIntervalo(rows *sql.Rows, op string) ([]arbol.Intervalo, error) {
	defer rows.Close()
	items := []arbol.Intervalo{}
	for rows.Next() {
		itv := arbol.Intervalo{}
		err := rows.Scan(
			&itv.NodoID, &itv.TsIni, &itv.TsFin,
		)
		if err != nil {
			return nil, gko.ErrInesperado.Err(err).Op(op)
		}
		items = append(items, itv)
	}
	return items, nil
}

//  ================================================================  //
//  ========== LIST_BY NODO_ID =====================================  //

func (s *Repositorio) ListIntervalosByNodoID(NodoID int) ([]arbol.Intervalo, error) {
	const op string = "ListIntervalosByNodoID"
	if NodoID == 0 {
		return nil, gko.ErrDatoIndef.Str("param_indefinido").Op(op).Msg("NodoID sin especificar")
	}
	rows, err := s.db.Query(
		"SELECT "+columnasIntervalo+" "+fromIntervalo+
			"WHERE nodo_id = ?",
		NodoID,
	)
	if err != nil {
		return nil, gko.ErrInesperado.Err(err).Op(op)
	}
	return s.scanRowsIntervalo(rows, op)
}
