package sqliteust

import (
	"database/sql"
	"errors"

	"github.com/pargomx/gecko/gko"

	"monorepo/historias_de_usuario/ust"
)

//  ================================================================  //
//  ========== MYSQL/CONSTANTES ====================================  //

// Lista de columnas separadas por coma para usar en consulta SELECT
// en conjunto con scanRow o scanRows, ya que las columnas coinciden
// con los campos escaneados.
const columnasNodoHistoria string = "his.historia_id, his.titulo, his.objetivo, his.prioridad, his.completada, coalesce(nod.padre_id, 0), coalesce(nod.padre_tbl, ''), coalesce(nod.nivel, 0), coalesce(nod.posicion, 0), coalesce((SELECT COUNT(nodo_id) FROM nodos WHERE padre_id = his.historia_id), 0) AS num_historias, coalesce((SELECT COUNT(tarea_id) FROM tareas WHERE historia_id = his.historia_id), 0) AS num_tareas"

// Origen de los datos de ust.NodoHistoria
//
// FROM historias his
// INNER JOIN nodos nod ON nodo_id = historia_id
const fromNodoHistoria string = "FROM historias his INNER JOIN nodos nod ON nodo_id = historia_id "

//  ================================================================  //
//  ========== MYSQL/SCAN-ROW ======================================  //

// Utilizar luego de un sql.QueryRow(). No es necesario hacer row.Close()
func (s *Repositorio) scanRowNodoHistoria(row *sql.Row, nhist *ust.NodoHistoria, op string) error {

	err := row.Scan(
		&nhist.HistoriaID, &nhist.Titulo, &nhist.Objetivo, &nhist.Prioridad, &nhist.Completada, &nhist.PadreID, &nhist.PadreTbl, &nhist.Nivel, &nhist.Posicion, &nhist.NumHistorias, &nhist.NumTareas,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return gko.ErrNoEncontrado().Msg("Historia de usuario no se encuentra").Op(op)
		}
		return gko.ErrInesperado().Err(err).Op(op)
	}

	return nil
}

//  ================================================================  //
//  ========== MYSQL/GET ===========================================  //

// GetNodoHistoria devuelve un NodoHistoria de la DB por su clave primaria.
// Error si no encuentra ninguno, o si encuentra más de uno.
func (s *Repositorio) GetNodoHistoria(HistoriaID int) (*ust.NodoHistoria, error) {
	const op string = "mysqlust.GetNodoHistoria"
	if HistoriaID == 0 {
		return nil, gko.ErrDatoInvalido().Msg("HistoriaID sin especificar").Ctx(op, "pk_indefinida")
	}
	row := s.db.QueryRow(
		"SELECT "+columnasNodoHistoria+" "+fromNodoHistoria+
			"WHERE his.historia_id = ?",
		HistoriaID,
	)
	nhist := &ust.NodoHistoria{}
	return nhist, s.scanRowNodoHistoria(row, nhist, op)
}

//  ================================================================  //
//  ========== MYSQL/SCAN-ROWS =====================================  //

// scanRowsNodoHistoria escanea cada row en la struct NodoHistoria
// y devuelve un slice con todos los items.
// Siempre se encarga de llamar rows.Close()
func (s *Repositorio) scanRowsNodoHistoria(rows *sql.Rows, op string) ([]ust.NodoHistoria, error) {
	defer rows.Close()
	items := []ust.NodoHistoria{}
	for rows.Next() {
		nhist := ust.NodoHistoria{}

		err := rows.Scan(
			&nhist.HistoriaID, &nhist.Titulo, &nhist.Objetivo, &nhist.Prioridad, &nhist.Completada, &nhist.PadreID, &nhist.PadreTbl, &nhist.Nivel, &nhist.Posicion, &nhist.NumHistorias, &nhist.NumTareas,
		)
		if err != nil {
			return nil, gko.ErrInesperado().Err(err).Op(op)
		}

		items = append(items, nhist)
	}
	return items, nil
}

//  ================================================================  //
//  ========== MYSQL/LIST_BY =======================================  //

func (s *Repositorio) ListNodoHistoriasByPadreID(PadreID int) ([]ust.NodoHistoria, error) {
	const op string = "mysqlust.ListNodoHistoriasByPadreID"
	if PadreID == 0 {
		return nil, gko.ErrDatoInvalido().Msg("PadreID sin especificar").Ctx(op, "param_indefinido")
	}
	rows, err := s.db.Query(
		"SELECT "+columnasNodoHistoria+" "+fromNodoHistoria+
			"WHERE nod.padre_id = ?"+" ORDER BY nod.posicion",
		PadreID,
	)
	if err != nil {
		return nil, gko.ErrInesperado().Err(err).Op(op)
	}
	return s.scanRowsNodoHistoria(rows, op)
}
