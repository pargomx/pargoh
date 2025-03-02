package sqliteust

import (
	"monorepo/ust"

	"github.com/pargomx/gecko/gko"
)

func (s *Repositorio) FullTextSearch(search string) ([]ust.SearchResult, error) {
	const op string = "FullTextSearch"
	if search == "" {
		return nil, gko.ErrDatoIndef().Msg("Búsqueda vacía").Op(op)
	}
	rows, err := s.db.Query(
		"SELECT historia_id, otro_id, origen, snippet(historias_fts, 3, 'ſ', 'ſ', '...', 64) "+
			"FROM historias_fts WHERE historias_fts MATCH ? ORDER BY rank LIMIT 150",
		search,
	)
	if err != nil {
		return nil, gko.ErrInesperado().Err(err).Op(op)
	}
	defer rows.Close()
	items := []ust.SearchResult{}
	for rows.Next() {
		item := ust.SearchResult{}
		err := rows.Scan(
			&item.HistoriaID, &item.OtroID, &item.Origen, &item.Texto,
		)
		if err != nil {
			return nil, gko.ErrInesperado().Err(err).Op(op)
		}
		items = append(items, item)
	}
	return items, nil
}
