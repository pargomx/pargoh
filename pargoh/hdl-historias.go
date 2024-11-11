package main

import (
	"monorepo/dhistorias"
	"monorepo/sqliteust"
	"monorepo/ust"

	"github.com/pargomx/gecko"
	"github.com/pargomx/gecko/gko"
)

// ================================================================ //
// ========== READ ================================================ //

func (s *servidor) getHistoria(c *gecko.Context) error {
	Historia, err := dhistorias.GetHistoria(c.PathInt("historia_id"), dhistorias.GetDescendientes, s.repo)
	if err != nil {
		return err
	}
	data := map[string]any{
		"Titulo":   Historia.Historia.Titulo,
		"Agregado": Historia,

		"ListaTipoTarea": ust.ListaTipoTarea,
	}
	return c.RenderOk("historia", data)
}

func (s *servidor) getHistoriaTablero(c *gecko.Context) error {
	Historia, err := dhistorias.GetHistoria(c.PathInt("historia_id"), dhistorias.GetDescendientes, s.repo)
	if err != nil {
		return err
	}
	data := map[string]any{
		"Titulo":   Historia.Historia.Titulo,
		"Agregado": Historia,
	}
	return c.RenderOk("hist_tablero", data)
}

// ================================================================ //
// ========== WRITE =============================================== //

func (s *servidor) postHistoriaDePersona(c *gecko.Context) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	repotx := sqliteust.NuevoRepo(tx)
	nuevaHistoria := ust.Historia{
		HistoriaID: ust.NewRandomID(),
		Titulo:     c.FormVal("titulo"),
		Objetivo:   c.FormVal("objetivo"),
		Prioridad:  c.FormInt("prioridad"),
		Completada: c.FormBool("completada"),
	}
	err = dhistorias.AgregarHistoria(c.PathInt("persona_id"), nuevaHistoria, repotx)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	defer s.reloader.brodcastReload(c)
	return c.Redirf("/personas/%v", c.PathInt("persona_id"))
}

func (s *servidor) postHistoriaDeHistoria(c *gecko.Context) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	repotx := sqliteust.NuevoRepo(tx)
	nuevaHistoria := ust.Historia{
		HistoriaID: ust.NewRandomID(),
		Titulo:     c.FormVal("titulo"),
		Objetivo:   c.FormVal("objetivo"),
		Prioridad:  c.FormInt("prioridad"),
		Completada: c.FormBool("completada"),
	}
	err = dhistorias.AgregarHistoria(c.PathInt("historia_id"), nuevaHistoria, repotx)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	defer s.reloader.brodcastReload(c)
	return c.Redirf("/historias/%v", c.PathInt("historia_id"))
}

func (s *servidor) updateHistoria(c *gecko.Context) error {
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
	return c.StatusOk("Historia actualizada")
}

func (s *servidor) patchHistoria(c *gecko.Context) error {
	err := dhistorias.ParcharHistoria(
		c.PathInt("historia_id"),
		c.PathVal("param"),
		c.FormValue("value"),
		s.repo,
	)
	if err != nil {
		return err
	}
	defer s.reloader.brodcastReload(c)
	return c.StatusOk("Historia parchada")
}

func (s *servidor) priorizarHistoria(c *gecko.Context) error {
	err := dhistorias.PriorizarHistoria(c.PathInt("historia_id"), c.FormInt("prioridad"), s.repo)
	if err != nil {
		return err
	}
	return c.StatusOk("Historia priorizada")
}

func (s *servidor) priorizarHistoriaNuevo(c *gecko.Context) error {
	err := dhistorias.PriorizarHistoria(c.PathInt("historia_id"), c.PathInt("prioridad"), s.repo)
	if err != nil {
		return err
	}
	return c.RefreshHTMX()
}

func (s *servidor) marcarHistoria(c *gecko.Context) error {
	err := dhistorias.MarcarHistoria(c.PathInt("historia_id"), c.FormBool("completada"), s.repo)
	if err != nil {
		return err
	}
	return c.StatusOk("Historia marcada")
}

func (s *servidor) marcarHistoriaNueva(c *gecko.Context) error {
	err := dhistorias.MarcarHistoria(c.PathInt("historia_id"), c.PathBool("completada"), s.repo)
	if err != nil {
		return err
	}
	return c.RefreshHTMX()
}

func (s *servidor) deleteHistoria(c *gecko.Context) error {
	padreID, err := dhistorias.EliminarHistoria(c.PathInt("historia_id"), s.repo)
	if err != nil {
		return err
	}
	padre, err := s.repo.GetNodo(padreID)
	if err != nil {
		return err
	}
	if padre.EsHistoria() {
		return c.Redirf("/historias/%v", padreID)
	} else if padre.EsPersona() {
		return c.Redirf("/personas/%v", padreID)
	} else {
		gko.LogWarnf("deleteHistoria: padre %v no es persona ni historia", padreID)
		return c.Redir("/proyectos")
	}
}

func (s *servidor) reordenarHistoria(c *gecko.Context) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	err = dhistorias.ReordenarNodo(c.FormInt("historia_id"), c.FormInt("new_pos"), sqliteust.NuevoRepo(tx))
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	hist, err := s.repo.GetNodoHistoria(c.FormInt("historia_id"))
	if err != nil {
		return err
	}
	defer s.reloader.brodcastReload(c)

	if hist.PadreTbl == ust.TipoNodoPersona {
		return c.Redirf("/personas/%v", hist.PadreID)
	} else if hist.PadreTbl == ust.TipoNodoHistoria {
		return c.Redirf("/historias/%v", hist.PadreID)
	} else {
		return gko.ErrInesperado().Msgf("reordenarHistoria: padre %v no es persona ni historia, sino %v", hist.PadreID, hist.PadreTbl)
	}
}

func (s *servidor) moverHistoria(c *gecko.Context) error {
	nuevoPadreID := c.FormInt("target_historia_id")
	if nuevoPadreID == 0 {
		nuevoPadreID = c.FormInt("target_persona_id")
		if nuevoPadreID == 0 {
			nuevoPadreID = c.FormInt("nuevo_padre_id")
		}
	}
	historiaID := c.FormInt("historia_id")
	if historiaID == 0 {
		historiaID = c.PathInt("historia_id")
	}
	err := dhistorias.MoverHistoria(historiaID, nuevoPadreID, s.repo)
	if err != nil {
		return err
	}
	defer s.reloader.brodcastReload(c)
	// TODO: enviar link a la nueva ubicación como sugerencia.
	return c.RefreshHTMX()
}
