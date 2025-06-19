package sqlitearbol

import (
	"github.com/pargomx/gecko/gko"

	"monorepo/arbol"
)

//  ================================================================  //
//  ========== INSERT ==============================================  //

func (s *Repositorio) InsertLatido(lat arbol.Latido) error {
	const op string = "InsertLatido"
	if lat.TsLatido == "" {
		return gko.ErrDatoIndef.Str("pk_indefinida").Op(op).Msg("TsLatido sin especificar")
	}
	_, err := s.db.Exec("INSERT INTO latidos "+
		"(ts_latido, nodo_id, segundos) "+
		"VALUES (?, ?, ?) ",
		lat.TsLatido, lat.NodoID, lat.Segundos,
	)
	if err != nil {
		return gko.ErrAlEscribir.Err(err).Op(op)
	}
	return nil
}
