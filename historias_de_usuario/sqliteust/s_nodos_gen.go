package sqliteust

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"

	"monorepo/historias_de_usuario/ust"

	"github.com/pargomx/gecko"
)

//  ================================================================  //
//  ========== MYSQL/CONSTANTES ====================================  //

// Lista de columnas separadas por coma para usar en consulta SELECT
// en conjunto con scanRow o scanRows, ya que las columnas coinciden
// con los campos escaneados.
const columnasNodo string = "nodo_id, nodo_tbl, padre_id, padre_tbl, nivel, posicion"

// Origen de los datos de ust.Nodo
//
// FROM nodos
const fromNodo string = "FROM nodos "

//  ================================================================  //
//  ========== MYSQL/TBL-INSERT ====================================  //

//  ================================================================  //
//  ========== MYSQL/TBL-UPDATE ====================================  //

// UpdateNodo valida y sobreescribe el registro en la base de datos.
func (s *Repositorio) UpdateNodo(nod ust.Nodo) error {
	const op string = "mysqlust.UpdateNodo"
	if nod.NodoID == 0 {
		return gecko.NewErr(http.StatusBadRequest).Msg("NodoID sin especificar").Ctx(op, "pk_indefinida")
	}
	if nod.NodoTbl == "" {
		return gecko.NewErr(http.StatusBadRequest).Msg("NodoTbl sin especificar").Ctx(op, "required_sin_valor")
	}
	if nod.PadreTbl == "" {
		return gecko.NewErr(http.StatusBadRequest).Msg("PadreTbl sin especificar").Ctx(op, "required_sin_valor")
	}
	err := nod.Validar()
	if err != nil {
		return gecko.NewErr(http.StatusBadRequest).Err(err).Op(op).Msg(err.Error())
	}
	_, err = s.db.Exec(
		"UPDATE nodos SET "+
			"nodo_id=?, nodo_tbl=?, padre_id=?, padre_tbl=?, nivel=?, posicion=? "+
			"WHERE nodo_id = ?",
		nod.NodoID, nod.NodoTbl, nod.PadreID, nod.PadreTbl, nod.Nivel, nod.Posicion,
		nod.NodoID,
	)
	if err != nil {
		return gecko.NewErr(http.StatusInternalServerError).Err(err).Op(op)
	}
	return nil
}

//  ================================================================  //
//  ========== MYSQL/TBL-DELETE ====================================  //

// DeleteNodo elimina permanentemente un registro del nodo de la base de datos.
// Error si el registro no existe o si no se da la clave primaria.
func (s *Repositorio) DeleteNodo(NodoID int) error {
	const op string = "mysqlust.DeleteNodo"
	if NodoID == 0 {
		return gecko.NewErr(http.StatusBadRequest).Msg("NodoID sin especificar").Ctx(op, "pk_indefinida")
	}
	// Verificar que solo se borre un registro.
	var num int
	err := s.db.QueryRow("SELECT COUNT(nodo_id) FROM nodos WHERE nodo_id = ?",
		NodoID,
	).Scan(&num)
	if err != nil {
		if err == sql.ErrNoRows {
			return gecko.NewErr(http.StatusNotFound).Err(ust.ErrNodoNotFound).Op(op)
		}
		return gecko.NewErr(http.StatusInternalServerError).Err(err).Op(op)
	}
	if num > 1 {
		return gecko.NewErr(http.StatusInternalServerError).Err(nil).Op(op).Msgf("abortado porque serían borrados %v registros", num)
	} else if num == 0 {
		return gecko.NewErr(http.StatusNotFound).Err(ust.ErrNodoNotFound).Op(op).Msg("cero resultados")
	}
	// Eliminar registro
	_, err = s.db.Exec(
		"DELETE FROM nodos WHERE nodo_id = ?",
		NodoID,
	)
	if err != nil {
		if strings.HasPrefix(err.Error(), "Error 1451 (23000)") {
			return gecko.NewErr(http.StatusConflict).Err(err).Op(op).Msg("Este registro es referenciado por otros y no se puede eliminar")
		} else {
			return gecko.NewErr(http.StatusInternalServerError).Err(err).Op(op)
		}
	}
	return nil
}

//  ================================================================  //
//  ========== MYSQL/SCAN-ROW ======================================  //

// Utilizar luego de un sql.QueryRow(). No es necesario hacer row.Close()
func (s *Repositorio) scanRowNodo(row *sql.Row, nod *ust.Nodo, op string) error {

	err := row.Scan(
		&nod.NodoID, &nod.NodoTbl, &nod.PadreID, &nod.PadreTbl, &nod.Nivel, &nod.Posicion,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return gecko.NewErr(http.StatusNotFound).Msg("el nodo no se encuentra").Op(op)
		}
		return gecko.NewErr(http.StatusInternalServerError).Err(err).Op(op)
	}

	return nil
}

//  ================================================================  //
//  ========== MYSQL/GET ===========================================  //

// GetNodo devuelve un Nodo de la DB por su clave primaria.
// Error si no encuentra ninguno, o si encuentra más de uno.
func (s *Repositorio) GetNodo(NodoID int) (*ust.Nodo, error) {
	const op string = "mysqlust.GetNodo"
	// if NodoID == 0 {
	// 	return nil, gecko.NewErr(http.StatusBadRequest).Msg("NodoID sin especificar").Ctx(op, "pk_indefinida")
	// }
	row := s.db.QueryRow(
		"SELECT "+columnasNodo+" "+fromNodo+
			"WHERE nodo_id = ?",
		NodoID,
	)
	nod := &ust.Nodo{}
	return nod, s.scanRowNodo(row, nod, op)
}

//  ================================================================  //
//  ========== MYSQL/SCAN-ROWS =====================================  //

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
			return nil, gecko.NewErr(http.StatusInternalServerError).Err(err).Op(op)
		}

		items = append(items, nod)
	}
	return items, nil
}

//  ================================================================  //
//  ========== MYSQL/LIST_BY =======================================  //

// ListNodosByPadreID retorna los registros a partir de PadreID.
func (s *Repositorio) ListNodosByPadreID(PadreID int) ([]ust.Nodo, error) {
	const op string = "mysqlust.ListNodosByPadreID"
	if PadreID == 0 {
		return nil, gecko.NewErr(http.StatusBadRequest).Msg("PadreID sin especificar").Ctx(op, "param_indefinido")
	}
	where := "WHERE padre_id = ?"
	argumentos := []any{}
	argumentos = append(argumentos, PadreID)

	rows, err := s.db.Query(
		"SELECT "+columnasNodo+" "+fromNodo+
			where,
		argumentos...,
	)
	if err != nil {
		return nil, gecko.NewErr(http.StatusInternalServerError).Err(err).Op(op)
	}
	return s.scanRowsNodo(rows, op)
}
