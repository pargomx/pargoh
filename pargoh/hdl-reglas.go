package main

import (
	"monorepo/arbol"
	"monorepo/ust"

	"github.com/pargomx/gecko"
)

// ================================================================ //
// ========== REGLAS DE NEGOCIO =================================== //

func (s *writehdl) postRegla(c *gecko.Context, tx *handlerTx) error {
	args := arbol.ArgsAgregarHoja{
		Tipo:    "REG",
		NodoID:  ust.NewRandomID(),
		PadreID: c.PathInt("historia_id"),
		Titulo:  c.FormValue("texto"),
	}
	err := tx.app.AgregarHoja(args)
	if err != nil {
		return err
	}
	defer s.reloader.brodcastReload(c)
	// TODO: Solo enviar el fragmento.
	return c.RedirOtrof("/historias/%v", args.PadreID)
}

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
	return c.RedirOtrof("/historias/%v", c.PathInt("historia_id"))
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
	return c.RedirOtrof("/historias/%v", c.PathInt("historia_id"))
}
