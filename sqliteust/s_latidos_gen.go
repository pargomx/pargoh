package sqliteust

import (
	"github.com/pargomx/gecko/gko"

	"monorepo/ust"
)

//  ================================================================  //
//  ========== INSERT ==============================================  //

func (s *Repositorio) InsertLatido(lat ust.Latido) error {
	const op string = "InsertLatido"
	if lat.Timestamp == "" {
		return gko.ErrDatoIndef.Str("pk_indefinida").Op(op).Msg("Timestamp sin especificar")
	}
	_, err := s.db.Exec("INSERT INTO latidos "+
		"(timestamp, persona_id, segundos) "+
		"VALUES (?, ?, ?) ",
		lat.Timestamp, lat.PersonaID, lat.Segundos,
	)
	if err != nil {
		return gko.ErrAlEscribir.Err(err).Op(op)
	}
	return nil
}
