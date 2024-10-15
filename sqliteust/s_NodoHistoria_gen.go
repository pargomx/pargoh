package sqliteust

import (
	"database/sql"
	"errors"

	"github.com/pargomx/gecko/gko"

	"monorepo/ust"
)

//  ================================================================  //
//  ========== CONSTANTES ==========================================  //

// Lista de columnas separadas por coma para usar en consulta SELECT
// en conjunto con scanRow o scanRows, ya que las columnas coinciden
// con los campos escaneados.
//
//	his.historia_id,
//	his.proyecto_id AS proyecto_id,
//	his.persona_id AS persona_id,
//	his.titulo,
//	his.objetivo,
//	his.prioridad,
//	his.completada,
//	coalesce(nod.padre_id, 0),
//	coalesce(nod.padre_tbl, ''),
//	coalesce(nod.nivel, 0),
//	coalesce(nod.posicion, 0),
//	his.segundos_presupuesto,
//	coalesce((SELECT COUNT(nodo_id) FROM nodos WHERE padre_id = his.historia_id), 0) AS num_historias,
//	coalesce((SELECT COUNT(tarea_id) FROM tareas WHERE historia_id = his.historia_id), 0) AS num_tareas,
//	coalesce((SELECT SUM(segundos_estimado) FROM tareas WHERE historia_id = his.historia_id), 0) AS segundos_estimado,
//	coalesce((SELECT SUM(unixepoch(coalesce(nullif(interv.fin,''),datetime('now','-6 hours'))) - unixepoch(interv.inicio)) FROM intervalos interv JOIN tareas tar ON tar.tarea_id = interv.tarea_id WHERE tar.historia_id = his.historia_id GROUP BY tar.historia_id), 0) AS segundos_real
const columnasNodoHistoria string = "his.historia_id, his.proyecto_id AS proyecto_id, his.persona_id AS persona_id, his.titulo, his.objetivo, his.prioridad, his.completada, coalesce(nod.padre_id, 0), coalesce(nod.padre_tbl, ''), coalesce(nod.nivel, 0), coalesce(nod.posicion, 0), his.segundos_presupuesto, coalesce((SELECT COUNT(nodo_id) FROM nodos WHERE padre_id = his.historia_id), 0) AS num_historias, coalesce((SELECT COUNT(tarea_id) FROM tareas WHERE historia_id = his.historia_id), 0) AS num_tareas, coalesce((SELECT SUM(segundos_estimado) FROM tareas WHERE historia_id = his.historia_id), 0) AS segundos_estimado, coalesce((SELECT SUM(unixepoch(coalesce(nullif(interv.fin,''),datetime('now','-6 hours'))) - unixepoch(interv.inicio)) FROM intervalos interv JOIN tareas tar ON tar.tarea_id = interv.tarea_id WHERE tar.historia_id = his.historia_id GROUP BY tar.historia_id), 0) AS segundos_real"

// Origen de los datos de ust.NodoHistoria
//
//	FROM historias his
//	INNER JOIN nodos nod ON nodo_id = historia_id
const fromNodoHistoria string = "FROM historias his INNER JOIN nodos nod ON nodo_id = historia_id "

//  ================================================================  //
//  ========== SCAN ================================================  //

// Utilizar luego de un sql.QueryRow(). No es necesario hacer row.Close()
func (s *Repositorio) scanRowNodoHistoria(row *sql.Row, nhist *ust.NodoHistoria) error {
	err := row.Scan(
		&nhist.HistoriaID, &nhist.ProyectoID, &nhist.PersonaID, &nhist.Titulo, &nhist.Objetivo, &nhist.Prioridad, &nhist.Completada, &nhist.PadreID, &nhist.PadreTbl, &nhist.Nivel, &nhist.Posicion, &nhist.SegundosPresupuesto, &nhist.NumHistorias, &nhist.NumTareas, &nhist.SegundosEstimado, &nhist.SegundosReal,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return gko.ErrNoEncontrado().Msg("Historia de usuario no se encuentra")
		}
		return gko.ErrInesperado().Err(err)
	}
	return nil
}

//  ================================================================  //
//  ========== GET =================================================  //

