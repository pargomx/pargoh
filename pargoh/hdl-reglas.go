package main

import (
	"monorepo/arbol"

	"github.com/pargomx/gecko"
)

// ================================================================ //
// ========== REGLAS DE NEGOCIO =================================== //

func (s *writehdl) patchRegla(c *gecko.Context, tx *handlerTx) error {
	err := tx.app.ParcharNodo(arbol.ArgsParcharNodo{
		NodoID: c.PathInt("regla_id"),
		Campo:  "texto",
		NewVal: c.FormValue("texto"),
	})
	if err != nil {
		return err
	}
	defer s.reloader.brodcastReload(c)
	// TODO: Solo enviar el fragmento.
	return c.RedirOtrof("/h/%v", c.PathInt("nodo_id"))
}

func (s *writehdl) marcarRegla(c *gecko.Context, tx *handlerTx) error {
	err := tx.app.ParcharNodo(arbol.ArgsParcharNodo{
		NodoID: c.PathInt("regla_id"),
		Campo:  "marcar_regla",
		NewVal: "",
	})
	if err != nil {
		return err
	}
	defer s.reloader.brodcastReload(c)
	// TODO: Solo enviar el fragmento.
	return c.RedirOtrof("/h/%v", c.PathInt("nodo_id"))
}
