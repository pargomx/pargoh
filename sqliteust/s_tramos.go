package sqliteust

import "github.com/pargomx/gecko/gko"

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
