package sqlitearbol

import (
	"monorepo/arbol"

	"github.com/pargomx/gecko/gko"
)

// Retorna los latidos agregados por minuto y nodo.
//
//	Ej. desde '2025-03-02 00:00:00' hasta '2025-03-04 23:59:59'.
func (s *Repositorio) ListLatidos(desde, hasta string) ([]arbol.Latido, error) {
	op := gko.Op("ListLatidos")
	if len(desde) != 19 {
		return nil, op.E(gko.ErrDatoInvalido).Strf("desde inválido: %v", desde)
	}
	if len(hasta) != 19 {
		return nil, op.E(gko.ErrDatoInvalido).Strf("hasta inválido: %v", hasta)
	}
	rows, err := s.db.Query(
		"SELECT substr(ts_latido,0,17)||':00' AS minuto, nodo_id, sum(segundos) as segundos FROM latidos WHERE minuto BETWEEN ? AND ? GROUP BY minuto, nodo_id;",
		desde, hasta,
	)
	if err != nil {
		return nil, op.Err(err)
	}
	defer rows.Close()
	items := []arbol.Latido{}
	for rows.Next() {
		lat := arbol.Latido{}
		err := rows.Scan(
			&lat.TsLatido, &lat.NodoID, &lat.Segundos,
		)
		if err != nil {
			return nil, op.Err(err)
		}
		items = append(items, lat)
	}
	return items, nil
}
