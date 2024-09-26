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
			&item.MinutosEstimado,
			&item.SegundosReal,
		)
		if err != nil {
			return nil, gko.ErrInesperado().Err(err).Op(op)
		}
		items = append(items, item)
	}
	return items, nil
}
