package main

import (
	"monorepo/dhistorias"

	"github.com/pargomx/gecko"
)

// ================================================================ //
// ========== VIAJE DE USUARIO ==================================== //

func (s *servidor) postTramoDeViaje(c *gecko.Context) error {
	err := dhistorias.AgregarTramoDeViaje(s.repoOld, c.PathInt("historia_id"), c.FormValue("texto"))
	if err != nil {
		return err
	}
	defer s.reloader.brodcastReload(c)
	return c.RedirOtrof("/historias/%v", c.PathInt("historia_id"))
}

func (s *servidor) deleteTramoDeViaje(c *gecko.Context) error {
	tx, err := s.newRepoTx()
	if err != nil {
		return err
	}
	err = dhistorias.EliminarTramoDeViaje(tx.repoOld, c.PathInt("historia_id"), c.PathInt("posicion"))
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	defer s.reloader.brodcastReload(c)
	return c.RedirOtrof("/historias/%v", c.PathInt("historia_id"))
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

func (s *servidor) moverTramo(c *gecko.Context) error {
	historiaID, err := dhistorias.MoverTramo(c.FormInt("historia_id"), c.FormInt("posicion"), c.FormInt("target_historia_id"), s.repoOld)
	if err != nil {
		return err
	}
	defer s.reloader.brodcastReload(c)
	return c.RedirOtrof("/historias/%v", historiaID)
}
