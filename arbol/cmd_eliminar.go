package arbol

import (
	"fmt"

	"github.com/pargomx/gecko/gko"
)

const EvNodoEliminado gko.EventKey = "nodo.eliminado"

type evNodoEliminado struct {
	NodoID int
	Tipo   string
}

func (e evNodoEliminado) ToMsg(t string) string {
	switch t {
	case "key":
		return string(EvNodoEliminado)
	default:
		return fmt.Sprintf("Nodo eliminado: %v (tipo: %s)", e.NodoID, e.Tipo)
	}
}

func (s *AppTx) EliminarNodo(NodoID int) (padre *Nodo, err error) {
	op := gko.Op("EliminarNodo")

	nod, err := s.repo.GetNodo(NodoID)
	if err != nil {
		return nil, op.Err(err)
	}

	hijos, err := s.repo.ListNodosByPadreID(nod.NodoID)
	if err != nil {
		return nil, op.Err(err)
	}

	if len(hijos) != 0 {
		return nil, op.E(gko.ErrHayHuerfanos).
			Msg("Para eliminar este nodo primero elimine sus descendientes")
	}

	err = s.repo.DeleteNodo(nod.NodoID)
	if err != nil {
		return nil, op.Err(err)
	}

	err = s.riseEvent(EvNodoEliminado, evNodoEliminado{
		NodoID: nod.NodoID,
		Tipo:   nod.Tipo,
	})
	if err != nil {
		return nil, op.Err(err)
	}

	// Eliminar imagen si la hay
	if nod.Imagen != "" {
		err = s.borrarImagen(argsBorrarImagen{Filename: nod.Imagen})
		if err != nil {
			return nil, op.Err(err)
		}
	}

	return s.repo.GetNodo(nod.PadreID)
}

// ================================================================ //

// Eliminar rama desde el nodo especificado junto con todos sus descendientes.
func (s *AppTx) EliminarRama(NodoID int) (padre *Nodo, err error) {
	op := gko.Op("EliminarRama")

	nod, err := s.repo.GetNodo(NodoID)
	if err != nil {
		return nil, op.Err(err)
	}

	err = s.eliminarRecursivo(*nod)
	if err != nil {
		return nil, op.Err(err)
	}

	return s.repo.GetNodo(nod.PadreID)
}

// Elimina desde el último descendiente de forma recursiva.
func (s *AppTx) eliminarRecursivo(nod Nodo) (err error) {

	hijos, err := s.repo.ListNodosByPadreID(nod.NodoID)
	if err != nil {
		return err
	}
	for _, hijo := range hijos {
		err = s.eliminarRecursivo(hijo)
		if err != nil {
			return err
		}
	}

	err = s.repo.DeleteNodo(nod.NodoID)
	if err != nil {
		return err
	}
	err = s.riseEvent(EvNodoEliminado, evNodoEliminado{
		NodoID: nod.NodoID,
		Tipo:   nod.Tipo,
	})
	if nod.Imagen != "" {
		err = s.borrarImagen(argsBorrarImagen{Filename: nod.Imagen})
		if err != nil {
			return err
		}
	}
	return nil
}
