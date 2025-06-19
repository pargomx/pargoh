package main

import (
	"monorepo/arbol"

	"github.com/pargomx/gecko"
	"github.com/pargomx/gecko/gko"
)

func (s *writehdl) reordenarPersona(c *gecko.Context, tx *handlerTx) error {
	err := tx.app.ReordenarEntidad(arbol.ArgsReordenar{
		NodoID: c.FormInt("persona_id"),
		NewPos: c.FormInt("new_pos"),
	})
	if err != nil {
		return err
	}
	pers, err := tx.repo.GetPersona(c.FormInt("persona_id"))
	if err != nil {
		return err
	}
	return c.RedirOtrof("/proyectos/%v", pers.ProyectoID)
}

func (s *writehdl) reordenarHistoria(c *gecko.Context, tx *handlerTx) error {
	err := tx.app.ReordenarEntidad(arbol.ArgsReordenar{
		NodoID: c.FormInt("historia_id"),
		NewPos: c.FormInt("new_pos"),
	})
	if err != nil {
		return err
	}

	his, err := tx.repo.GetHistoria(c.FormInt("historia_id"))
	if err != nil {
		return err
	}
	defer s.reloader.brodcastReload(c)

	padre, err := tx.repo.GetNodo(his.PadreID)
	if err != nil {
		return err
	}
	if padre.EsPersona() {
		return c.RedirOtrof("/personas/%v", his.PadreID)
	} else if padre.EsHistoriaDeUsuario() {
		return c.RedirOtrof("/historias/%v", his.PadreID)
	}
	return gko.ErrInesperado.Msgf("reordenarHistoria: padre %v no es persona ni historia, sino %v",
		padre.NodoID, padre.Tipo)
}

func (s *writehdl) reordenarTramo(c *gecko.Context, tx *handlerTx) error {
	err := tx.app.ReordenarEntidad(arbol.ArgsReordenar{
		NodoID: c.FormInt("historia_id"),
		NewPos: c.FormInt("new_pos"),
	})
	if err != nil {
		return err
	}
	defer s.reloader.brodcastReload(c)
	return c.RedirOtrof("/historias/%v", c.FormInt("historia_id"))
}

func (s *writehdl) reordenarRegla(c *gecko.Context, tx *handlerTx) error {
	err := tx.app.ReordenarEntidad(arbol.ArgsReordenar{
		NodoID: c.FormInt("historia_id"),
		NewPos: c.FormInt("new_pos"),
	})
	if err != nil {
		return err
	}
	defer s.reloader.brodcastReload(c)
	// TODO: Solo enviar el fragmento.
	return c.RedirOtrof("/historias/%v", c.FormInt("historia_id"))
}

// ================================================================ //
// ========== Mover =============================================== //

func (s *writehdl) moverHistoria(c *gecko.Context, tx *handlerTx) error {
	newPadreID := c.FormInt("target_historia_id")
	if newPadreID == 0 {
		newPadreID = c.FormInt("target_persona_id")
		if newPadreID == 0 {
			newPadreID = c.FormInt("nuevo_padre_id")
		}
	}
	historiaID := c.FormInt("historia_id")
	if historiaID == 0 {
		historiaID = c.PathInt("historia_id")
	}

	err := tx.app.MoverHoja(arbol.ArgsMover{
		NodoID:     historiaID,
		NewPadreID: newPadreID,
	})
	if err != nil {
		return err
	}
	defer s.reloader.brodcastReload(c)
	// TODO: enviar link a la nueva ubicación como sugerencia.
	return c.RefreshHTMX()
}

func (s *writehdl) moverTramo(c *gecko.Context, tx *handlerTx) error {
	// historiaID, err := dhistorias.MoverTramo(c.FormInt("historia_id"), c.FormInt("posicion"), c.FormInt("target_historia_id"), s.repoOld)
	args := arbol.ArgsMover{
		NodoID:     c.FormInt("nodo_id"),
		NewPadreID: c.FormInt("new_padre_id"),
	}
	err := tx.app.MoverHoja(args)
	if err != nil {
		return err
	}
	defer s.reloader.brodcastReload(c)
	return c.RedirOtrof("/historias/%v", args.NewPadreID)
}

func (s *writehdl) moverTarea(c *gecko.Context, tx *handlerTx) error {
	args := arbol.ArgsMover{
		NodoID:     c.FormInt("nodo_id"),
		NewPadreID: c.FormInt("new_padre_id"),
	}
	err := tx.app.MoverHoja(args)
	if err != nil {
		return err
	}
	defer s.reloader.brodcastReload(c)
	return c.AskedForFallback("/historias/%v", args.NewPadreID)
}