// GetNodoHistoria devuelve un NodoHistoria de la DB por su clave primaria.
// Error si no encuentra ninguno, o si encuentra más de uno.
func (s *Repositorio) GetNodoHistoria(HistoriaID int) (*ust.NodoHistoria, error) {
	const op string = "GetNodoHistoria"
	if HistoriaID == 0 {
		return nil, gko.ErrDatoIndef().Op(op).Msg("HistoriaID sin especificar").Str("pk_indefinida")
	}
	row := s.db.QueryRow(
		"SELECT "+columnasNodoHistoria+" "+fromNodoHistoria+
			"WHERE his.historia_id = ?",
		HistoriaID,
	)
	nhist := &ust.NodoHistoria{}
	err := s.scanRowNodoHistoria(row, nhist)
	if err != nil {
		return nil, err
	}
	return nhist, nil
}

//  ================================================================  //
//  ========== SCAN ================================================  //

// scanRowsNodoHistoria escanea cada row en la struct NodoHistoria
// y devuelve un slice con todos los items.
// Siempre se encarga de llamar rows.Close()
func (s *Repositorio) scanRowsNodoHistoria(rows *sql.Rows, op string) ([]ust.NodoHistoria, error) {
	defer rows.Close()
	items := []ust.NodoHistoria{}
	for rows.Next() {
		nhist := ust.NodoHistoria{}
		err := rows.Scan(
			&nhist.HistoriaID, &nhist.ProyectoID, &nhist.PersonaID, &nhist.Titulo, &nhist.Objetivo, &nhist.Prioridad, &nhist.Completada, &nhist.PadreID, &nhist.PadreTbl, &nhist.Nivel, &nhist.Posicion, &nhist.SegundosPresupuesto, &nhist.NumHistorias, &nhist.NumTareas, &nhist.SegundosEstimado, &nhist.SegundosReal,
		)
		if err != nil {
			return nil, gko.ErrInesperado().Err(err).Op(op)
		}
		items = append(items, nhist)
	}
	return items, nil
}

//  ================================================================  //
//  ========== LIST ================================================  //

func (s *Repositorio) ListNodoHistorias() ([]ust.NodoHistoria, error) {
	const op string = "ListNodoHistorias"
	rows, err := s.db.Query(
		"SELECT " + columnasNodoHistoria + " " + fromNodoHistoria,
	)
	if err != nil {
		return nil, gko.ErrInesperado().Err(err).Op(op)
	}
	return s.scanRowsNodoHistoria(rows, op)
}

//  ================================================================  //
//  ========== LIST BYPROYECTOID ===================================  //

func (s *Repositorio) ListNodoHistoriasByProyectoID(ProyectoID string) ([]ust.NodoHistoria, error) {
	const op string = "ListNodoHistoriasByProyectoID"
	rows, err := s.db.Query(
		"SELECT "+columnasNodoHistoria+" "+fromNodoHistoria+
			"WHERE his.proyecto_id = ? ORDER BY (his.prioridad * nod.nivel) + 20 - nod.posicion DESC",
		ProyectoID,
	)
	if err != nil {
		return nil, gko.ErrInesperado().Err(err).Op(op)
	}
	return s.scanRowsNodoHistoria(rows, op)
}

//  ================================================================  //
//  ========== LIST BYPADREID ======================================  //

func (s *Repositorio) ListNodoHistoriasByPadreID(PadreID int) ([]ust.NodoHistoria, error) {
	const op string = "ListNodoHistoriasByPadreID"
	rows, err := s.db.Query(
		"SELECT "+columnasNodoHistoria+" "+fromNodoHistoria+
			"WHERE nod.padre_id = ? ORDER BY nod.posicion",
		PadreID,
	)
	if err != nil {
		return nil, gko.ErrInesperado().Err(err).Op(op)
	}
	return s.scanRowsNodoHistoria(rows, op)
}

//  ================================================================  //
//  ========== LIST PRIORITARIAS ===================================  //

func (s *Repositorio) ListNodoHistoriasPrioritarias() ([]ust.NodoHistoria, error) {
	const op string = "ListNodoHistoriasPrioritarias"
	rows, err := s.db.Query(
		"SELECT " + columnasNodoHistoria + " " + fromNodoHistoria +
			"WHERE his.prioridad > 0 AND completada == 0 ORDER BY (his.prioridad * nod.nivel) + 20 - nod.posicion DESC LIMIT 50",
	)
	if err != nil {
		return nil, gko.ErrInesperado().Err(err).Op(op)
	}
	return s.scanRowsNodoHistoria(rows, op)
}
