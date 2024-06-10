package sqliteust

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"

	"monorepo/gecko"
	"monorepo/historias_de_usuario/ust"
)

//  ================================================================  //
//  ========== MYSQL/CONSTANTES ====================================  //

// Lista de columnas separadas por coma para usar en consulta SELECT
// en conjunto con scanRow o scanRows, ya que las columnas coinciden
// con los campos escaneados.
const columnasTarea string = "tarea_id, historia_id, tipo, descripcion, impedimentos, tiempo_estimado, tiempo_real, estatus"

// Origen de los datos de ust.Tarea
//
// FROM tareas
const fromTarea string = "FROM tareas "

//  ================================================================  //
//  ========== MYSQL/TBL-INSERT ====================================  //

// InsertTarea valida el registro y lo inserta en la base de datos.
func (s *Repositorio) InsertTarea(tar ust.Tarea) error {
	const op string = "mysqlust.InsertTarea"
	if tar.TareaID == 0 {
		return gecko.NewErr(http.StatusBadRequest).Msg("TareaID sin especificar").Ctx(op, "pk_indefinida")
	}
	err := tar.Validar()
	if err != nil {
		return gecko.NewErr(http.StatusBadRequest).Err(err).Op(op).Msg(err.Error())
	}
	_, err = s.db.Exec("INSERT INTO tareas "+
		"(tarea_id, historia_id, tipo, descripcion, impedimentos, tiempo_estimado, tiempo_real, estatus) "+
		"VALUES (?, ?, ?, ?, ?, ?, ?, ?) ",
		tar.TareaID, tar.HistoriaID, tar.Tipo.String, tar.Descripcion, tar.Impedimentos, tar.TiempoEstimado, tar.TiempoReal, tar.Estatus,
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

// UpdateTarea valida y sobreescribe el registro en la base de datos.
func (s *Repositorio) UpdateTarea(tar ust.Tarea) error {
	const op string = "mysqlust.UpdateTarea"
	if tar.TareaID == 0 {
		return gecko.NewErr(http.StatusBadRequest).Msg("TareaID sin especificar").Ctx(op, "pk_indefinida")
	}
	err := tar.Validar()
	if err != nil {
		return gecko.NewErr(http.StatusBadRequest).Err(err).Op(op).Msg(err.Error())
	}
	_, err = s.db.Exec(
		"UPDATE tareas SET "+
			"tarea_id=?, historia_id=?, tipo=?, descripcion=?, impedimentos=?, tiempo_estimado=?, tiempo_real=?, estatus=? "+
			"WHERE tarea_id = ?",
		tar.TareaID, tar.HistoriaID, tar.Tipo.String, tar.Descripcion, tar.Impedimentos, tar.TiempoEstimado, tar.TiempoReal, tar.Estatus,
		tar.TareaID,
	)
	if err != nil {
		return gecko.NewErr(http.StatusInternalServerError).Err(err).Op(op)
	}
	return nil
}

//  ================================================================  //
//  ========== MYSQL/TBL-DELETE ====================================  //

func (s *Repositorio) DeleteTarea(TareaID int) error {
	const op string = "mysqlust.DeleteTarea"
	if TareaID == 0 {
		return gecko.NewErr(http.StatusBadRequest).Msg("TareaID sin especificar").Ctx(op, "pk_indefinida")
	}
	// Verificar que solo se borre un registro.
	var num int
	err := s.db.QueryRow("SELECT COUNT(tarea_id) FROM tareas WHERE tarea_id = ?",
		TareaID,
	).Scan(&num)
	if err != nil {
		if err == sql.ErrNoRows {
			return gecko.NewErr(http.StatusNotFound).Err(ust.ErrTareaNotFound).Op(op)
		}
		return gecko.NewErr(http.StatusInternalServerError).Err(err).Op(op)
	}
	if num > 1 {
		return gecko.NewErr(http.StatusInternalServerError).Err(nil).Op(op).Msgf("abortado porque serían borrados %v registros", num)
	} else if num == 0 {
		return gecko.NewErr(http.StatusNotFound).Err(ust.ErrTareaNotFound).Op(op).Msg("cero resultados")
	}
	// Eliminar registro
	_, err = s.db.Exec(
		"DELETE FROM tareas WHERE tarea_id = ?",
		TareaID,
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
func (s *Repositorio) scanRowTarea(row *sql.Row, tar *ust.Tarea, op string) error {
	var tipo string
	err := row.Scan(
		&tar.TareaID, &tar.HistoriaID, &tipo, &tar.Descripcion, &tar.Impedimentos, &tar.TiempoEstimado, &tar.TiempoReal, &tar.Estatus,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return gecko.NewErr(http.StatusNotFound).Msg("Tarea no se encuentra").Op(op)
		}
		return gecko.NewErr(http.StatusInternalServerError).Err(err).Op(op)
	}
	tar.Tipo = ust.SetTipoTareaDB(tipo)
	return nil
}

//  ================================================================  //
//  ========== MYSQL/GET ===========================================  //

// GetTarea devuelve un Tarea de la DB por su clave primaria.
// Error si no encuentra ninguno, o si encuentra más de uno.
func (s *Repositorio) GetTarea(TareaID int) (*ust.Tarea, error) {
	const op string = "mysqlust.GetTarea"
	if TareaID == 0 {
		return nil, gecko.NewErr(http.StatusBadRequest).Msg("TareaID sin especificar").Ctx(op, "pk_indefinida")
	}
	row := s.db.QueryRow(
		"SELECT "+columnasTarea+" "+fromTarea+
			"WHERE tarea_id = ?",
		TareaID,
	)
	tar := &ust.Tarea{}
	return tar, s.scanRowTarea(row, tar, op)
}

//  ================================================================  //
//  ========== MYSQL/SCAN-ROWS =====================================  //

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
			&tar.TareaID, &tar.HistoriaID, &tipo, &tar.Descripcion, &tar.Impedimentos, &tar.TiempoEstimado, &tar.TiempoReal, &tar.Estatus,
		)
		if err != nil {
			return nil, gecko.NewErr(http.StatusInternalServerError).Err(err).Op(op)
		}
		tar.Tipo = ust.SetTipoTareaDB(tipo)
		items = append(items, tar)
	}
	return items, nil
}

//  ================================================================  //
//  ========== MYSQL/LIST_BY =======================================  //

func (s *Repositorio) ListTareasByHistoriaID(HistoriaID int) ([]ust.Tarea, error) {
	const op string = "mysqlust.ListTareasByHistoriaID"
	if HistoriaID == 0 {
		return nil, gecko.NewErr(http.StatusBadRequest).Msg("HistoriaID sin especificar").Ctx(op, "param_indefinido")
	}
	rows, err := s.db.Query(
		"SELECT "+columnasTarea+" "+fromTarea+
			"WHERE historia_id = ?",
		HistoriaID,
	)
	if err != nil {
		return nil, gecko.NewErr(http.StatusInternalServerError).Err(err).Op(op)
	}
	return s.scanRowsTarea(rows, op)
}
