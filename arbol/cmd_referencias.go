package arbol

import (
	"github.com/pargomx/gecko/gko"
)

const EvReferenciaAgregada gko.EventKey = "referencia_agregada"
const EvReferenciaEliminada gko.EventKey = "referencia_eliminada"

type ArgsReferencia struct {
	NodoID    int
	RefNodoID int
}

func (s *AppTx) AgregarReferencia(args ArgsReferencia) error {
	op := gko.Op("AgregarReferencia")
	if args.NodoID == 0 {
		return op.Str("falta NodoID")
	}
	if args.RefNodoID == 0 {
		return op.Str("falta RefNodoID")
	}
	if args.NodoID == args.RefNodoID {
		return op.Msg("No se puede referenciar a sí misma")
	}
	err := s.repo.ExisteNodo(args.RefNodoID)
	if err != nil {
		return op.Err(err)
	}
	err = s.repo.ExisteNodo(args.RefNodoID)
	if err != nil {
		return op.Err(err)
	}
	// No debe existir la referencia en ningún sentido.
	if err = s.repo.ExisteReferencia(args.NodoID, args.RefNodoID); err == nil {
		return op.Msg("Ya existe esta referencia")
	}
	if err = s.repo.ExisteReferencia(args.RefNodoID, args.NodoID); err == nil {
		return op.Msg("Ya existe esta referencia en sentido contrario")
	}
	err = s.repo.InsertReferencia(Referencia(args))
	if err != nil {
		return op.Err(err)
	}
	s.Results.Add(EvReferenciaAgregada.WithArgs(args))
	return nil
}

func (s *AppTx) EliminarReferencia(args ArgsReferencia) error {
	op := gko.Op("EliminarReferencia")
	if args.NodoID == 0 {
		return op.Str("falta NodoID")
	}
	if args.RefNodoID == 0 {
		return op.Str("falta RefNodoID")
	}
	err := s.repo.ExisteReferencia(args.NodoID, args.RefNodoID)
	if err != nil { // Puede que la referencia esté volteada.
		nodoID := args.NodoID
		args.NodoID = args.RefNodoID
		args.RefNodoID = nodoID
	}
	err = s.repo.DeleteReferencia(args.NodoID, args.RefNodoID)
	if err != nil {
		return op.Err(err)
	}
	s.Results.Add(EvReferenciaEliminada.WithArgs(args))
	return nil
}
