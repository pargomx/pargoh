package arbol

import "github.com/pargomx/gecko/gko"

const EvEliminarNodo gko.EventKey = "nodo_eliminado"

type ArgsEliminarNodo struct {
	NodoID int // Nuevo ID aleatorio.
}

func (s *AppTx) EliminarNodo(args ArgsEliminarNodo) (padre *Nodo, err error) {
	op := gko.Op("EliminarNodo")

	nod, err := s.repo.GetNodo(args.NodoID)
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
	s.Results.Add(EvEliminarNodo.WithArgs(args).
		Msgf("Deleted %v %v", nod.Tipo, nod.NodoID))

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
func (s *AppTx) EliminarRama(args ArgsEliminarNodo) (padre *Nodo, err error) {
	op := gko.Op("EliminarRama")

	nod, err := s.repo.GetNodo(args.NodoID)
	if err != nil {
		return nil, op.Err(err)
	}

	err = s.eliminarRecursivo(*nod)
	if err != nil {
		return nil, op.Err(err)
	}
	s.Results.Add(EvEliminarNodo.WithArgs(args).
		Msgf("Deleted %v %v (rama completa)", nod.Tipo, nod.NodoID))

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
		s.Results.Add(EvEliminarNodo.WithArgs(ArgsEliminarNodo{NodoID: hijo.NodoID}).
			Msgf("Deleted %v %v in order to delete %v", hijo.Tipo, hijo.NodoID, nod.NodoID))
	}

	err = s.repo.DeleteNodo(nod.NodoID)
	if err != nil {
		return err
	}
	if nod.Imagen != "" {
		err = s.borrarImagen(argsBorrarImagen{Filename: nod.Imagen})
		if err != nil {
			return err
		}
	}
	return nil
}
