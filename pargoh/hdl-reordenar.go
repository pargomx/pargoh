package main

import (
	"monorepo/arbol"

	"github.com/pargomx/gecko"
)

// ================================================================ //
// ========== Reordenar =========================================== //

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

func (s *writehdl) moverNodo(c *gecko.Context, tx *handlerTx) error {
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
