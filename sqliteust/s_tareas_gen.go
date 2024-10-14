package sqliteust

import (
	"database/sql"
	"errors"

	"github.com/pargomx/gecko/gko"

	"monorepo/ust"
)

//  ================================================================  //
//  ========== INSERT ==============================================  //

func (s *Repositorio) InsertTarea(tar ust.Tarea) error {
	const op string = "InsertTarea"
	if tar.TareaID == 0 {
		return gko.ErrDatoIndef().Op(op).Msg("TareaID sin especificar").Str("pk_indefinida")
	}
	_, err := s.db.Exec("INSERT INTO tareas "+
		"(tarea_id, historia_id, tipo, descripcion, impedimentos, segundos_estimado, segundos_real, estatus, importancia) "+
		"VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?) ",
		tar.TareaID, tar.HistoriaID, tar.Tipo.String, tar.Descripcion, tar.Impedimentos, tar.SegundosEstimado, tar.SegundosReal, tar.Estatus, tar.Importancia,
	)
	if err != nil {
		return gko.ErrAlEscribir().Err(err).Op(op)
	}
	return nil
}

//  ================================================================  //
//  ========== UPDATE ==============================================  //

// UpdateTarea valida y sobreescribe el registro en la base de datos.
func (s *Repositorio) UpdateTarea(tar ust.Tarea) error {
	const op string = "UpdateTarea"
	if tar.TareaID == 0 {
		return gko.ErrDatoIndef().Op(op).Msg("TareaID sin especificar").Str("pk_indefinida")
	}
	_, err := s.db.Exec(
		"UPDATE tareas SET "+
			"tarea_id=?, historia_id=?, tipo=?, descripcion=?, impedimentos=?, segundos_estimado=?, segundos_real=?, estatus=?, importancia=? "+
			"WHERE tarea_id = ?",
		tar.TareaID, tar.HistoriaID, tar.Tipo.String, tar.Descripcion, tar.Impedimentos, tar.SegundosEstimado, tar.SegundosReal, tar.Estatus, tar.Importancia,
		tar.TareaID,
	)
	if err != nil {
		return gko.ErrInesperado().Err(err).Op(op)
	}
	return nil
}

//  ================================================================  //
//  ========== EXISTE ==============================================  //

