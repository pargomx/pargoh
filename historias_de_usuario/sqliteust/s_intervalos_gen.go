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
const columnasIntervalo string = "tarea_id, inicio, fin"

// Origen de los datos de ust.Intervalo
//
// FROM intervalos
const fromIntervalo string = "FROM intervalos "

//  ================================================================  //
//  ========== MYSQL/TBL-INSERT ====================================  //

// InsertIntervalo valida el registro y lo inserta en la base de datos.
func (s *Repositorio) InsertIntervalo(interv ust.Intervalo) error {
	const op string = "mysqlust.InsertIntervalo"
	if interv.TareaID == 0 {
		return gecko.NewErr(http.StatusBadRequest).Msg("TareaID sin especificar").Ctx(op, "pk_indefinida")
	}
	if interv.Inicio == "" {
		return gecko.NewErr(http.StatusBadRequest).Msg("Inicio sin especificar").Ctx(op, "pk_indefinida")
	}
	err := interv.Validar()
	if err != nil {
		return gecko.NewErr(http.StatusBadRequest).Err(err).Op(op).Msg(err.Error())
	}
	_, err = s.db.Exec("INSERT INTO intervalos "+
		"(tarea_id, inicio, fin) "+
		"VALUES (?, ?, ?) ",
		interv.TareaID, interv.Inicio, interv.Fin,
	)
	if err != nil {
		if strings.HasPrefix(err.Error(), "Error 1062 (23000)") {
			return gecko.NewErr(http.StatusConflict).Err(err).Op(op)
		} else if strings.HasPrefix(err.Error(), "Error 1452 (23000)") {
			return gecko.NewErr(http.StatusBadRequest).Err(err).Op(op).Msg("No se puede insertar la información porque el registro asociado no existe")
		} else {
			return gecko.NewErr(http.StatusInternalServerError).Err(err).Op(op)
		}
	}
	return nil
}

//  ================================================================  //
//  ========== MYSQL/TBL-UPDATE ====================================  //

// UpdateIntervalo valida y sobreescribe el registro en la base de datos.
func (s *Repositorio) UpdateIntervalo(interv ust.Intervalo) error {
	const op string = "mysqlust.UpdateIntervalo"
	if interv.TareaID == 0 {
		return gecko.NewErr(http.StatusBadRequest).Msg("TareaID sin especificar").Ctx(op, "pk_indefinida")
	}
	if interv.Inicio == "" {
		return gecko.NewErr(http.StatusBadRequest).Msg("Inicio sin especificar").Ctx(op, "pk_indefinida")
	}
	err := interv.Validar()
	if err != nil {
		return gecko.NewErr(http.StatusBadRequest).Err(err).Op(op).Msg(err.Error())
	}
	_, err = s.db.Exec(
		"UPDATE intervalos SET "+
			"tarea_id=?, inicio=?, fin=? "+
			"WHERE tarea_id = ? AND inicio = ?",
		interv.TareaID, interv.Inicio, interv.Fin,
		interv.TareaID, interv.Inicio,
	)
	if err != nil {
		return gecko.NewErr(http.StatusInternalServerError).Err(err).Op(op)
	}
	return nil
}

//  ================================================================  //
//  ========== MYSQL/SCAN-ROW ======================================  //

// Utilizar luego de un sql.QueryRow(). No es necesario hacer row.Close()
func (s *Repositorio) scanRowIntervalo(row *sql.Row, interv *ust.Intervalo, op string) error {

	err := row.Scan(
		&interv.TareaID, &interv.Inicio, &interv.Fin,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return gecko.NewErr(http.StatusNotFound).Msg("Intervalo no se encuentra").Op(op)
		}
		return gecko.NewErr(http.StatusInternalServerError).Err(err).Op(op)
	}

	return nil
}

//  ================================================================  //
//  ========== MYSQL/GET ===========================================  //

// GetIntervalo devuelve un Intervalo de la DB por su clave primaria.
// Error si no encuentra ninguno, o si encuentra más de uno.
func (s *Repositorio) GetIntervalo(TareaID int, Inicio string) (*ust.Intervalo, error) {
	const op string = "mysqlust.GetIntervalo"
	if TareaID == 0 {
		return nil, gecko.NewErr(http.StatusBadRequest).Msg("TareaID sin especificar").Ctx(op, "pk_indefinida")
	}
	if Inicio == "" {
		return nil, gecko.NewErr(http.StatusBadRequest).Msg("Inicio sin especificar").Ctx(op, "pk_indefinida")
	}
	row := s.db.QueryRow(
		"SELECT "+columnasIntervalo+" "+fromIntervalo+
			"WHERE tarea_id = ? AND inicio = ?",
		TareaID, Inicio,
	)
	interv := &ust.Intervalo{}
	return interv, s.scanRowIntervalo(row, interv, op)
}

//  ================================================================  //
//  ========== MYSQL/SCAN-ROWS =====================================  //

// scanRowsIntervalo escanea cada row en la struct Intervalo
// y devuelve un slice con todos los items.
// Siempre se encarga de llamar rows.Close()
func (s *Repositorio) scanRowsIntervalo(rows *sql.Rows, op string) ([]ust.Intervalo, error) {
	defer rows.Close()
	items := []ust.Intervalo{}
	for rows.Next() {
		interv := ust.Intervalo{}

		err := rows.Scan(
			&interv.TareaID, &interv.Inicio, &interv.Fin,
		)
		if err != nil {
			return nil, gecko.NewErr(http.StatusInternalServerError).Err(err).Op(op)
		}

		items = append(items, interv)
	}
	return items, nil
}

//  ================================================================  //
//  ========== MYSQL/LIST_BY =======================================  //

func (s *Repositorio) ListIntervalosByTareaID(TareaID int) ([]ust.Intervalo, error) {
	const op string = "mysqlust.ListIntervalosByTareaID"
	if TareaID == 0 {
		return nil, gecko.NewErr(http.StatusBadRequest).Msg("TareaID sin especificar").Ctx(op, "param_indefinido")
	}
	rows, err := s.db.Query(
		"SELECT "+columnasIntervalo+" "+fromIntervalo+
			"WHERE tarea_id = ?",
		TareaID,
	)
	if err != nil {
		return nil, gecko.NewErr(http.StatusInternalServerError).Err(err).Op(op)
	}
	return s.scanRowsIntervalo(rows, op)
}
