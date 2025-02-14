package sqliteust

import (
	"database/sql"
	"errors"

	"github.com/pargomx/gecko/gko"

	"monorepo/ust"
)

//  ================================================================  //
//  ========== INSERT ==============================================  //

func (s *Repositorio) InsertHistoria(his ust.Historia) error {
	const op string = "InsertHistoria"
	if his.HistoriaID == 0 {
		return gko.ErrDatoIndef().Op(op).Msg("HistoriaID sin especificar").Str("pk_indefinida")
	}
	if his.Titulo == "" {
		return gko.ErrDatoIndef().Op(op).Msg("Titulo sin especificar").Str("required_sin_valor")
	}
	_, err := s.db.Exec("INSERT INTO historias "+
		"(historia_id, titulo, objetivo, prioridad, completada, persona_id, proyecto_id, segundos_presupuesto, descripcion) "+
		"VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?) ",
		his.HistoriaID, his.Titulo, his.Objetivo, his.Prioridad, his.Completada, his.PersonaID, his.ProyectoID, his.SegundosPresupuesto, his.Descripcion,
	)
	if err != nil {
		return gko.ErrAlEscribir().Err(err).Op(op)
	}
	return nil
}

//  ================================================================  //
//  ========== UPDATE ==============================================  //

// UpdateHistoria valida y sobreescribe el registro en la base de datos.
func (s *Repositorio) UpdateHistoria(his ust.Historia) error {
	const op string = "UpdateHistoria"
	if his.HistoriaID == 0 {
		return gko.ErrDatoIndef().Op(op).Msg("HistoriaID sin especificar").Str("pk_indefinida")
	}
	if his.Titulo == "" {
		return gko.ErrDatoIndef().Op(op).Msg("Titulo sin especificar").Str("required_sin_valor")
	}
	_, err := s.db.Exec(
		"UPDATE historias SET "+
			"historia_id=?, titulo=?, objetivo=?, prioridad=?, completada=?, persona_id=?, proyecto_id=?, segundos_presupuesto=?, descripcion=? "+
			"WHERE historia_id = ?",
		his.HistoriaID, his.Titulo, his.Objetivo, his.Prioridad, his.Completada, his.PersonaID, his.ProyectoID, his.SegundosPresupuesto, his.Descripcion,
		his.HistoriaID,
	)
	if err != nil {
		return gko.ErrInesperado().Err(err).Op(op)
	}
	return nil
}

//  ================================================================  //
//  ========== EXISTE ==============================================  //

// Retorna error nil si existe solo un registro con esta clave primaria.
func (s *Repositorio) ExisteHistoria(HistoriaID int) error {
	const op string = "ExisteHistoria"
	var num int
	err := s.db.QueryRow("SELECT COUNT(historia_id) FROM historias WHERE historia_id = ?",
		HistoriaID,
	).Scan(&num)
	if err != nil {
		if err == sql.ErrNoRows {
			return gko.ErrNoEncontrado().Msg("Historia de usuario no encontrado").Op(op)
		}
		return gko.ErrInesperado().Err(err).Op(op)
	}
	if num > 1 {
		return gko.ErrInesperado().Err(nil).Op(op).Str("existen más de un registro para la pk").Ctx("registros", num)
	} else if num == 0 {
		return gko.ErrNoEncontrado().Msg("Historia de usuario no encontrado").Op(op)
	}
	return nil
}

//  ================================================================  //
//  ========== DELETE ==============================================  //

