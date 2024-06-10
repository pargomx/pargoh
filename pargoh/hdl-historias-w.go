package main

import (
	"monorepo/gecko"
	"monorepo/historias_de_usuario/dhistorias"
	"monorepo/historias_de_usuario/sqliteust"
	"monorepo/historias_de_usuario/ust"
)

func (s *servidor) postHistoria(c *gecko.Context) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	repotx := sqliteust.NuevoRepositorio(tx)
	nuevaHistoria := ust.Historia{
		HistoriaID: ust.NewRandomID(),
		Titulo:     c.FormVal("titulo"),
		Objetivo:   c.FormVal("objetivo"),
		Prioridad:  c.FormInt("prioridad"),
		Completada: c.FormBool("completada"),
	}
	err = dhistorias.AgregarHistoria(c.PathInt("nodo_id"), nuevaHistoria, repotx)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return c.StatusOk("Historia creada")
}

func (s *servidor) patchHistoria(c *gecko.Context) error {
	err := dhistorias.ActualizarHistoria(
		c.PathInt("historia_id"),
		ust.Historia{
			HistoriaID: c.FormInt("historia_id"),
			Titulo:     c.FormValue("titulo"),
			Objetivo:   c.FormValue("objetivo"),
			Prioridad:  c.FormInt("prioridad"),
			Completada: c.FormBool("completada"),
		},
		s.repo,
	)
	if err != nil {
		return err
	}
	return c.StatusAccepted("Historia actualizada")
}

func (s *servidor) priorizarHistoria(c *gecko.Context) error {
	err := dhistorias.PriorizarHistoria(c.PathInt("historia_id"), c.FormInt("prioridad"), s.repo)
	if err != nil {
		return err
	}
	return c.StatusAccepted("Historia priorizada")
}

func (s *servidor) marcarHistoria(c *gecko.Context) error {
	err := dhistorias.MarcarHistoria(c.PathInt("historia_id"), c.FormBool("completada"), s.repo)
	if err != nil {
		return err
	}
	return c.StatusAccepted("Historia marcada")
}

func (s *servidor) deleteHistoria(c *gecko.Context) error {
	err := dhistorias.EliminarHistoria(c.PathInt("historia_id"), s.repo)
	if err != nil {
		return err
	}
	return c.StatusAccepted("Historia eliminada")
}

func (s *servidor) moverHistoria(c *gecko.Context) error {
	err := dhistorias.MoverHistoria(c.PathInt("historia_id"), c.FormInt("nuevo_padre_id"), s.repo)
	if err != nil {
		return err
	}
	return c.Redir("/historias/%v/mover", c.PathInt("historia_id"))
	// return c.StatusAccepted("Historia movida")
}
