package main

import (
	"monorepo/arbol"
	"monorepo/ust"

	"github.com/pargomx/gecko"
)

// ================================================================ //
// ========== VIAJE DE USUARIO ==================================== //

func (s *writehdl) postTramoDeViaje(c *gecko.Context, tx *handlerTx) error {
	args := arbol.ArgsAgregarHoja{
		Tipo:    "VIA",
		NodoID:  ust.NewRandomID(),
		PadreID: c.PathInt("nodo_id"),
		Titulo:  c.FormValue("texto"),
	}
	err := tx.app.AgregarHoja(args)
	if err != nil {
		return err
	}
	defer s.reloader.brodcastReload(c)
	// TODO: Solo enviar el fragmento.
	return c.RedirOtrof("/h/%v", args.PadreID)
}

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
