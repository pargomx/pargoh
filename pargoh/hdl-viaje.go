package main

import (
	"monorepo/dhistorias"
	"monorepo/sqliteust"

	"github.com/pargomx/gecko"
)

// ================================================================ //
// ========== VIAJE DE USUARIO ==================================== //

func (s *servidor) postTramoDeViaje(c *gecko.Context) error {
	err := dhistorias.AgregarTramoDeViaje(s.repo, c.PathInt("historia_id"), c.FormValue("texto"))
	if err != nil {
		return err
	}
	defer s.reloader.brodcastReload(c)
	return c.RedirOtrof("/historias/%v", c.PathInt("historia_id"))
}

func (s *servidor) deleteTramoDeViaje(c *gecko.Context) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	err = dhistorias.EliminarTramoDeViaje(sqliteust.NuevoRepo(tx), c.PathInt("historia_id"), c.PathInt("posicion"))
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	defer s.reloader.brodcastReload(c)
	return c.RedirOtrof("/historias/%v", c.PathInt("historia_id"))
}

func (s *servidor) patchTramoDeViaje(c *gecko.Context) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	err = dhistorias.EditarTramoDeViaje(sqliteust.NuevoRepo(tx), c.PathInt("historia_id"), c.PathInt("posicion"), c.FormValue("texto"))
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	defer s.reloader.brodcastReload(c)
	return c.RedirOtrof("/historias/%v", c.PathInt("historia_id"))
}

func (s *servidor) reordenarTramo(c *gecko.Context) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	err = dhistorias.ReordenarTramo(sqliteust.NuevoRepo(tx), c.FormInt("historia_id"), c.FormInt("old_pos"), c.FormInt("new_pos"))
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	defer s.reloader.brodcastReload(c)
	return c.RedirOtrof("/historias/%v", c.FormInt("historia_id"))
}

func (s *servidor) moverTramo(c *gecko.Context) error {
	historiaID, err := dhistorias.MoverTramo(c.FormInt("historia_id"), c.FormInt("posicion"), c.FormInt("target_historia_id"), s.repo)
	if err != nil {
		return err
	}
	defer s.reloader.brodcastReload(c)
	return c.RedirOtrof("/historias/%v", historiaID)
}