func (s *Repositorio) DeleteHistoria(HistoriaID int) error {
	const op string = "DeleteHistoria"
	if HistoriaID == 0 {
		return gko.ErrDatoIndef().Op(op).Msg("HistoriaID sin especificar").Str("pk_indefinida")
	}
	err := s.ExisteHistoria(HistoriaID)
	if err != nil {
		return gko.Err(err).Op(op)
	}
	_, err = s.db.Exec(
		"DELETE FROM historias WHERE historia_id = ?",
		HistoriaID,
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
//	historia_id,
//	titulo,
//	objetivo,
//	prioridad,
//	completada,
//	persona_id,
//	proyecto_id,
//	segundos_presupuesto,
//	descripcion
const columnasHistoria string = "historia_id, titulo, objetivo, prioridad, completada, persona_id, proyecto_id, segundos_presupuesto, descripcion"

// Origen de los datos de ust.Historia
//
//	FROM historias
const fromHistoria string = "FROM historias "

//  ================================================================  //
//  ========== SCAN ================================================  //

// Utilizar luego de un sql.QueryRow(). No es necesario hacer row.Close()
func (s *Repositorio) scanRowHistoria(row *sql.Row, his *ust.Historia) error {
	err := row.Scan(
		&his.HistoriaID, &his.Titulo, &his.Objetivo, &his.Prioridad, &his.Completada, &his.PersonaID, &his.ProyectoID, &his.SegundosPresupuesto, &his.Descripcion,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return gko.ErrNoEncontrado().Msg("Historia de usuario no encontrado")
		}
		return gko.ErrInesperado().Err(err)
	}
	return nil
}

//  ================================================================  //
//  ========== GET =================================================  //

// GetHistoria devuelve un Historia de la DB por su clave primaria.
// Error si no encuentra ninguno, o si encuentra más de uno.
func (s *Repositorio) GetHistoria(HistoriaID int) (*ust.Historia, error) {
	const op string = "GetHistoria"
	if HistoriaID == 0 {
		return nil, gko.ErrDatoIndef().Op(op).Msg("HistoriaID sin especificar").Str("pk_indefinida")
	}
	row := s.db.QueryRow(
		"SELECT "+columnasHistoria+" "+fromHistoria+
			"WHERE historia_id = ?",
		HistoriaID,
	)
	his := &ust.Historia{}
	err := s.scanRowHistoria(row, his)
	if err != nil {
		return nil, err
	}
	return his, nil
}

//  ================================================================  //
//  ========== SCAN ================================================  //

// scanRowsHistoria escanea cada row en la struct Historia
// y devuelve un slice con todos los items.
// Siempre se encarga de llamar rows.Close()
func (s *Repositorio) scanRowsHistoria(rows *sql.Rows, op string) ([]ust.Historia, error) {
	defer rows.Close()
	items := []ust.Historia{}
	for rows.Next() {
		his := ust.Historia{}
		err := rows.Scan(
			&his.HistoriaID, &his.Titulo, &his.Objetivo, &his.Prioridad, &his.Completada, &his.PersonaID, &his.ProyectoID, &his.SegundosPresupuesto, &his.Descripcion,
		)
		if err != nil {
			return nil, gko.ErrInesperado().Err(err).Op(op)
		}
		items = append(items, his)
	}
	return items, nil
}

//  ================================================================  //
//  ========== LIST ================================================  //

func (s *Repositorio) ListHistorias() ([]ust.Historia, error) {
	const op string = "ListHistorias"
	rows, err := s.db.Query(
		"SELECT " + columnasHistoria + " " + fromHistoria,
	)
	if err != nil {
		return nil, gko.ErrInesperado().Err(err).Op(op)
	}
	return s.scanRowsHistoria(rows, op)
}

//  ================================================================  //
//  ========== LIST_BY PROYECTO_ID =================================  //

func (s *Repositorio) ListHistoriasByProyectoID(ProyectoID string) ([]ust.Historia, error) {
	const op string = "ListHistoriasByProyectoID"
	if ProyectoID == "" {
		return nil, gko.ErrDatoIndef().Op(op).Msg("ProyectoID sin especificar").Str("param_indefinido")
	}
	rows, err := s.db.Query(
		"SELECT "+columnasHistoria+" "+fromHistoria+
			"WHERE proyecto_id = ?",
		ProyectoID,
	)
	if err != nil {
		return nil, gko.ErrInesperado().Err(err).Op(op)
	}
	return s.scanRowsHistoria(rows, op)
}

//  ================================================================  //
//  ========== LIST BYPADREID ======================================  //

func (s *Repositorio) ListHistoriasByPadreID(nodoID int) ([]ust.Historia, error) {
	const op string = "ListHistoriasByPadreID"
	rows, err := s.db.Query(
		"SELECT "+columnasHistoria+" "+fromHistoria+
			"JOIN nodos ON nodo_id = historia_id WHERE padre_id = ? ORDER BY posicion",
		nodoID,
	)
	if err != nil {
		return nil, gko.ErrInesperado().Err(err).Op(op)
	}
	return s.scanRowsHistoria(rows, op)
}
