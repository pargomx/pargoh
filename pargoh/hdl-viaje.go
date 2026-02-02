package main

import (
	"monorepo/arbol"

	"github.com/pargomx/gecko"
)

// ================================================================ //
// ========== VIAJE DE USUARIO ==================================== //

func (s *writehdl) patchTramoDeViaje(c *gecko.Context, tx *handlerTx) error {
	err := tx.app.ParcharNodo(arbol.ArgsParcharNodo{
		NodoID: c.PathInt("regla_id"),
		Campo:  "texto",
		NewVal: c.FormValue("texto"),
	})
	if err != nil {
		return err
	}
	defer s.reloader.brodcastReload(c)
	return c.RedirOtrof("/h/%v", c.PathInt("historia_id"))
}
