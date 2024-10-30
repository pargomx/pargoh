package main

import (
	"monorepo/dhistorias"
	"monorepo/sqliteust"

	"github.com/pargomx/gecko"
	"github.com/pargomx/gecko/gko"
)

func (s *servidor) postReferencia(c *gecko.Context) error {
	refHistoriaID := c.FormInt("target_historia_id")
	if refHistoriaID == 0 {
		return gko.ErrDatoIndef().Msg("Debe seleccionar una historia de usuario")
	}
	err := dhistorias.AgregarReferencia(s.repo, c.PathInt("historia_id"), refHistoriaID)
	if err != nil {
		return err
	}
	defer s.reloader.brodcastReload(c)
	return c.Redir("/historias/%v", c.PathInt("historia_id"))
}

func (s *servidor) deleteReferencia(c *gecko.Context) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	err = dhistorias.EliminarReferencia(sqliteust.NuevoRepo(tx), c.PathInt("historia_id"), c.PathInt("ref_historia_id"))
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	defer s.reloader.brodcastReload(c)
	return c.Redir("/historias/%v", c.PathInt("historia_id"))
}
