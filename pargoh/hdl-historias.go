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

func (s *servidor) updateHistoria(c *gecko.Context) error {
	err := dhistorias.ActualizarHistoria(
		c.PathInt("historia_id"),
		ust.Historia{
			HistoriaID: c.FormInt("historia_id"),
			Titulo:     c.FormValue("titulo"),
			Objetivo:   c.FormValue("objetivo"),
			Prioridad:  c.FormInt("prioridad"),
			Completada: c.FormBool("completada"),
		},
		s.repoOld,
	)
	if err != nil {
		return err
	}
	return c.AskedFor("Historia actualizada")
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

func (s *servidor) priorizarHistoria(c *gecko.Context) error {
	err := dhistorias.PriorizarHistoria(c.PathInt("historia_id"), c.FormInt("prioridad"), s.repoOld)
	if err != nil {
		return err
	}
	return c.AskedFor("Historia priorizada")
}

func (s *servidor) priorizarHistoriaNuevo(c *gecko.Context) error {
	err := dhistorias.PriorizarHistoria(c.PathInt("historia_id"), c.PathInt("prioridad"), s.repoOld)
	if err != nil {
		return err
	}
	return c.AskedFor("Historia priorizada")
}

func (s *servidor) marcarHistoria(c *gecko.Context) error {
	err := dhistorias.MarcarHistoria(c.PathInt("historia_id"), c.FormBool("completada"), s.repoOld)
	if err != nil {
		return err
	}
	return c.AskedFor("Historia marcada")
}

func (s *servidor) marcarHistoriaNueva(c *gecko.Context) error {
	err := dhistorias.MarcarHistoria(c.PathInt("historia_id"), c.PathBool("completada"), s.repoOld)
	if err != nil {
		return err
	}
	return c.AskedFor("Historia marcada")
}
