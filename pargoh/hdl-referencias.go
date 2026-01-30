package main

import (
	"monorepo/arbol"

	"github.com/pargomx/gecko"
)

func (s *writehdl) postReferencia(c *gecko.Context, tx *handlerTx) error {
	args := arbol.ArgsReferencia{
		NodoID:    c.PathInt("nodo_id"),
		RefNodoID: c.FormInt("target_nodo_id"),
	}
	err := tx.app.AgregarReferencia(args)
	if err != nil {
		return err
	}
	defer s.reloader.brodcastReload(c)
	return c.RedirOtrof("/h/%v", args.NodoID)
}

func (s *writehdl) deleteReferencia(c *gecko.Context, tx *handlerTx) error {
	args := arbol.ArgsReferencia{
		NodoID:    c.PathInt("nodo_id"),
		RefNodoID: c.PathInt("ref_nodo_id"),
	}
	err := tx.app.EliminarReferencia(args)
	if err != nil {
		return err
	}
	defer s.reloader.brodcastReload(c)
	return c.RedirOtrof("/h/%v", args.NodoID)
}
