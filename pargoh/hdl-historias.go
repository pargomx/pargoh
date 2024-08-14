package main

import (
	"monorepo/dhistorias"
	"monorepo/ust"

	"github.com/pargomx/gecko"
)

func (s *servidor) getHistoriasLista(c *gecko.Context) error {
	agg, err := dhistorias.GetHistoriasDePadre(c.PathInt("nodo_id"), s.repo)
	if err != nil {
		return err
	}
	titulo := "Nodo"
	if agg.Abuelo != nil {
		titulo = "Historia" // agg.Abuelo.Titulo
	} else {
		titulo = agg.Persona.Nombre
	}
	data := map[string]any{
		"Titulo":   titulo,
		"Agregado": agg,
	}
	return c.RenderOk("hist_lista", data)
}

func (s *servidor) getHistoriasTablero(c *gecko.Context) error {
	agg, err := dhistorias.GetHistoriasDePadre(c.PathInt("nodo_id"), s.repo)
	if err != nil {
		return err
	}
	titulo := "Nodo"
	if agg.Abuelo != nil {
		titulo = "Historia" // agg.Abuelo.Titulo
	} else {
		titulo = agg.Persona.Nombre
	}
	data := map[string]any{
		"Titulo":   titulo,
		"Agregado": agg,
	}
	return c.RenderOk("hist_tablero", data)
}

func (s *servidor) getHistoriasPrioritarias(c *gecko.Context) error {
	Historias, err := s.repo.ListNodoHistoriasPrioritarias()
	if err != nil {
		return err
	}
	data := map[string]any{
		"Titulo":    "Historias prioritarias",
		"Historias": Historias,
	}
	return c.RenderOk("hist_prioritarias", data)
}

func (s *servidor) formHistoria(c *gecko.Context) error {
	historia, err := s.repo.GetHistoria(c.PathInt("historia_id"))
	if err != nil {
		return err
	}
	data := map[string]any{
		"Titulo":   historia.Titulo,
		"Historia": historia,
	}
	return c.RenderOk("hist_form", data)
}

func (s *servidor) moverHistoriaForm(c *gecko.Context) error {
	historia, err := s.repo.GetNodoHistoria(c.PathInt("historia_id"))
	if err != nil {
		return err
	}
	arboles, err := dhistorias.GetArbolCompleto(s.repo)
	if err != nil {
		return err
	}
	data := map[string]any{
		"Titulo":   "Mover historia",
		"Arboles":  arboles,
		"Historia": historia,
	}
	return c.RenderOk("hist_mover", data)
}

func (s *servidor) getTareasDeHistoria(c *gecko.Context) error {
	historia, err := s.repo.GetNodoHistoria(c.PathInt("historia_id"))
	if err != nil {
		return err
	}
	tareas, err := s.repo.ListTareasByHistoriaID(historia.HistoriaID)
	if err != nil {
		return err
	}
	agg, err := dhistorias.GetHistoriasDePadre(historia.HistoriaID, s.repo)
	if err != nil {
		return err
	}
	data := map[string]any{
		"Titulo":   "Tareas",
		"Historia": historia,
		"Tareas":   tareas,
		"Agregado": agg,

		"ListaTipoTarea": ust.ListaTipoTarea,
	}
	return c.RenderOk("hist_tareas", data)
}

func (s *servidor) getHistoria(c *gecko.Context) error {
	agg, err := dhistorias.GetHistoria(c.PathInt("historia_id"), s.repo)
	if err != nil {
		return err
	}
	data := map[string]any{
		"Titulo":   agg.Historia.Titulo,
		"Agregado": agg,

		"ListaTipoTarea": ust.ListaTipoTarea,
	}
	return c.RenderOk("historia", data)
}

func (s *servidor) postTramoDeViaje(c *gecko.Context) error {
	err := dhistorias.NuevoTramoDeViaje(s.repo, c.PathInt("historia_id"), c.FormValue("texto"))
	if err != nil {
		return err
	}
	return c.RefreshHTMX()
}

func (s *servidor) deleteTramoDeViaje(c *gecko.Context) error {
	err := dhistorias.EliminarTramoDeViaje(s.repo, c.PathInt("historia_id"), c.PathInt("posicion"))
	if err != nil {
		return err
	}
	return c.RefreshHTMX()
}
