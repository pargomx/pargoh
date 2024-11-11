package main

import (
	"monorepo/dhistorias"
	"monorepo/sqliteust"

	"github.com/pargomx/gecko"
)

// ================================================================ //
// ========== REGLAS DE NEGOCIO =================================== //

func (s *servidor) postRegla(c *gecko.Context) error {
	err := dhistorias.AgregarRegla(s.repo, c.PathInt("historia_id"), c.FormValue("texto"))
	if err != nil {
		return err
	}
	defer s.reloader.brodcastReload(c)
	return c.Redirf("/historias/%v", c.PathInt("historia_id"))
}

func (s *servidor) deleteRegla(c *gecko.Context) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	err = dhistorias.EliminarRegla(sqliteust.NuevoRepo(tx), c.PathInt("historia_id"), c.PathInt("posicion"))
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	defer s.reloader.brodcastReload(c)
	return c.Redirf("/historias/%v", c.PathInt("historia_id"))
}

func (s *servidor) patchRegla(c *gecko.Context) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	err = dhistorias.EditarRegla(sqliteust.NuevoRepo(tx), c.PathInt("historia_id"), c.PathInt("posicion"), c.FormValue("texto"))
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	defer s.reloader.brodcastReload(c)
	return c.Redirf("/historias/%v", c.PathInt("historia_id"))
}

func (s *servidor) reordenarRegla(c *gecko.Context) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	err = dhistorias.ReordenarRegla(sqliteust.NuevoRepo(tx), c.FormInt("historia_id"), c.FormInt("old_pos"), c.FormInt("new_pos"))
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	defer s.reloader.brodcastReload(c)
	return c.Redirf("/historias/%v", c.FormInt("historia_id"))
}
