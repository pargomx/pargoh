package sqliteust

import "github.com/pargomx/gecko/gko"

func (s *Repositorio) ReordenarTramo(historiaID int, oldPos int, newPos int) error {
	var op = gko.Op("sqliteust.ReordenarTramo").Ctx("hist", historiaID).Ctx("old", oldPos).Ctx("new", newPos)
	if oldPos == newPos {
		return nil
	}
	var hermanos int
	err := s.db.QueryRow("SELECT COUNT(*) FROM tramos WHERE historia_id = ?", historiaID).Scan(&hermanos)
	if err != nil {
		return op.Op("contar_hermanos").Err(err)
	}
	if newPos < 1 || newPos > hermanos {
		return op.Msgf("posición %v inválida para tramo %v que tiene %v hermanos", newPos, historiaID, hermanos)
	}

	// Se utilizan posiciones negativas temporales porque no se puede confiar en el orden
	// con el que el update modifica las filas para que mantenga el UNIQUE con la posición.
	_, err = s.db.Exec(
		"UPDATE tramos SET posicion = -(?) WHERE historia_id = ? AND posicion = ?",
		newPos, historiaID, oldPos,
	)
	if err != nil {
		return op.Op("set_pos_negativa").Err(err)
	}

	// Dependiendo si se mueve hacia arriba o abajo, recorrer a los hermanos.
	if oldPos < newPos {
		_, err = s.db.Exec(
			"UPDATE tramos SET posicion = -(posicion - 1) WHERE historia_id = ? AND posicion > ? AND posicion <= ?",
			historiaID, oldPos, newPos,
		)
	} else {
		_, err = s.db.Exec(
			"UPDATE tramos SET posicion = -(posicion + 1) WHERE historia_id = ? AND posicion >= ? AND posicion <= ?",
			historiaID, newPos, oldPos,
		)
	}
	if err != nil {
		return op.Op("set_pos_hermanos").Err(err)
	}

	// Volver a positivo todas las posiciones.
	_, err = s.db.Exec(
		"UPDATE tramos SET posicion = -(posicion) WHERE historia_id = ? AND posicion < 0",
		historiaID,
	)
	if err != nil {
		return op.Op("set_pos_positiva").Err(err)
	}
	return nil
}

// ================================================================ //

func (s *Repositorio) MoverTramo(historiaID int, posicion int, newHistoriaID int) error {
	const op string = "MoverTramo"
	if historiaID == 0 {
		return gko.ErrDatoIndef.Msg("HistoriaID sin especificar").Str("pk_indefinida").Op(op)
	}
	if posicion == 0 {
		return gko.ErrDatoIndef.Msg("Posicion sin especificar").Str("pk_indefinida").Op(op)
	}
	err := s.ExisteTramo(historiaID, posicion)
	if err != nil {
		return gko.Err(err).Op(op)
	}
	_, err = s.db.Exec(
		"UPDATE tramos SET historia_id = ?, posicion = (SELECT COUNT(*) FROM tramos WHERE historia_id = ?)+1 WHERE historia_id = ? AND posicion = ?",
		newHistoriaID, newHistoriaID, historiaID, posicion,
	)
	if err != nil {
		return gko.ErrAlEscribir.Err(err).Op(op).Ctx("historia_id", historiaID).Ctx("Pos", posicion)
	}

	// Recorrer los tramos de la historia origen.
	// Se utilizan posiciones negativas temporales porque no se puede confiar en el orden
	// con el que el update modifica las filas para que mantenga el UNIQUE con la posición.
	_, err = s.db.Exec(
		"UPDATE tramos SET posicion = -(posicion - 1) WHERE historia_id = ? AND posicion > ?",
		historiaID, posicion,
	)
	if err != nil {
		return gko.ErrAlEscribir.Err(err).Op(op).Op("set_pos_negativa")
	}
	_, err = s.db.Exec(
		"UPDATE tramos SET posicion = -(posicion) WHERE historia_id = ? AND posicion < 0",
		historiaID,
	)
	if err != nil {
		return gko.ErrAlEscribir.Err(err).Op(op).Ctx("historia_id", historiaID).Ctx("Pos", posicion)
	}
	return nil
}

// ================================================================ //

func (s *Repositorio) DeleteAllTramos(HistoriaID int) error {
	const op string = "DeleteAllTramos"
	if HistoriaID == 0 {
		return gko.ErrDatoIndef.Msg("HistoriaID sin especificar").Str("pk_indefinida").Op(op)
	}
	err := s.ExisteHistoria(HistoriaID)
	if err != nil {
		return gko.Err(err).Op(op)
	}
	_, err = s.db.Exec(
		"DELETE FROM tramos WHERE historia_id = ?",
		HistoriaID,
	)
	if err != nil {
		return gko.ErrAlEscribir.Err(err).Op(op).Ctx("historia_id", HistoriaID)
	}
	return nil
}
