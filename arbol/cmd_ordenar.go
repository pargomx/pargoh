package arbol

import "github.com/pargomx/gecko/gko"

type ArgsReordenar struct {
	NodoID int // Nodo a reordenar respecto a sus hermanos
	NewPos int // Nueva posición: desde 1 hasta el número total de hermanos.
}

func (s *AppTx) ReordenarEntidad(args ArgsReordenar) error {
	op := gko.Op("ReordenarEntidad")
	nod, err := s.repo.GetNodo(args.NodoID)
	if err != nil {
		return op.Err(err)
	}
	err = s.repo.ReordenarNodo(*nod, args.NewPos)
	if err != nil {
		return op.Err(err)
	}
	s.Results.Add(EvNodoReordenado.WithArgs(args))
	// .Msgf("Reordenado %v de %v a %v", nod.NodoID, nod.Posicion, args.NewPos))
	return nil
}
