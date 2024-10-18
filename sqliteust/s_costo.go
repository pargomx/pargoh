package sqliteust

import (
	"monorepo/ust"

	"github.com/pargomx/gecko/gko"

	_ "embed"
)

//go:embed qry_costo.sql
var qryHistoriasCosto string

func (s *Repositorio) ListHistoriasCosto(personaID int) ([]ust.HistoriaCosto, error) {
	const op string = "ListHistoriasCosto"
	rows, err := s.db.Query(qryHistoriasCosto, personaID)
	if err != nil {
		return nil, gko.ErrInesperado().Err(err).Op(op)
	}
	defer rows.Close()
	items := []ust.HistoriaCosto{}
	for rows.Next() {
		item := ust.HistoriaCosto{}
		err := rows.Scan(
			&item.HistoriaID,
			&item.PadreID,
			&item.Nivel,
			&item.Posicion,
			&item.Titulo,
			&item.Prioridad,
			&item.Completada,
			&item.SegundosEstimado,
			&item.SegundosUtilizado,
		)
		if err != nil {
			return nil, gko.ErrInesperado().Err(err).Op(op)
		}
		items = append(items, item)
	}
	return items, nil
}

// ================================================================ //
// ========== DIAS ================================================ //

const qryDias = `WITH RECURSIVE date_range AS (
    SELECT (SELECT date(min(inicio), "-5 hours") FROM intervalos) AS dia
    UNION ALL
    SELECT date(dia, '+1 day')
    FROM date_range
    WHERE dia < date('now')
)
SELECT dia FROM date_range;`

func (s *Repositorio) ListDias() ([]string, error) {
	const op string = "ListDias"
	rows, err := s.db.Query(qryDias)
	if err != nil {
		return nil, gko.ErrInesperado().Err(err).Op(op)
	}
	defer rows.Close()
	items := []string{}
	for rows.Next() {
		var item string
		err := rows.Scan(&item)
		if err != nil {
			return nil, gko.ErrInesperado().Err(err).Op(op)
		}
		items = append(items, item)
	}
	return items, nil
}
