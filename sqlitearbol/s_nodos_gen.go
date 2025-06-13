package sqlitearbol

import (
	"database/sql"
	"errors"

	"github.com/pargomx/gecko/gko"

	"monorepo/arbol"
)

//  ================================================================  //
//  ========== INSERT ==============================================  //

func (s *Repositorio) InsertNodo(nod arbol.Nodo) error {
	const op string = "InsertNodo"
	if nod.NodoID == 0 {
		return gko.ErrDatoIndef.Str("pk_indefinida").Op(op).Msg("NodoID sin especificar")
	}
	if nod.Tipo == "" {
		return gko.ErrDatoIndef.Str("required_sin_valor").Op(op).Msg("Tipo sin especificar")
	}
	if nod.Título == "" {
		return gko.ErrDatoIndef.Str("required_sin_valor").Op(op).Msg("Título sin especificar")
	}
	_, err := s.db.Exec("INSERT INTO nodos "+
		"(nodo_id, padre_id, tipo, posicion, titulo, descripcion, objetivo, notas, color, imagen, prioridad, estatus, segundos, centavos) "+
		"VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) ",
		nod.NodoID, nod.PadreID, nod.Tipo, nod.Posicion, nod.Título, nod.Descripcion, nod.Objetivo, nod.Notas, nod.Color, nod.Imagen, nod.Prioridad, nod.Estatus, nod.Segundos, nod.Centavos,
	)
	if err != nil {
		return gko.ErrAlEscribir.Err(err).Op(op)
	}
	return nil
}

//  ================================================================  //
//  ========== UPDATE ==============================================  //

// UpdateNodo valida y sobreescribe el registro en la base de datos.
func (s *Repositorio) UpdateNodo(NodoID int, nod arbol.Nodo) error {
	const op string = "UpdateNodo"
	if nod.NodoID == 0 {
		return gko.ErrDatoIndef.Str("pk_indefinida").Op(op).Msg("NodoID sin especificar")
	}
	if nod.Tipo == "" {
		return gko.ErrDatoIndef.Str("required_sin_valor").Op(op).Msg("Tipo sin especificar")
	}
	if nod.Título == "" {
		return gko.ErrDatoIndef.Str("required_sin_valor").Op(op).Msg("Título sin especificar")
	}
	_, err := s.db.Exec(
		"UPDATE nodos SET "+
			"nodo_id=?, padre_id=?, tipo=?, posicion=?, titulo=?, descripcion=?, objetivo=?, notas=?, color=?, imagen=?, prioridad=?, estatus=?, segundos=?, centavos=? "+
			"WHERE nodo_id = ?",
		nod.NodoID, nod.PadreID, nod.Tipo, nod.Posicion, nod.Título, nod.Descripcion, nod.Objetivo, nod.Notas, nod.Color, nod.Imagen, nod.Prioridad, nod.Estatus, nod.Segundos, nod.Centavos,
		NodoID,
	)
	if err != nil {
		return gko.ErrInesperado.Err(err).Op(op)
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
			return gko.ErrNoEncontrado.Msg("Nodo no encontrado").Op(op)
		}
		return gko.ErrInesperado.Err(err).Op(op)
	}
	if num > 1 {
		return gko.ErrInesperado.Err(nil).Op(op).Str("existen más de un registro para la pk").Ctx("registros", num)
	} else if num == 0 {
		return gko.ErrNoEncontrado.Msg("Nodo no encontrado").Op(op)
	}
	return nil
}

//  ================================================================  //
//  ========== DELETE ==============================================  //

