package sqlitearbol

import (
	"database/sql"

	"github.com/pargomx/gecko/gko"

	"monorepo/arbol"
)

//  ================================================================  //
//  ========== INSERT ==============================================  //

func (s *Repositorio) InsertReferencia(ref arbol.Referencia) error {
	const op string = "InsertReferencia"
	if ref.NodoID == 0 {
		return gko.ErrDatoIndef.Str("pk_indefinida").Op(op).Msg("NodoID sin especificar")
	}
	if ref.RefNodoID == 0 {
		return gko.ErrDatoIndef.Str("pk_indefinida").Op(op).Msg("RefNodoID sin especificar")
	}
	_, err := s.db.Exec("INSERT INTO referencias "+
		"(nodo_id, ref_nodo_id) "+
		"VALUES (?, ?) ",
		ref.NodoID, ref.RefNodoID,
	)
	if err != nil {
		return gko.ErrAlEscribir.Err(err).Op(op)
	}
	return nil
}

//  ================================================================  //
//  ========== EXISTE ==============================================  //

// Retorna error nil si existe solo un registro con esta clave primaria.
func (s *Repositorio) ExisteReferencia(NodoID int, RefNodoID int) error {
	const op string = "ExisteReferencia"
	var num int
	err := s.db.QueryRow("SELECT COUNT(nodo_id) FROM referencias WHERE nodo_id = ? AND ref_nodo_id = ?",
		NodoID, RefNodoID,
	).Scan(&num)
	if err != nil {
		if err == sql.ErrNoRows {
			return gko.ErrNoEncontrado.Msg("Referencia no encontrado").Op(op)
		}
		return gko.ErrInesperado.Err(err).Op(op)
	}
	if num > 1 {
		return gko.ErrInesperado.Err(nil).Op(op).Str("existen más de un registro para la pk").Ctx("registros", num)
	} else if num == 0 {
		return gko.ErrNoEncontrado.Msg("Referencia no encontrado").Op(op)
	}
	return nil
}

//  ================================================================  //
//  ========== DELETE ==============================================  //

func (s *Repositorio) DeleteReferencia(NodoID int, RefNodoID int) error {
	const op string = "DeleteReferencia"
	if NodoID == 0 {
		return gko.ErrDatoIndef.Str("pk_indefinida").Op(op).Msg("NodoID sin especificar")
	}
	if RefNodoID == 0 {
		return gko.ErrDatoIndef.Str("pk_indefinida").Op(op).Msg("RefNodoID sin especificar")
	}
	err := s.ExisteReferencia(NodoID, RefNodoID)
	if err != nil {
		return gko.Err(err).Op(op)
	}
	_, err = s.db.Exec(
		"DELETE FROM referencias WHERE nodo_id = ? AND ref_nodo_id = ?",
		NodoID, RefNodoID,
	)
	if err != nil {
		return gko.ErrAlEscribir.Err(err).Op(op)
	}
	return nil
}
