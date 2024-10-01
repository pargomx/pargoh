package sqliteust

import (
	"database/sql"
	"errors"

	"github.com/pargomx/gecko/gko"

	"monorepo/ust"
)

//  ================================================================  //
//  ========== INSERT ==============================================  //

func (s *Repositorio) InsertIntervalo(interv ust.Intervalo) error {
	const op string = "InsertIntervalo"
	if interv.TareaID == 0 {
		return gko.ErrDatoIndef().Op(op).Msg("TareaID sin especificar").Str("pk_indefinida")
	}
	if interv.Inicio == "" {
		return gko.ErrDatoIndef().Op(op).Msg("Inicio sin especificar").Str("pk_indefinida")
	}
	_, err := s.db.Exec("INSERT INTO intervalos "+
		"(tarea_id, inicio, fin) "+
		"VALUES (?, ?, ?) ",
		interv.TareaID, interv.Inicio, interv.Fin,
	)
	if err != nil {
		return gko.ErrAlEscribir().Err(err).Op(op)
	}
	return nil
}

//  ================================================================  //
//  ========== UPDATE ==============================================  //

// UpdateIntervalo valida y sobreescribe el registro en la base de datos.
func (s *Repositorio) UpdateIntervalo(TareaID int, Inicio string, interv ust.Intervalo) error {
	const op string = "UpdateIntervalo"
	if interv.TareaID == 0 {
		return gko.ErrDatoIndef().Op(op).Msg("TareaID sin especificar").Str("pk_indefinida")
	}
	if interv.Inicio == "" {
		return gko.ErrDatoIndef().Op(op).Msg("Inicio sin especificar").Str("pk_indefinida")
	}
	_, err := s.db.Exec(
		"UPDATE intervalos SET "+
			"tarea_id=?, inicio=?, fin=? "+
			"WHERE tarea_id = ? AND inicio = ?",
		interv.TareaID, interv.Inicio, interv.Fin,
		TareaID, Inicio,
	)
	if err != nil {
		return gko.ErrInesperado().Err(err).Op(op)
	}
	return nil
}

//  ================================================================  //
//  ========== EXISTE ==============================================  //

// Retorna error nil si existe solo un registro con esta clave primaria.
func (s *Repositorio) ExisteIntervalo(TareaID int, Inicio string) error {
	const op string = "ExisteIntervalo"
	var num int
	err := s.db.QueryRow("SELECT COUNT(tarea_id) FROM intervalos WHERE tarea_id = ? AND inicio = ?",
		TareaID, Inicio,
	).Scan(&num)
	if err != nil {
		if err == sql.ErrNoRows {
			return gko.ErrNoEncontrado().Err(ust.ErrIntervaloNotFound).Op(op)
		}
		return gko.ErrInesperado().Err(err).Op(op)
	}
	if num > 1 {
		return gko.ErrInesperado().Err(nil).Op(op).Str("existen más de un registro para la pk").Ctx("registros", num)
	} else if num == 0 {
		return gko.ErrNoEncontrado().Err(ust.ErrIntervaloNotFound).Op(op)
	}
	return nil
}

//  ================================================================  //
//  ========== DELETE ==============================================  //

func (s *Repositorio) DeleteIntervalo(TareaID int, Inicio string) error {
	const op string = "DeleteIntervalo"
	if TareaID == 0 {
		return gko.ErrDatoIndef().Op(op).Msg("TareaID sin especificar").Str("pk_indefinida")
	}
	if Inicio == "" {
		return gko.ErrDatoIndef().Op(op).Msg("Inicio sin especificar").Str("pk_indefinida")
	}
	err := s.ExisteIntervalo(TareaID, Inicio)
	if err != nil {
		return gko.Err(err).Op(op)
	}
	_, err = s.db.Exec(
		"DELETE FROM intervalos WHERE tarea_id = ? AND inicio = ?",
		TareaID, Inicio,
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
//	tarea_id,
//	inicio,
//	fin
const columnasIntervalo string = "tarea_id, inicio, fin"

// Origen de los datos de ust.Intervalo
//
//	FROM intervalos
const fromIntervalo string = "FROM intervalos "

//  ================================================================  //
//  ========== SCAN ================================================  //

// Utilizar luego de un sql.QueryRow(). No es necesario hacer row.Close()
func (s *Repositorio) scanRowIntervalo(row *sql.Row, interv *ust.Intervalo) error {
	err := row.Scan(
		&interv.TareaID, &interv.Inicio, &interv.Fin,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return gko.ErrNoEncontrado().Msg("Intervalo no se encuentra")
		}
		return gko.ErrInesperado().Err(err)
	}
	return nil
}

//  ================================================================  //
//  ========== GET =================================================  //

// GetIntervalo devuelve un Intervalo de la DB por su clave primaria.
// Error si no encuentra ninguno, o si encuentra más de uno.
func (s *Repositorio) GetIntervalo(TareaID int, Inicio string) (*ust.Intervalo, error) {
	const op string = "GetIntervalo"
	if TareaID == 0 {
		return nil, gko.ErrDatoIndef().Op(op).Msg("TareaID sin especificar").Str("pk_indefinida")
	}
	if Inicio == "" {
		return nil, gko.ErrDatoIndef().Op(op).Msg("Inicio sin especificar").Str("pk_indefinida")
	}
	row := s.db.QueryRow(
		"SELECT "+columnasIntervalo+" "+fromIntervalo+
			"WHERE tarea_id = ? AND inicio = ?",
		TareaID, Inicio,
	)
	interv := &ust.Intervalo{}
	err := s.scanRowIntervalo(row, interv)
	if err != nil {
		return nil, err
	}
	return interv, nil
}

//  ================================================================  //
//  ========== SCAN ================================================  //

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
			return nil, gko.ErrInesperado().Err(err).Op(op)
		}
		items = append(items, interv)
	}
	return items, nil
}

//  ================================================================  //
//  ========== LIST_BY TAREA_ID ====================================  //

func (s *Repositorio) ListIntervalosByTareaID(TareaID int) ([]ust.Intervalo, error) {
	const op string = "ListIntervalosByTareaID"
	if TareaID == 0 {
		return nil, gko.ErrDatoIndef().Op(op).Msg("TareaID sin especificar").Str("param_indefinido")
	}
	rows, err := s.db.Query(
		"SELECT "+columnasIntervalo+" "+fromIntervalo+
			"WHERE tarea_id = ?",
		TareaID,
	)
	if err != nil {
		return nil, gko.ErrInesperado().Err(err).Op(op)
	}
	return s.scanRowsIntervalo(rows, op)
}
