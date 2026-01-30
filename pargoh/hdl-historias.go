package main

import (
	"monorepo/arbol"
	"monorepo/dhistorias"
	"monorepo/ust"

	"github.com/pargomx/gecko"
)

// ================================================================ //
// ========== READ ================================================ //

func (s *readhdl) getHistoria(c *gecko.Context) error {
	Historia, err := s.repo.GetHistoria(c.PathInt("historia_id"))
	if err != nil {
		return err
	}

	data := map[string]any{
		"Titulo":          Historia.Titulo,
		"Agregado":        Historia,
		"ScriptsHistoria": true,
		"OldGrafico":      c.QueryBool("old"),
		"ListaTipoTarea":  ust.ListaTipoTarea,
	}
	return c.RenderOk("historia", data)
}

func (s *readhdl) getHistoriaTablero(c *gecko.Context) error {
	Historia, err := dhistorias.GetHistoria(c.PathInt("historia_id"), dhistorias.GetDescendientes, s.repoOld)
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
	args := arbol.ArgsAgregarHoja{
		Tipo:    c.FormVal("tipo"),
		NodoID:  ust.NewRandomID(),
		PadreID: c.FormInt("padre_id"),
		Titulo:  c.FormVal("titulo"),
	}
	err := tx.app.AgregarHoja(args)
	if err != nil {
		return err
	}
	return c.RedirOtrof("/n/%v", c.FormInt("padre_id"))
}

func (s *writehdl) postProyecto(c *gecko.Context, tx *handlerTx) error {
	args := arbol.ArgsAgregarHoja{
		Tipo:    "PRY",
		NodoID:  ust.NewRandomID(),
		PadreID: arbol.NodoProyectosActivos,
		Titulo:  c.FormVal("titulo"),
	}
	err := tx.app.AgregarHoja(args)
	if err != nil {
		return err
	}
	return c.RedirOtro("/")
}

func (s *writehdl) postHistoriaDePersona(c *gecko.Context, tx *handlerTx) error {
	args := arbol.ArgsAgregarHoja{
		Tipo:    "HIS",
		NodoID:  ust.NewRandomID(),
		PadreID: c.PathInt("persona_id"),
		Titulo:  c.FormVal("titulo"),
	}
	err := tx.app.AgregarHoja(args)
	if err != nil {
		return err
	}
	defer s.reloader.brodcastReload(c)
	return c.AskedForFallback("/personas/%v", args.PadreID)
}

func (s *writehdl) postHistoriaDeHistoria(c *gecko.Context, tx *handlerTx) error {
	args := arbol.ArgsAgregarHoja{
		Tipo:    "HIS",
		NodoID:  ust.NewRandomID(),
		PadreID: c.PathInt("historia_id"),
		Titulo:  c.FormVal("titulo"),
	}
	err := tx.app.AgregarHoja(args)
	if err != nil {
		return err
	}
	defer s.reloader.brodcastReload(c)
	return c.AskedForFallback("/historias/%v", args.PadreID)
}

// Agregar historia de usuario como padre de la actual.
func (s *writehdl) postPadreParaHistoria(c *gecko.Context, tx *handlerTx) error {
	nod, err := tx.repo.GetNodo(c.PathInt("historia_id"))
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
	return c.AskedForFallback("/historias/%v", newPadre.NodoID)
}

func (s *writehdl) patchHistoria(c *gecko.Context, tx *handlerTx) error {
	err := tx.app.ParcharNodo(arbol.ArgsParcharNodo{
		NodoID: c.PathInt("historia_id"),
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
		NodoID: c.PathInt("historia_id"),
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
		NodoID: c.PathInt("historia_id"),
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