func (s *Repositorio) DeleteNodo(NodoID int) error {
	const op string = "DeleteNodo"
	if NodoID == 0 {
		return gko.ErrDatoIndef.Str("pk_indefinida").Op(op).Msg("NodoID sin especificar")
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
//	padre_id,
//	tipo,
//	posicion,
//	titulo,
//	descripcion,
//	objetivo,
//	notas,
//	color,
//	imagen,
//	prioridad,
//	estatus,
//	segundos,
//	centavos
const columnasNodo string = "nodo_id, padre_id, tipo, posicion, titulo, descripcion, objetivo, notas, color, imagen, prioridad, estatus, segundos, centavos"

// Origen de los datos de arbol.Nodo
//
//	FROM nodos
const fromNodo string = "FROM nodos "

//  ================================================================  //
//  ========== SCAN ================================================  //

// Utilizar luego de un sql.QueryRow(). No es necesario hacer row.Close()
func (s *Repositorio) scanRowNodo(row *sql.Row, nod *arbol.Nodo) error {
	err := row.Scan(
		&nod.NodoID, &nod.PadreID, &nod.Tipo, &nod.Posicion, &nod.Título, &nod.Descripcion, &nod.Objetivo, &nod.Notas, &nod.Color, &nod.Imagen, &nod.Prioridad, &nod.Estatus, &nod.Segundos, &nod.Centavos,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return gko.ErrNoEncontrado.Msg("Nodo no encontrado")
		}
		return gko.ErrInesperado.Err(err)
	}
	return nil
}

//  ================================================================  //
//  ========== GET =================================================  //

// GetNodo devuelve un Nodo de la DB por su clave primaria.
// Error si no encuentra ninguno, o si encuentra más de uno.
func (s *Repositorio) GetNodo(NodoID int) (*arbol.Nodo, error) {
	const op string = "GetNodo"
	if NodoID == 0 {
		return nil, gko.ErrDatoIndef.Str("pk_indefinida").Op(op).Msg("NodoID sin especificar")
	}
	row := s.db.QueryRow(
		"SELECT "+columnasNodo+" "+fromNodo+
			"WHERE nodo_id = ?",
		NodoID,
	)
	nod := &arbol.Nodo{}
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
func (s *Repositorio) scanRowsNodo(rows *sql.Rows, op string) ([]arbol.Nodo, error) {
	defer rows.Close()
	items := []arbol.Nodo{}
	for rows.Next() {
		nod := arbol.Nodo{}
		err := rows.Scan(
			&nod.NodoID, &nod.PadreID, &nod.Tipo, &nod.Posicion, &nod.Título, &nod.Descripcion, &nod.Objetivo, &nod.Notas, &nod.Color, &nod.Imagen, &nod.Prioridad, &nod.Estatus, &nod.Segundos, &nod.Centavos,
		)
		if err != nil {
			return nil, gko.ErrInesperado.Err(err).Op(op)
		}
		items = append(items, nod)
	}
	return items, nil
}

//  ================================================================  //
//  ========== LIST_BY PADRE_ID ====================================  //

func (s *Repositorio) ListNodosByPadreID(PadreID int) ([]arbol.Nodo, error) {
	const op string = "ListNodosByPadreID"
	if PadreID == 0 {
		return nil, gko.ErrDatoIndef.Str("param_indefinido").Op(op).Msg("PadreID sin especificar")
	}
	rows, err := s.db.Query(
		"SELECT "+columnasNodo+" "+fromNodo+
			"WHERE padre_id = ?",
		PadreID,
	)
	if err != nil {
		return nil, gko.ErrInesperado.Err(err).Op(op)
	}
	return s.scanRowsNodo(rows, op)
}

//  ================================================================  //
//  ========== LIST_BY PADRE_ID TIPO ===============================  //

func (s *Repositorio) ListNodosByPadreIDTipo(PadreID int, Tipo string) ([]arbol.Nodo, error) {
	const op string = "ListNodosByPadreIDTipo"
	if PadreID == 0 {
		return nil, gko.ErrDatoIndef.Str("param_indefinido").Op(op).Msg("PadreID sin especificar")
	}
	if Tipo == "" {
		return nil, gko.ErrDatoIndef.Str("param_indefinido").Op(op).Msg("Tipo sin especificar")
	}
	rows, err := s.db.Query(
		"SELECT "+columnasNodo+" "+fromNodo+
			"WHERE padre_id = ? AND tipo = ?",
		PadreID, Tipo,
	)
	if err != nil {
		return nil, gko.ErrInesperado.Err(err).Op(op)
	}
	return s.scanRowsNodo(rows, op)
}
