package main

import (
	"monorepo/arbol"
	"monorepo/dhistorias"
	"monorepo/ust"

	"github.com/pargomx/gecko"
	"github.com/pargomx/gecko/gko"
)

// ================================================================ //
// ========== READ ================================================ //

func (s *readhdl) getNodoTablero(c *gecko.Context) error {
	Historia, err := dhistorias.GetHistoria(c.PathInt("nodo_id"), dhistorias.GetDescendientes, s.repoOld)
	if err != nil {
		return err
	}
	data := map[string]any{
		"Titulo":   Historia.Historia.Titulo,
		"Agregado": Historia,
	}
	return c.RenderOk("hist_tablero", data)
}

// ================================================================ //
// ========== WRITE =============================================== //

func (s *writehdl) postNodo(c *gecko.Context, tx *handlerTx) error {
	padreID := c.PathInt("nodo_id")
	args := arbol.ArgsAgregarHoja{
		Tipo:    c.FormVal("tipo"),
		NodoID:  ust.NewRandomID(),
		PadreID: padreID,
		Titulo:  c.FormVal("titulo"),
	}
	err := tx.app.AgregarHoja(args)
	if err != nil {
		return err
	}
	// defer s.reloader.brodcastReload(c)
	return c.RedirOtrof("/h/%v", padreID)
}

func (s *writehdl) postNodoDeTipo(c *gecko.Context, tx *handlerTx) error {
	padreID := c.PathInt("nodo_id")

	tipo := ""
	switch c.PathVal("tipo") {
	case "grupo":
		tipo = arbol.TipoGrupo
	case "proyecto":
		tipo = arbol.TipoProyecto
	case "persona":
		tipo = arbol.TipoPersona
	case "historia":
		tipo = arbol.TipoHistoria
	case "tarea":
		tipo = arbol.TipoTarea
	case "regla":
		tipo = arbol.TipoRegla
	case "viaje":
		tipo = arbol.TipoViaje
	default:
		return gko.ErrDatoInvalido.Msgf("No se puede agregar nodo de tipo '%v'", c.PathVal("tipo"))
	}

	args := arbol.ArgsAgregarHoja{
		Tipo:    tipo,
		NodoID:  ust.NewRandomID(),
		PadreID: padreID,
		Titulo:  c.FormVal("titulo"),
	}
	err := tx.app.AgregarHoja(args)
	if err != nil {
		return err
	}
	// defer s.reloader.brodcastReload(c)
	return c.RedirOtrof("/h/%v", padreID)
}

// Agregar historia de usuario como padre de la actual.
func (s *writehdl) postNodoPadre(c *gecko.Context, tx *handlerTx) error {
	nod, err := tx.repo.GetNodo(c.PathInt("nodo_id"))
	if err != nil {
		return err
	}
	newPadre := arbol.ArgsAgregarHoja{
		Tipo:    "HIS",
		NodoID:  ust.NewRandomID(),
		PadreID: nod.PadreID,
		Titulo:  c.PromptVal(),
	}
	err = tx.app.AgregarHoja(newPadre)
	if err != nil {
		return err
	}

	err = tx.app.MoverHoja(arbol.ArgsMover{
		NodoID:     nod.NodoID,
		NewPadreID: newPadre.NodoID,
	})
	if err != nil {
		return err
	}

	// Mover el padre a la misma posición en que estaba el otro nodo.
	if nod.Posicion > 1 {
		tx.app.ReordenarEntidad(arbol.ArgsReordenar{
			NodoID: newPadre.NodoID,
			NewPos: nod.Posicion,
		})
	}

	defer s.reloader.brodcastReload(c)
	return c.AskedForFallback("/h/%v", newPadre.NodoID)
}

func (s *writehdl) patchHistoria(c *gecko.Context, tx *handlerTx) error {
	err := tx.app.ParcharNodo(arbol.ArgsParcharNodo{
		NodoID: c.PathInt("nodo_id"),
		Campo:  c.PathVal("param"),
		NewVal: c.FormValue("value"),
	})
	if err != nil {
		return err
	}
	defer s.reloader.brodcastReload(c)
	return c.AskedFor("Historia parchada")
}

func (s *writehdl) priorizarHistoria(c *gecko.Context, tx *handlerTx) error {
	args := arbol.ArgsParcharNodo{
		NodoID: c.PathInt("nodo_id"),
		Campo:  "prioridad",
		NewVal: c.FormVal("prioridad"),
	}
	if args.NewVal == "" {
		args.NewVal = c.PathVal("prioridad")
	}
	err := tx.app.ParcharNodo(args)
	if err != nil {
		return err
	}
	return c.AskedFor("Historia priorizada")
}

func (s *writehdl) marcarHistoria(c *gecko.Context, tx *handlerTx) error {
	args := arbol.ArgsParcharNodo{
		NodoID: c.PathInt("nodo_id"),
		Campo:  "completada",
		NewVal: c.FormVal("completada"),
	}
	if args.NewVal == "" {
		args.NewVal = c.PathVal("completada")
	}
	err := tx.app.ParcharNodo(args)
	if err != nil {
		return err
	}
	return c.AskedFor("Historia marcada")
}
