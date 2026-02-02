package main

import (
	"monorepo/arbol"

	"github.com/pargomx/gecko"
)

func (s *writehdl) reordenarNodo(c *gecko.Context, tx *handlerTx) error {
	err := tx.app.ReordenarEntidad(arbol.ArgsReordenar{
		NodoID: c.FormInt("nodo_id"),
		NewPos: c.FormInt("new_pos"),
	})
	if err != nil {
		return err
	}
	return c.AskedForFallback("/h/%v", c.FormInt("nodo_id"))
}

// ================================================================ //
// ========== Mover =============================================== //

func (s *writehdl) moverHistoria(c *gecko.Context, tx *handlerTx) error {
	newPadreID := c.FormInt("target_nodo_id")
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
	// historiaID, err := dhistorias.MoverTramo(c.FormInt("historia_id"), c.FormInt("posicion"), c.FormInt("target_nodo_id"), s.repoOld)
	args := arbol.ArgsMover{
		NodoID:     c.FormInt("nodo_id"),
		NewPadreID: c.FormInt("new_padre_id"),
	}
	err := tx.app.MoverHoja(args)
	if err != nil {
		return err
	}
	defer s.reloader.brodcastReload(c)
	return c.RedirOtrof("/h/%v", args.NewPadreID)
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
	return c.AskedForFallback("/h/%v", args.NewPadreID)
}
