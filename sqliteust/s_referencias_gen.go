package sqliteust

import (
	"database/sql"

	"github.com/pargomx/gecko/gko"

	"monorepo/ust"
)

//  ================================================================  //
//  ========== INSERT ==============================================  //

func (s *Repositorio) InsertReferencia(ref ust.Referencia) error {
	const op string = "InsertReferencia"
	if ref.HistoriaID == 0 {
		return gko.ErrDatoIndef().Op(op).Msg("HistoriaID sin especificar").Str("pk_indefinida")
	}
	if ref.RefHistoriaID == 0 {
		return gko.ErrDatoIndef().Op(op).Msg("RefHistoriaID sin especificar").Str("pk_indefinida")
	}
	_, err := s.db.Exec("INSERT INTO referencias "+
		"(historia_id, ref_historia_id) "+
		"VALUES (?, ?) ",
		ref.HistoriaID, ref.RefHistoriaID,
	)
	if err != nil {
		return gko.ErrAlEscribir().Err(err).Op(op)
	}
	return nil
}

//  ================================================================  //
//  ========== EXISTE ==============================================  //

// Retorna error nil si existe solo un registro con esta clave primaria.
func (s *Repositorio) ExisteReferencia(HistoriaID int, RefHistoriaID int) error {
	const op string = "ExisteReferencia"
	var num int
	err := s.db.QueryRow("SELECT COUNT(historia_id) FROM referencias WHERE historia_id = ? AND ref_historia_id = ?",
		HistoriaID, RefHistoriaID,
	).Scan(&num)
	if err != nil {
		if err == sql.ErrNoRows {
			return gko.ErrNoEncontrado().Err(ust.ErrReferenciaNotFound).Op(op)
		}
		return gko.ErrInesperado().Err(err).Op(op)
	}
	if num > 1 {
		return gko.ErrInesperado().Err(nil).Op(op).Str("existen más de un registro para la pk").Ctx("registros", num)
	} else if num == 0 {
		return gko.ErrNoEncontrado().Err(ust.ErrReferenciaNotFound).Op(op)
	}
	return nil
}

//  ================================================================  //
//  ========== DELETE ==============================================  //

func (s *Repositorio) DeleteReferencia(HistoriaID int, RefHistoriaID int) error {
	const op string = "DeleteReferencia"
	if HistoriaID == 0 {
		return gko.ErrDatoIndef().Op(op).Msg("HistoriaID sin especificar").Str("pk_indefinida")
	}
	if RefHistoriaID == 0 {
		return gko.ErrDatoIndef().Op(op).Msg("RefHistoriaID sin especificar").Str("pk_indefinida")
	}
	err := s.ExisteReferencia(HistoriaID, RefHistoriaID)
	if err != nil {
		return gko.Err(err).Op(op)
	}
	_, err = s.db.Exec(
		"DELETE FROM referencias WHERE historia_id = ? AND ref_historia_id = ?",
		HistoriaID, RefHistoriaID,
	)
	if err != nil {
		return gko.ErrAlEscribir().Err(err).Op(op)
	}
	return nil
}
