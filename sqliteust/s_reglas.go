package sqliteust

import "github.com/pargomx/gecko/gko"

func (s *Repositorio) ReordenarRegla(historiaID int, oldPos int, newPos int) error {
	var op = gko.Op("sqliteust.ReordenarRegla").Ctx("hist", historiaID).Ctx("old", oldPos).Ctx("new", newPos)
	if oldPos == newPos {
		return nil
	}
	var hermanos int
	err := s.db.QueryRow("SELECT COUNT(*) FROM reglas WHERE historia_id = ?", historiaID).Scan(&hermanos)
	if err != nil {
		return op.Op("contar_hermanos").Err(err)
	}
	if newPos < 1 || newPos > hermanos {
		return op.Msgf("posición %v inválida para regla %v que tiene %v hermanos", newPos, historiaID, hermanos)
	}

	// Se utilizan posiciones negativas temporales porque no se puede confiar en el orden
	// con el que el update modifica las filas para que mantenga el UNIQUE con la posición.
	_, err = s.db.Exec(
		"UPDATE reglas SET posicion = -(?) WHERE historia_id = ? AND posicion = ?",
		newPos, historiaID, oldPos,
	)
	if err != nil {
		return op.Op("set_pos_negativa").Err(err)
	}

	// Dependiendo si se mueve hacia arriba o abajo, recorrer a los hermanos.
	if oldPos < newPos {
		_, err = s.db.Exec(
			"UPDATE reglas SET posicion = -(posicion - 1) WHERE historia_id = ? AND posicion > ? AND posicion <= ?",
			historiaID, oldPos, newPos,
		)
	} else {
		_, err = s.db.Exec(
			"UPDATE reglas SET posicion = -(posicion + 1) WHERE historia_id = ? AND posicion >= ? AND posicion <= ?",
			historiaID, newPos, oldPos,
		)
	}
	if err != nil {
		return op.Op("set_pos_hermanos").Err(err)
	}

	// Volver a positivo todas las posiciones.
	_, err = s.db.Exec(
		"UPDATE reglas SET posicion = -(posicion) WHERE historia_id = ? AND posicion < 0",
		historiaID,
	)
	if err != nil {
		return op.Op("set_pos_positiva").Err(err)
	}
	return nil
}
