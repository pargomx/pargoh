package arbol

import "github.com/pargomx/gecko/gko"

const EvNodoMovido gko.EventKey = "nodo_movido"

type ArgsMover struct {
	NodoID     int // Nodo que se moverá.
	NewPadreID int // Nuevo padre.
}

func (s *AppTx) MoverHoja(args ArgsMover) error {
	op := gko.Op("MoverHoja")
	nod, err := s.repo.GetNodo(args.NodoID)
	if err != nil {
		return op.Err(err)
	}
	if args.NewPadreID == 0 {
		return op.Msg("No se especificó a dónde se moverá")
	}
	if nod.PadreID == args.NewPadreID {
		return op.Msg("No se moverá porque sigue siendo el mismo padre")
	}
	newPadre, err := s.repo.GetNodo(args.NewPadreID)
	if err != nil {
		return op.Err(err).Msg("El nuevo padre no existe")
	}

	// TODO: verificar que no sea hijo de su propio descendiente.
	// for _, ancestro := range newPadre.Ancestros {
	// 	if ancestro.HistoriaID == historiaID {
	// 		return op.Msg("La historia no puede ser hija de su propio descendiente")
	// 	}
	// }

	// TODO: verificar que el nuevo padre sea un tipo de entidad válida.

	err = s.repo.MoverNodo(*nod, newPadre.NodoID)
	if err != nil {
		return op.Err(err)
	}
	s.Results.Add(EvNodoMovido.WithArgs(args))
	return nil
}
