package main

import (
	"monorepo/dhistorias"
	"monorepo/ust"

	"github.com/pargomx/gecko"
	"github.com/pargomx/gecko/gko"
)

func (s *servidor) postTarea(c *gecko.Context) error {
	estimado, err := ust.NuevaDuración(c.FormVal("tiempo_estimado"))
	if err != nil {
		return err
	}
	tarea := ust.Tarea{
		TareaID:        ust.NewRandomID(),
		HistoriaID:     c.PathInt("historia_id"),
		Tipo:           ust.SetTipoTareaDB(c.FormVal("tipo")),
		Descripcion:    c.FormVal("descripcion"),
		Impedimentos:   c.FormVal("impedimentos"),
		TiempoEstimado: estimado,
	}
	err = dhistorias.AgregarTarea(tarea, s.repo)
	if err != nil {
		return err
	}
	defer s.reloader.brodcastReload(c)
	return c.Redir("/historias/%v", tarea.HistoriaID)
}

func (s *servidor) modificarTarea(c *gecko.Context) error {
	estimado, err := ust.NuevaDuración(c.FormVal("tiempo_estimado"))
	if err != nil {
		return err
	}
	tarea := ust.Tarea{
		TareaID:        c.FormInt("tarea_id"),
		HistoriaID:     c.FormInt("historia_id"),
		Tipo:           ust.SetTipoTareaDB(c.FormVal("tipo")),
		Descripcion:    c.FormVal("descripcion"),
		Impedimentos:   c.FormVal("impedimentos"),
		TiempoEstimado: estimado,
	}
	err = dhistorias.ActualizarTarea(c.PathInt("tarea_id"), tarea, s.repo)
	if err != nil {
		return err
	}
	defer s.reloader.brodcastReload(c)
	return c.Redir("/historias/%v", tarea.HistoriaID)
}

func (s *servidor) moverTarea(c *gecko.Context) error {
	historiaID, err := dhistorias.MoverTarea(c.FormInt("tarea_id"), c.FormInt("target_historia_id"), s.repo)
	if err != nil {
		return err
	}
	defer s.reloader.brodcastReload(c)
	return c.Redir("/historias/%v", historiaID)
}

func (s *servidor) eliminarTarea(c *gecko.Context) error {
	historiaID, err := dhistorias.EliminarTarea(c.PathInt("tarea_id"), s.repo)
	if err != nil {
		return err
	}
	gko.LogInfof("Tarea %d eliminada", c.PathInt("tarea_id"))
	defer s.reloader.brodcastReload(c)
	return c.Redir("/historias/%v", historiaID)
}

func (s *servidor) iniciarTarea(c *gecko.Context) error {
	historiaID, err := dhistorias.IniciarTarea(c.PathInt("tarea_id"), s.repo)
	if err != nil {
		return err
	}
	defer s.reloader.brodcastReload(c)
	return c.Redir("/historias/%v", historiaID)
}
func (s *servidor) pausarTarea(c *gecko.Context) error {
	historiaID, err := dhistorias.PausarTarea(c.PathInt("tarea_id"), s.repo)
	if err != nil {
		return err
	}
	defer s.reloader.brodcastReload(c)
	return c.Redir("/historias/%v", historiaID)
}
func (s *servidor) terminarTarea(c *gecko.Context) error {
	historiaID, err := dhistorias.FinalizarTarea(c.PathInt("tarea_id"), s.repo)
	if err != nil {
		return err
	}
	defer s.reloader.brodcastReload(c)
	return c.Redir("/historias/%v", historiaID)
}

func (s *servidor) getTarea(c *gecko.Context) error {
	tarea, err := s.repo.GetTarea(c.PathInt("tarea_id"))
	if err != nil {
		return err
	}
	intervalos, err := s.repo.ListIntervalosByTareaID(tarea.TareaID)
	if err != nil {
		return err
	}
	data := map[string]any{
		"Tarea":      tarea,
		"Intervalos": intervalos,
	}
	return c.RenderOk("tarea", data)
}

func (s *servidor) getIntervalos(c *gecko.Context) error {
	recientes, err := s.repo.ListIntervalosRecientes()
	if err != nil {
		return err
	}
	abiertos, err := s.repo.ListIntervalosRecientesAbiertos()
	if err != nil {
		return err
	}
	data := map[string]any{
		"Titulo":    "Intervalos de trabajo",
		"Recientes": recientes,
		"Abiertos":  abiertos,
	}
	return c.RenderOk("intervalos", data)
}
