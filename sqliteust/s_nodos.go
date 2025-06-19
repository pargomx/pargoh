package sqliteust

import (
	"monorepo/ust"
	"strings"

	"github.com/pargomx/gecko/gko"
)

// InsertNodo valida el registro y lo inserta en la base de datos.
func (s *Repositorio) InsertNodo(nod ust.Nodo) error {
	const op string = "sqliteust.InsertNodo"
	if nod.NodoID == 0 {
		return gko.ErrDatoInvalido.Msg("NodoID sin especificar").Ctx(op, "pk_indefinida")
	}
	if nod.NodoTbl == "" {
		return gko.ErrDatoInvalido.Msg("NodoTbl sin especificar").Ctx(op, "required_sin_valor")
	}
	if nod.PadreTbl == "" {
		return gko.ErrDatoInvalido.Msg("PadreTbl sin especificar").Ctx(op, "required_sin_valor")
	}
	err := nod.Validar()
	if err != nil {
		return gko.ErrDatoInvalido.Err(err).Op(op).Msg(err.Error())
	}
	_, err = s.db.Exec("INSERT INTO nodos "+
		"(nodo_id, nodo_tbl, padre_id, padre_tbl, nivel, posicion) "+
		"VALUES (?, ?, ?, ?, (SELECT nivel+1 FROM nodos WHERE nodo_id = ?), (SELECT count(nodo_id)+1 FROM nodos WHERE padre_id = ?)) ",
		nod.NodoID, nod.NodoTbl, nod.PadreID, nod.PadreTbl, nod.PadreID, nod.PadreID,
	)
	if err != nil {
		if strings.HasPrefix(err.Error(), "Error 1062 (23000)") {
			return gko.ErrYaExiste.Err(err).Op(op)
		} else if strings.HasPrefix(err.Error(), "Error 1452 (23000)") {
			return gko.ErrDatoInvalido.Err(err).Op(op).Msg("No se puede insertar la información porque el registro asociado no existe")
		} else {
			return gko.ErrInesperado.Err(err).Op(op)
		}
	}
	return nil
}

func (s *Repositorio) EliminarNodo(nodoID int) error {
	op := gko.Op("sqliteust.EliminarNodo")
	if nodoID == 0 {
		return op.Msg("NodoID sin especificar")
	}
	nodo, err := s.GetNodo(nodoID)
	if err != nil {
		return op.Err(err)
	}
	_, err = s.db.Exec("DELETE FROM nodos WHERE nodo_id = ?", nodoID)
	if err != nil {
		return op.Err(err)
	}
	_, err = s.db.Exec("UPDATE nodos SET posicion = posicion - 1 WHERE padre_id = ? AND posicion > ?", nodo.PadreID, nodo.Posicion)
	if err != nil {
		return op.Err(err)
	}
	return nil
}

func (s *Repositorio) MoverNodo(nodoID int, nuevoPadreID int) error {
	op := gko.Op("sqliteust.MoverNodo")
	nodoMovido, err := s.GetNodo(nodoID)
	if err != nil {
		return op.Err(err)
	}
	newPadre, err := s.GetNodo(nuevoPadreID)
	if err != nil {
		return op.Err(err)
	}
	_, err = s.db.Exec(
		"UPDATE nodos SET padre_id = ?, padre_tbl = ?, nivel = (SELECT nivel+1 FROM nodos WHERE nodo_id = ?), posicion = (SELECT count(nodo_id)+1 FROM nodos WHERE padre_id = ?) WHERE nodo_id = ?",
		newPadre.NodoID, newPadre.NodoTbl, newPadre.NodoID, newPadre.NodoID, nodoMovido.NodoID,
	)
	if err != nil {
		return op.Err(err).Op("update_nodo")
	}
	_, err = s.db.Exec(
		"UPDATE nodos SET posicion = posicion - 1 WHERE padre_id = ? AND posicion > ?", nodoMovido.PadreID, nodoMovido.Posicion,
	)
	if err != nil {
		return op.Err(err).Op("update_old_hermanos")
	}
	err = s.actualizarNivelDeDescendientes(nodoID)
	if err != nil {
		return op.Err(err).Op("update_nivel_descendientes")
	}
	return nil
}

func (s *Repositorio) actualizarNivelDeDescendientes(nodoID int) error {
	hijos, err := s.ListNodosByPadreID(nodoID)
	if err != nil {
		return err
	}
	if len(hijos) > 0 {
		_, err = s.db.Exec(
			"UPDATE nodos SET nivel = (SELECT nivel+1 FROM nodos WHERE nodo_id = ?) WHERE padre_id = ?",
			nodoID, nodoID,
		)
		if err != nil {
			return err
		}
		for _, hijo := range hijos {
			err = s.actualizarNivelDeDescendientes(hijo.NodoID)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
