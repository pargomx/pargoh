package arbol

import (
	"fmt"

	"monorepo/ust"

	"github.com/pargomx/gecko/gko"
	"github.com/pargomx/gecko/gkt"
)

const EvNodoAgregado gko.EventKey = "nodo.agregado"

type evNodoAgregado struct {
	NodoID  int
	PadreID int
	Tipo    string
	Titulo  string
}

func (e evNodoAgregado) ToMsg(t string) string {
	switch t {
	case "key":
		return string(EvNodoAgregado)
	default:
		return fmt.Sprintf("Padre: %v ID: %v Tipo: '%s' Título: '%s'", e.PadreID, e.NodoID, e.Tipo, e.Titulo)
	}
}

type ArgsAgregarHoja struct {
	Tipo    string
	NodoID  int
	PadreID int
	Titulo  string
}

func (s *AppTx) AgregarHoja(args ArgsAgregarHoja) error {
	op := gko.Op("AgregarHoja")

	if args.NodoID == 0 {
		return op.Str("el nuevo nodoID aleatorio debe ser definido por quien invoca")
	}
	if !esTipoValido(args.Tipo) {
		return op.Strf("tipo de nodo '%v' inválido", args.Tipo)
	}

	padre, err := s.repo.GetNodo(args.PadreID)
	if err != nil {
		return op.Err(err).Msg("Padre indefinido")
	}

	args.Titulo = gkt.SinEspaciosExtra(args.Titulo)
	if args.Titulo == "" {
		return op.Msg("No indicó ningún texto para crear la entidad")
	}

	if padre.EsTarea() || padre.EsTramo() || padre.EsRegla() {
		return op.Msg("el nodo padre no puede tener descendientes")
	}

	err = s.repo.InsertNodo(Nodo{
		NodoID:  args.NodoID,
		PadreID: args.PadreID,
		Tipo:    args.Tipo,
		Titulo:  args.Titulo,
	})
	if err != nil {
		return op.Err(err)
	}

	err = s.riseEvent(EvNodoAgregado, evNodoAgregado{
		NodoID:  args.NodoID,
		PadreID: args.PadreID,
		Tipo:    args.Tipo,
		Titulo:  args.Titulo,
	})
	if err != nil {
		return op.Err(err)
	}

	return nil
}

const EvTareaAgregada gko.EventKey = "tarea.agregada"

type evTareaAgregada struct {
	NodoID  int
	PadreID int
	Tipo    string
	Titulo  string
}

func (e evTareaAgregada) ToMsg(t string) string {
	switch t {
	case "key":
		return string(EvTareaAgregada)
	default:
		return fmt.Sprintf("Tarea agregada ID: %v Título: '%s'", e.NodoID, e.Titulo)
	}
}

type ArgsAgregarTarea struct {
	Tipo     string
	NodoID   int
	PadreID  int
	Titulo   string
	Estimado string
}

func (s *AppTx) AgregarTarea(args ArgsAgregarTarea) error {
	op := gko.Op("AgregarTarea")

	if args.NodoID == 0 {
		return op.Str("el nuevo nodoID aleatorio debe ser definido por quien invoca")
	}
	if !esTipoValido(args.Tipo) {
		return op.Strf("tipo de nodo '%v' inválido", args.Tipo)
	}

	padre, err := s.repo.GetNodo(args.PadreID)
	if err != nil {
		return op.Err(err)
	}

	args.Titulo = gkt.SinEspaciosExtra(args.Titulo)
	if args.Titulo == "" {
		return op.Msg("No indicó ningún texto para crear la entidad")
	}

	if padre.EsTarea() || padre.EsTramo() || padre.EsRegla() {
		return op.Msg("el nodo padre no puede tener descendientes")
	}

	estimado, err := ust.NuevaDuraciónSegundos(args.Estimado)
	if err != nil {
		return op.Err(err)
	}

	err = s.repo.InsertNodo(Nodo{
		NodoID:   args.NodoID,
		PadreID:  args.PadreID,
		Tipo:     args.Tipo,
		Titulo:   args.Titulo,
		Segundos: estimado,
	})
	if err != nil {
		return op.Err(err)
	}

	err = s.riseEvent(EvTareaAgregada, evTareaAgregada{
		NodoID:  args.NodoID,
		PadreID: args.PadreID,
		Tipo:    args.Tipo,
		Titulo:  args.Titulo,
	})
	if err != nil {
		return op.Err(err)
	}

	return nil
}
