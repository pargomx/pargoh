package sqliteust

import (
	"monorepo/historias_de_usuario/ust"
	"strings"

	"github.com/pargomx/gecko/gko"
)

// InsertNodo valida el registro y lo inserta en la base de datos.
func (s *Repositorio) InsertNodo(nod ust.Nodo) error {
	const op string = "sqliteust.InsertNodo"
	if nod.NodoID == 0 {
		return gko.ErrDatoInvalido().Msg("NodoID sin especificar").Ctx(op, "pk_indefinida")
	}
	if nod.NodoTbl == "" {
		return gko.ErrDatoInvalido().Msg("NodoTbl sin especificar").Ctx(op, "required_sin_valor")
	}
	if nod.PadreTbl == "" {
		return gko.ErrDatoInvalido().Msg("PadreTbl sin especificar").Ctx(op, "required_sin_valor")
	}
	err := nod.Validar()
	if err != nil {
		return gko.ErrDatoInvalido().Err(err).Op(op).Msg(err.Error())
	}
	_, err = s.db.Exec("INSERT INTO nodos "+
		"(nodo_id, nodo_tbl, padre_id, padre_tbl, nivel, posicion) "+
		"VALUES (?, ?, ?, ?, (SELECT nivel+1 FROM nodos WHERE nodo_id = ?), (SELECT count(nodo_id)+1 FROM nodos WHERE padre_id = ?)) ",
		nod.NodoID, nod.NodoTbl, nod.PadreID, nod.PadreTbl, nod.PadreID, nod.PadreID,
	)
	if err != nil {
		if strings.HasPrefix(err.Error(), "Error 1062 (23000)") {
			return gko.ErrYaExiste().Err(err).Op(op)
		} else if strings.HasPrefix(err.Error(), "Error 1452 (23000)") {
			return gko.ErrDatoInvalido().Err(err).Op(op).Msg("No se puede insertar la información porque el registro asociado no existe")
		} else {
			return gko.ErrInesperado().Err(err).Op(op)
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
	nodo, err := s.GetNodo(nodoID)
	if err != nil {
		return op.Err(err)
	}
	_, err = s.db.Exec(
		"UPDATE nodos SET padre_id = ?, nivel = (SELECT nivel+1 FROM nodos WHERE nodo_id = ?), posicion = (SELECT count(nodo_id)+1 FROM nodos WHERE padre_id = ?) WHERE nodo_id = ?",
		nuevoPadreID, nuevoPadreID, nuevoPadreID, nodoID,
	)
	if err != nil {
		return op.Err(err).Op("update_nodo")
	}
	_, err = s.db.Exec(
		"UPDATE nodos SET posicion = posicion - 1 WHERE padre_id = ? AND posicion > ?", nodo.PadreID, nodo.Posicion,
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

// La cambia de posición respecto a sus hermanos, o sea en el mismo padre.
func (s *Repositorio) ReordenarNodo(nodoID int, oldPosicion int, newPosicion int) error {
	var op = gko.Op("sqliteust.ReordenarNodo").Ctx("id", nodoID).Ctx("old", oldPosicion).Ctx("new", newPosicion)
	if oldPosicion == newPosicion {
		return nil
	}

	var padreID int
	err := s.db.QueryRow("SELECT padre_id FROM nodos WHERE nodo_id = ?", nodoID).Scan(&padreID)
	if err != nil {
		return op.Op("get_padreID").Err(err)
	}

	var hermanos int
	err = s.db.QueryRow("SELECT COUNT(nodo_id) FROM nodos WHERE padre_id = ?", padreID).Scan(&hermanos)
	if err != nil {
		return op.Op("contar_hermanos").Err(err)
	}
	if newPosicion < 1 || newPosicion > hermanos {
		return op.Msgf("posición %v inválida para nodo %v que tiene %v hermanos", newPosicion, nodoID, hermanos)
	}

	// Se utilizan posiciones negativas temporales porque no se puede confiar en el orden
	// con el que el update modifica las filas para que mantenga el UNIQUE con la posición.
	_, err = s.db.Exec(
		"UPDATE nodos SET posicion = -(?) WHERE nodo_id = ?",
		newPosicion, nodoID,
	)
	if err != nil {
		return op.Op("set_pos_negativa").Err(err)
	}

	// Dependiendo si se mueve hacia arriba o abajo, recorrer a los hermanos.
	if oldPosicion < newPosicion {
		_, err = s.db.Exec(
			"UPDATE nodos SET posicion = -(posicion - 1) WHERE padre_id = ? AND posicion > ? AND posicion <= ?",
			padreID, oldPosicion, newPosicion,
		)
	} else {
		_, err = s.db.Exec(
			"UPDATE nodos SET posicion = -(posicion + 1) WHERE padre_id = ? AND posicion >= ? AND posicion <= ?",
			padreID, newPosicion, oldPosicion,
		)
	}
	if err != nil {
		return op.Op("set_pos_hermanos").Err(err)
	}

	// Volver a positivo todas las posiciones.
	_, err = s.db.Exec(
		"UPDATE nodos SET posicion = -(posicion) WHERE posicion < 0",
	)
	if err != nil {
		return op.Op("set_pos_positiva").Err(err)
	}
	return nil
}
