package sqliteust

import (
	"monorepo/ust"

	"github.com/pargomx/gecko/gko"
)

// Retorna los latidos agregados por minuto y persona.
//
//	Ej. desde '2025-03-02 00:00:00' hasta '2025-03-04 23:59:59'.
func (s *Repositorio) ListLatidos(desde, hasta string) ([]ust.Latido, error) {
	op := gko.Op("ListLatidos")
	if len(desde) != 19 {
		return nil, op.E(gko.ErrDatoInvalido).Strf("desde inválido: %v", desde)
	}
	if len(hasta) != 19 {
		return nil, op.E(gko.ErrDatoInvalido).Strf("hasta inválido: %v", hasta)
	}
	rows, err := s.db.Query(
		"SELECT substr(timestamp,0,17)||':00' AS minuto, persona_id, sum(segundos) as segundos FROM latidos WHERE minuto BETWEEN ? AND ? GROUP BY minuto, persona_id;",
		desde, hasta,
	)
	if err != nil {
		return nil, op.Err(err)
	}
	defer rows.Close()
	items := []ust.Latido{}
	for rows.Next() {
		lat := ust.Latido{}
		err := rows.Scan(
			&lat.Timestamp, &lat.PersonaID, &lat.Segundos,
		)
		if err != nil {
			return nil, op.Err(err)
		}
		items = append(items, lat)
	}
	return items, nil
}
