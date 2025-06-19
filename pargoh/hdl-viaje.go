package main

import (
	"monorepo/arbol"
	"monorepo/dhistorias"
	"monorepo/ust"

	"github.com/pargomx/gecko"
)

// ================================================================ //
// ========== VIAJE DE USUARIO ==================================== //

func (s *writehdl) postTramoDeViaje(c *gecko.Context, tx *handlerTx) error {
	args := arbol.ArgsAgregarHoja{
		Tipo:    "VIA",
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

func (s *servidor) patchTramoDeViaje(c *gecko.Context) error {
	tx, err := s.newRepoTx()
	if err != nil {
		return err
	}
	err = dhistorias.EditarTramoDeViaje(tx.repoOld, c.PathInt("historia_id"), c.PathInt("posicion"), c.FormValue("texto"))
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	defer s.reloader.brodcastReload(c)
	return c.RedirOtrof("/historias/%v", c.PathInt("historia_id"))
}
