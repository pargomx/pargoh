package sqliteust

import "github.com/pargomx/gecko/gko"

func (s *Repositorio) DeleteAllTareas(HistoriaID int) error {
	const op string = "DeleteAllTareas"
	if HistoriaID == 0 {
		return gko.ErrDatoIndef().Op(op).Msg("HistoriaID sin especificar").Str("pk_indefinida")
	}
	err := s.ExisteHistoria(HistoriaID)
	if err != nil {
		return gko.Err(err).Op(op)
	}
	_, err = s.db.Exec(
		"DELETE FROM tareas WHERE historia_id = ?",
		HistoriaID,
	)
	if err != nil {
		return gko.ErrAlEscribir().Err(err).Op(op)
	}
	return nil
}
