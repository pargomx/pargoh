package arbol

import (
	"github.com/pargomx/gecko/gko"
	"github.com/pargomx/gecko/gkt"
)

const EvNodoAgregado gko.EventKey = "nodo_added"

type ArgsAgregarHoja struct {
	Tipo    string // Tipo de nodo para agregar.
	NodoID  int    // Nuevo ID aleatorio.
	PadreID int    // PAdre al que se agregará la hoja.
	Titulo  string // Nombre o título de la entidad.
}

func (s *AppTx) AgregarHoja(args ArgsAgregarHoja) error {
	op := gko.Op("AgregarHoja")

	// Validar nuevo nodo
	if args.NodoID == 0 {
		return op.Str("el nuevo nodoID aleatorio debe ser definido por quien invoca")
	}
	if !esTipoValido(args.Tipo) {
		return op.Strf("tipo de nodo '%v' inválido", args.Tipo)
	}

	// Validar padre
	padre, err := s.repo.GetNodo(args.PadreID)
	if err != nil {
		return op.Err(err)
	}

	// Título
	args.Titulo = gkt.SinEspaciosExtra(args.Titulo)
	if args.Titulo == "" {
		return op.Msg("No indicó ningún texto para crear la entidad")
	}

	// Validar jerarquía
	if padre.EsTarea() {
		return op.Msg("el nodo padre es una tarea y no puede tener descendientes")
	}
	if padre.EsTramo() {
		return op.Msg("el nodo padre es un tramo y no puede tener descendientes")
	}
	if padre.EsRegla() {
		return op.Msg("el nodo padre es una regla y no puede tener descendientes")
	}

	// Insertar en la base de datos
	err = s.repo.InsertNodo(Nodo{
		NodoID:  args.NodoID,
		PadreID: args.PadreID,
		Tipo:    args.Tipo,
		Titulo:  args.Titulo,
	})
	if err != nil {
		return op.Err(err)
	}

	s.Results.Add(EvNodoAgregado.WithArgs(args))
	return nil
}

// ================================================================ //
// ========== TAREA =============================================== //

const EvTareaAgregado gko.EventKey = "tarea_added"

type ArgsAgregarTarea struct {
	Tipo     string // Tipo de nodo para agregar.
	NodoID   int    // Nuevo ID aleatorio.
	PadreID  int    // PAdre al que se agregará la hoja.
	Titulo   string // Nombre o título de la entidad.
	Estimado string // Estimado en segundos (opcional).
}

func (s *AppTx) AgregarTarea(args ArgsAgregarTarea) error {
	op := gko.Op("AgregarTarea")

	// Validar nuevo nodo
	if args.NodoID == 0 {
		return op.Str("el nuevo nodoID aleatorio debe ser definido por quien invoca")
	}
	if !esTipoValido(args.Tipo) {
		return op.Strf("tipo de nodo '%v' inválido", args.Tipo)
	}

	// Validar padre
	padre, err := s.repo.GetNodo(args.PadreID)
	if err != nil {
		return op.Err(err)
	}

	// Título
	args.Titulo = gkt.SinEspaciosExtra(args.Titulo)
	if args.Titulo == "" {
		return op.Msg("No indicó ningún texto para crear la entidad")
	}

	// Validar jerarquía
	if padre.EsTarea() {
		return op.Msg("el nodo padre es una tarea y no puede tener descendientes")
	}
	if padre.EsTramo() {
		return op.Msg("el nodo padre es un tramo y no puede tener descendientes")
	}
	if padre.EsRegla() {
		return op.Msg("el nodo padre es una regla y no puede tener descendientes")
	}

	// Insertar en la base de datos
	err = s.repo.InsertNodo(Nodo{
		NodoID:  args.NodoID,
		PadreID: args.PadreID,
		Tipo:    args.Tipo,
		Titulo:  args.Titulo,
	})
	if err != nil {
		return op.Err(err)
	}

	s.Results.Add(EvTareaAgregado.WithArgs(args))
	return nil
}
