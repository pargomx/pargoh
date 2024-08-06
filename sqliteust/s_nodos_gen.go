package sqliteust

import (
	"database/sql"
	"errors"

	"github.com/pargomx/gecko/gko"

	"monorepo/ust"
)

//  ================================================================  //
//  ========== UPDATE ==============================================  //

// UpdateNodo valida y sobreescribe el registro en la base de datos.
func (s *Repositorio) UpdateNodo(nod ust.Nodo) error {
	const op string = "UpdateNodo"
	if nod.NodoID == 0 {
		return gko.ErrDatoIndef().Op(op).Msg("NodoID sin especificar").Str("pk_indefinida")
	}
	if nod.NodoTbl == "" {
		return gko.ErrDatoIndef().Op(op).Msg("NodoTbl sin especificar").Str("required_sin_valor")
	}
	if nod.PadreTbl == "" {
		return gko.ErrDatoIndef().Op(op).Msg("PadreTbl sin especificar").Str("required_sin_valor")
	}
	_, err := s.db.Exec(
		"UPDATE nodos SET "+
			"nodo_id=?, nodo_tbl=?, padre_id=?, padre_tbl=?, nivel=?, posicion=? "+
			"WHERE nodo_id = ?",
		nod.NodoID, nod.NodoTbl, nod.PadreID, nod.PadreTbl, nod.Nivel, nod.Posicion,
		nod.NodoID,
	)
	if err != nil {
		return gko.ErrInesperado().Err(err).Op(op)
	}
	return nil
}

//  ================================================================  //
//  ========== EXISTE ==============================================  //

// Retorna error nil si existe solo un registro con esta clave primaria.
func (s *Repositorio) ExisteNodo(NodoID int) error {
	const op string = "ExisteNodo"
	var num int
	err := s.db.QueryRow("SELECT COUNT(nodo_id) FROM nodos WHERE nodo_id = ?",
		NodoID,
	).Scan(&num)
	if err != nil {
		if err == sql.ErrNoRows {
			return gko.ErrNoEncontrado().Err(ust.ErrNodoNotFound).Op(op)
		}
		return gko.ErrInesperado().Err(err).Op(op)
	}
	if num > 1 {
		return gko.ErrInesperado().Err(nil).Op(op).Str("existen más de un registro para la pk").Ctx("registros", num)
	} else if num == 0 {
		return gko.ErrNoEncontrado().Err(ust.ErrNodoNotFound).Op(op)
	}
	return nil
}

//  ================================================================  //
//  ========== DELETE ==============================================  //

func (s *Repositorio) DeleteNodo(NodoID int) error {
	const op string = "DeleteNodo"
	if NodoID == 0 {
		return gko.ErrDatoIndef().Op(op).Msg("NodoID sin especificar").Str("pk_indefinida")
	}
	err := s.ExisteNodo(NodoID)
	if err != nil {
		return gko.Err(err).Op(op)
	}
	_, err = s.db.Exec(
		"DELETE FROM nodos WHERE nodo_id = ?",
		NodoID,
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
//	nodo_id,
//	nodo_tbl,
//	padre_id,
//	padre_tbl,
//	nivel,
//	posicion
const columnasNodo string = "nodo_id, nodo_tbl, padre_id, padre_tbl, nivel, posicion"

// Origen de los datos de ust.Nodo
//
//	FROM nodos
const fromNodo string = "FROM nodos "

//  ================================================================  //
//  ========== SCAN ================================================  //

// Utilizar luego de un sql.QueryRow(). No es necesario hacer row.Close()
func (s *Repositorio) scanRowNodo(row *sql.Row, nod *ust.Nodo) error {
	err := row.Scan(
		&nod.NodoID, &nod.NodoTbl, &nod.PadreID, &nod.PadreTbl, &nod.Nivel, &nod.Posicion,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return gko.ErrNoEncontrado().Msg("Nodo no se encuentra")
		}
		return gko.ErrInesperado().Err(err)
	}
	return nil
}

//  ================================================================  //
//  ========== GET =================================================  //

// GetNodo devuelve un Nodo de la DB por su clave primaria.
// Error si no encuentra ninguno, o si encuentra más de uno.
func (s *Repositorio) GetNodo(NodoID int) (*ust.Nodo, error) {
	// const op string = "GetNodo"
	// if NodoID == 0 {
	// 	return nil, gko.ErrDatoIndef().Op(op).Msg("NodoID sin especificar").Str("pk_indefinida")
	// }
	row := s.db.QueryRow(
		"SELECT "+columnasNodo+" "+fromNodo+
			"WHERE nodo_id = ?",
		NodoID,
	)
	nod := &ust.Nodo{}
	err := s.scanRowNodo(row, nod)
	if err != nil {
		return nil, err
	}
	return nod, nil
}

//  ================================================================  //
//  ========== SCAN ================================================  //

// scanRowsNodo escanea cada row en la struct Nodo
// y devuelve un slice con todos los items.
// Siempre se encarga de llamar rows.Close()
func (s *Repositorio) scanRowsNodo(rows *sql.Rows, op string) ([]ust.Nodo, error) {
	defer rows.Close()
	items := []ust.Nodo{}
	for rows.Next() {
		nod := ust.Nodo{}
		err := rows.Scan(
			&nod.NodoID, &nod.NodoTbl, &nod.PadreID, &nod.PadreTbl, &nod.Nivel, &nod.Posicion,
		)
		if err != nil {
			return nil, gko.ErrInesperado().Err(err).Op(op)
		}
		items = append(items, nod)
	}
	return items, nil
}

//  ================================================================  //
//  ========== LIST_BY =============================================  //

func (s *Repositorio) ListNodosByPadreID(PadreID int) ([]ust.Nodo, error) {
	const op string = "ListNodosByPadreID"
	if PadreID == 0 {
		return nil, gko.ErrDatoIndef().Op(op).Msg("PadreID sin especificar").Str("param_indefinido")
	}
	rows, err := s.db.Query(
		"SELECT "+columnasNodo+" "+fromNodo+
			"WHERE padre_id = ?",
		PadreID,
	)
	if err != nil {
		return nil, gko.ErrInesperado().Err(err).Op(op)
	}
	return s.scanRowsNodo(rows, op)
}
