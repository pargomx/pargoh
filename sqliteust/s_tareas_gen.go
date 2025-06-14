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
		return gko.ErrDatoIndef.Str("pk_indefinida").Op(op).Msg("TareaID sin especificar")
	}
	_, err := s.db.Exec("INSERT INTO tareas "+
		"(tarea_id, historia_id, descripcion, importancia, tipo, estatus, impedimentos, segundos_estimado, segundos_real) "+
		"VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?) ",
		tar.TareaID, tar.HistoriaID, tar.Descripcion, tar.Importancia.String, tar.Tipo.String, tar.Estatus, tar.Impedimentos, tar.SegundosEstimado, tar.SegundosUtilizado,
	)
	if err != nil {
		return gko.ErrAlEscribir.Err(err).Op(op)
	}
	return nil
}

//  ================================================================  //
//  ========== UPDATE ==============================================  //

// UpdateTarea valida y sobreescribe el registro en la base de datos.
func (s *Repositorio) UpdateTarea(tar ust.Tarea) error {
	const op string = "UpdateTarea"
	if tar.TareaID == 0 {
		return gko.ErrDatoIndef.Str("pk_indefinida").Op(op).Msg("TareaID sin especificar")
	}
	_, err := s.db.Exec(
		"UPDATE tareas SET "+
			"tarea_id=?, historia_id=?, descripcion=?, importancia=?, tipo=?, estatus=?, impedimentos=?, segundos_estimado=?, segundos_real=? "+
			"WHERE tarea_id = ?",
		tar.TareaID, tar.HistoriaID, tar.Descripcion, tar.Importancia.String, tar.Tipo.String, tar.Estatus, tar.Impedimentos, tar.SegundosEstimado, tar.SegundosUtilizado,
		tar.TareaID,
	)
	if err != nil {
		return gko.ErrInesperado.Err(err).Op(op)
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
			return gko.ErrNoEncontrado.Msg("Tarea no encontrado").Op(op)
		}
		return gko.ErrInesperado.Err(err).Op(op)
	}
	if num > 1 {
		return gko.ErrInesperado.Err(nil).Op(op).Str("existen más de un registro para la pk").Ctx("registros", num)
	} else if num == 0 {
		return gko.ErrNoEncontrado.Msg("Tarea no encontrado").Op(op)
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
//	descripcion,
//	importancia,
//	tipo,
//	estatus,
//	impedimentos,
//	segundos_estimado,
//	segundos_real
const columnasTarea string = "tarea_id, historia_id, descripcion, importancia, tipo, estatus, impedimentos, segundos_estimado, segundos_real"

// Origen de los datos de ust.Tarea
//
//	FROM tareas
const fromTarea string = "FROM tareas "

//  ================================================================  //
//  ========== SCAN ================================================  //

// Utilizar luego de un sql.QueryRow(). No es necesario hacer row.Close()
func (s *Repositorio) scanRowTarea(row *sql.Row, tar *ust.Tarea) error {
	var importancia string
	var tipo string
	err := row.Scan(
		&tar.TareaID, &tar.HistoriaID, &tar.Descripcion, &importancia, &tipo, &tar.Estatus, &tar.Impedimentos, &tar.SegundosEstimado, &tar.SegundosUtilizado,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return gko.ErrNoEncontrado.Msg("Tarea no encontrado")
		}
		return gko.ErrInesperado.Err(err)
	}
	tar.Importancia = ust.SetImportanciaTareaDB(importancia)
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
		return nil, gko.ErrDatoIndef.Str("pk_indefinida").Op(op).Msg("TareaID sin especificar")
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
		var importancia string
		var tipo string
		err := rows.Scan(
			&tar.TareaID, &tar.HistoriaID, &tar.Descripcion, &importancia, &tipo, &tar.Estatus, &tar.Impedimentos, &tar.SegundosEstimado, &tar.SegundosUtilizado,
		)
		if err != nil {
			return nil, gko.ErrInesperado.Err(err).Op(op)
		}
		tar.Importancia = ust.SetImportanciaTareaDB(importancia)
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
		"SELECT " + columnasTarea + " " + fromTarea,
	)
	if err != nil {
		return nil, gko.ErrInesperado.Err(err).Op(op)
	}
	return s.scanRowsTarea(rows, op)
}

//  ================================================================  //
//  ========== LIST_BY HISTORIA_ID =================================  //

func (s *Repositorio) ListTareasByHistoriaID(HistoriaID int) ([]ust.Tarea, error) {
	const op string = "ListTareasByHistoriaID"
	if HistoriaID == 0 {
		return nil, gko.ErrDatoIndef.Str("param_indefinido").Op(op).Msg("HistoriaID sin especificar")
	}
	rows, err := s.db.Query(
		"SELECT "+columnasTarea+" "+fromTarea+
			"WHERE historia_id = ?",
		HistoriaID,
	)
	if err != nil {
		return nil, gko.ErrInesperado.Err(err).Op(op)
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
		return nil, gko.ErrInesperado.Err(err).Op(op)
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
		return nil, gko.ErrInesperado.Err(err).Op(op)
	}
	return s.scanRowsTarea(rows, op)
}
