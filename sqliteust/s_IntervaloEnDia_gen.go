package sqliteust

import (
	"database/sql"

	"github.com/pargomx/gecko/gko"

	"monorepo/ust"
)

//  ================================================================  //
//  ========== CONSTANTES ==========================================  //

// Lista de columnas separadas por coma para usar en consulta SELECT
// en conjunto con scanRow o scanRows, ya que las columnas coinciden
// con los campos escaneados.
//
//	coalesce(his.proyecto_id, ''),
//	coalesce(his.persona_id, 0),
//	coalesce(his.historia_id, 0),
//	interv.tarea_id,
//	interv.inicio,
//	interv.fin,
//	coalesce(date(interv.inicio,'-5 hours'), '') AS fecha,
//	coalesce(unixepoch(interv.fin) - unixepoch(interv.inicio), 0) AS segundos
const columnasIntervaloEnDia string = "coalesce(his.proyecto_id, ''), coalesce(his.persona_id, 0), coalesce(his.historia_id, 0), interv.tarea_id, interv.inicio, interv.fin, coalesce(date(interv.inicio,'-5 hours'), '') AS fecha, coalesce(unixepoch(interv.fin) - unixepoch(interv.inicio), 0) AS segundos"

// Origen de los datos de ust.IntervaloEnDia
//
//	FROM intervalos interv
//	INNER JOIN tareas tar ON tar.tarea_id = interv.tarea_id
//	INNER JOIN historias his ON his.historia_id = tar.historia_id
const fromIntervaloEnDia string = "FROM intervalos interv INNER JOIN tareas tar ON tar.tarea_id = interv.tarea_id INNER JOIN historias his ON his.historia_id = tar.historia_id "

//  ================================================================  //
//  ========== SCAN ================================================  //

// scanRowsIntervaloEnDia escanea cada row en la struct IntervaloEnDia
// y devuelve un slice con todos los items.
// Siempre se encarga de llamar rows.Close()
func (s *Repositorio) scanRowsIntervaloEnDia(rows *sql.Rows, op string) ([]ust.IntervaloEnDia, error) {
	defer rows.Close()
	items := []ust.IntervaloEnDia{}
	for rows.Next() {
		intervd := ust.IntervaloEnDia{}
		err := rows.Scan(
			&intervd.ProyectoID, &intervd.PersonaID, &intervd.HistoriaID, &intervd.TareaID, &intervd.Inicio, &intervd.Fin, &intervd.Fecha, &intervd.Segundos,
		)
		if err != nil {
			return nil, gko.ErrInesperado().Err(err).Op(op)
		}
		items = append(items, intervd)
	}
	return items, nil
}

//  ================================================================  //
//  ========== LIST ================================================  //

func (s *Repositorio) ListIntervalosEnDias() ([]ust.IntervaloEnDia, error) {
	const op string = "ListIntervalosEnDias"
	rows, err := s.db.Query(
		"SELECT " + columnasIntervaloEnDia + " " + fromIntervaloEnDia,
	)
	if err != nil {
		return nil, gko.ErrInesperado().Err(err).Op(op)
	}
	return s.scanRowsIntervaloEnDia(rows, op)
}

//  ================================================================  //
//  ========== LIST ENTRE ==========================================  //

func (s *Repositorio) ListIntervalosEnDiasEntre(desde string, hasta string) ([]ust.IntervaloEnDia, error) {
	const op string = "ListIntervalosEnDiasEntre"
	rows, err := s.db.Query(
		"SELECT "+columnasIntervaloEnDia+" "+fromIntervaloEnDia+
			"WHERE workday BETWEEN ? AND ? ORDER BY interv.inicio",
		desde, hasta,
	)
	if err != nil {
		return nil, gko.ErrInesperado().Err(err).Op(op)
	}
	return s.scanRowsIntervaloEnDia(rows, op)
}