// Retorna error nil si existe solo un registro con esta clave primaria.
func (s *Repositorio) ExisteTarea(TareaID int) error {
	const op string = "ExisteTarea"
	var num int
	err := s.db.QueryRow("SELECT COUNT(tarea_id) FROM tareas WHERE tarea_id = ?",
		TareaID,
	).Scan(&num)
	if err != nil {
		if err == sql.ErrNoRows {
			return gko.ErrNoEncontrado().Err(ust.ErrTareaNotFound).Op(op)
		}
		return gko.ErrInesperado().Err(err).Op(op)
	}
	if num > 1 {
		return gko.ErrInesperado().Err(nil).Op(op).Str("existen más de un registro para la pk").Ctx("registros", num)
	} else if num == 0 {
		return gko.ErrNoEncontrado().Err(ust.ErrTareaNotFound).Op(op)
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
//	historia_id,
//	tipo,
//	descripcion,
//	impedimentos,
//	segundos_estimado,
//	segundos_real,
//	estatus,
//	importancia
const columnasTarea string = "tarea_id, historia_id, tipo, descripcion, impedimentos, segundos_estimado, segundos_real, estatus, importancia"

// Origen de los datos de ust.Tarea
//
//	FROM tareas
const fromTarea string = "FROM tareas "

//  ================================================================  //
//  ========== SCAN ================================================  //

// Utilizar luego de un sql.QueryRow(). No es necesario hacer row.Close()
func (s *Repositorio) scanRowTarea(row *sql.Row, tar *ust.Tarea) error {
	var tipo string
	err := row.Scan(
		&tar.TareaID, &tar.HistoriaID, &tipo, &tar.Descripcion, &tar.Impedimentos, &tar.SegundosEstimado, &tar.SegundosReal, &tar.Estatus, &tar.Importancia,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return gko.ErrNoEncontrado().Msg("Tarea no se encuentra")
		}
		return gko.ErrInesperado().Err(err)
	}
	tar.Tipo = ust.SetTipoTareaDB(tipo)
	return nil
}

//  ================================================================  //
//  ========== GET =================================================  //

// GetTarea devuelve un Tarea de la DB por su clave primaria.
// Error si no encuentra ninguno, o si encuentra más de uno.
func (s *Repositorio) GetTarea(TareaID int) (*ust.Tarea, error) {
	const op string = "GetTarea"
	if TareaID == 0 {
		return nil, gko.ErrDatoIndef().Op(op).Msg("TareaID sin especificar").Str("pk_indefinida")
	}
	row := s.db.QueryRow(
		"SELECT "+columnasTarea+" "+fromTarea+
			"WHERE tarea_id = ?",
		TareaID,
	)
	tar := &ust.Tarea{}
	err := s.scanRowTarea(row, tar)
	if err != nil {
		return nil, err
	}
	return tar, nil
}

//  ================================================================  //
//  ========== SCAN ================================================  //

// scanRowsTarea escanea cada row en la struct Tarea
// y devuelve un slice con todos los items.
// Siempre se encarga de llamar rows.Close()
func (s *Repositorio) scanRowsTarea(rows *sql.Rows, op string) ([]ust.Tarea, error) {
	defer rows.Close()
	items := []ust.Tarea{}
	for rows.Next() {
		tar := ust.Tarea{}
		var tipo string
		err := rows.Scan(
			&tar.TareaID, &tar.HistoriaID, &tipo, &tar.Descripcion, &tar.Impedimentos, &tar.SegundosEstimado, &tar.SegundosReal, &tar.Estatus, &tar.Importancia,
		)
		if err != nil {
			return nil, gko.ErrInesperado().Err(err).Op(op)
		}
		tar.Tipo = ust.SetTipoTareaDB(tipo)
		items = append(items, tar)
	}
	return items, nil
}

//  ================================================================  //
//  ========== LIST ================================================  //

func (s *Repositorio) ListTareas() ([]ust.Tarea, error) {
	const op string = "ListTareas"
	rows, err := s.db.Query(
		"SELECT " + columnasTarea + " " + fromTarea + "ORDER BY importancia DESC",
	)
	if err != nil {
		return nil, gko.ErrInesperado().Err(err).Op(op)
	}
	return s.scanRowsTarea(rows, op)
}

//  ================================================================  //
//  ========== LIST_BY HISTORIA_ID =================================  //

func (s *Repositorio) ListTareasByHistoriaID(HistoriaID int) ([]ust.Tarea, error) {
	const op string = "ListTareasByHistoriaID"
	if HistoriaID == 0 {
		return nil, gko.ErrDatoIndef().Op(op).Msg("HistoriaID sin especificar").Str("param_indefinido")
	}
	rows, err := s.db.Query(
		"SELECT "+columnasTarea+" "+fromTarea+
			"WHERE historia_id = ? ORDER BY importancia DESC",
		HistoriaID,
	)
	if err != nil {
		return nil, gko.ErrInesperado().Err(err).Op(op)
	}
	return s.scanRowsTarea(rows, op)
}

//  ================================================================  //
//  ========== LIST BUGS ===========================================  //

func (s *Repositorio) ListTareasBugs() ([]ust.Tarea, error) {
	const op string = "ListTareasBugs"
	rows, err := s.db.Query(
		"SELECT " + columnasTarea + " " + fromTarea +
			"WHERE tipo = 'BUG' AND estatus < 3",
	)
	if err != nil {
		return nil, gko.ErrInesperado().Err(err).Op(op)
	}
	return s.scanRowsTarea(rows, op)
}

//  ================================================================  //
//  ========== LIST ENCURSO ========================================  //

func (s *Repositorio) ListTareasEnCurso() ([]ust.Tarea, error) {
	const op string = "ListTareasEnCurso"
	rows, err := s.db.Query(
		"SELECT " + columnasTarea + " " + fromTarea +
			"WHERE estatus = 1",
	)
	if err != nil {
		return nil, gko.ErrInesperado().Err(err).Op(op)
	}
	return s.scanRowsTarea(rows, op)
}
