package arbol

import (
	"fmt"

	"github.com/pargomx/gecko/gko"
)

const EvNodoReordenado gko.EventKey = "nodo.reordenado"

type evNodoReordenado struct {
	NodoID int
	OldPos int
	NewPos int
}

func (e evNodoReordenado) ToMsg(t string) string {
	switch t {
	case "key":
		return string(EvNodoReordenado)
	default:
		return fmt.Sprintf("Nodo reposicionado %v de %v a %v", e.NodoID, e.OldPos, e.NewPos)
	}
}

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
	err = s.riseEvent(EvNodoReordenado, evNodoReordenado{
		NodoID: nod.NodoID,
		OldPos: nod.Posicion,
		NewPos: args.NewPos,
	})
	if err != nil {
		return op.Err(err)
	}
	return nil
}

// ================================================================ //
// ================================================================ //

const EvNodoMovido gko.EventKey = "nodo.movido"

type evNodoMovido struct {
	NodoID     int
	NewPadreID int
}

func (e evNodoMovido) ToMsg(t string) string {
	switch t {
	case "key":
		return string(EvNodoMovido)
	default:
		return fmt.Sprintf("Nodo movido: %v a nuevo padre: %v", e.NodoID, e.NewPadreID)
	}
}

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

	err = s.riseEvent(EvNodoMovido, evNodoMovido{
		NodoID:     args.NodoID,
		NewPadreID: args.NewPadreID,
	})
	if err != nil {
		return op.Err(err)
	}

	return nil
}
