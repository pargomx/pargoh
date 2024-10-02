package main

import (
	"monorepo/dhistorias"

	"github.com/pargomx/gecko"
)

func (s *servidor) navDesdeRoot(c *gecko.Context) error {
	proyectos, err := s.repo.ListProyectos()
	if err != nil {
		return err
	}
	data := map[string]any{
		"Proyectos": proyectos,
	}
	return c.RenderOk("nav_root", data)
}

func (s *servidor) navDesdeProyecto(c *gecko.Context) error {
	proyecto, err := s.repo.GetProyecto(c.PathVal("proyecto_id"))
	if err != nil {
		return err
	}
	personas, err := s.repo.ListNodosPersonas(proyecto.ProyectoID)
	if err != nil {
		return err
	}
	data := map[string]any{
		"Proyecto": proyecto,
		"Personas": personas,
	}
	return c.RenderOk("nav_proy", data)
}

func (s *servidor) navDesdePersona(c *gecko.Context) error {
	persona, err := s.repo.GetPersona(c.PathInt("persona_id"))
	if err != nil {
		return err
	}
	proyecto, err := s.repo.GetProyecto(persona.ProyectoID)
	if err != nil {
		return err
	}
	historias, err := s.repo.ListNodoHistoriasByPadreID(persona.PersonaID)
	if err != nil {
		return err
	}
	data := map[string]any{
		"Proyecto":  proyecto,
		"Persona":   persona,
		"Historias": historias,
	}
	return c.RenderOk("nav_pers", data)
}

func (s *servidor) navDesdeHistoria(c *gecko.Context) error {
	agg, err := dhistorias.GetHistoria(c.PathInt("historia_id"), s.repo)
	if err != nil {
		return err
	}
	data := map[string]any{
		"Agregado": agg,
	}
	return c.RenderOk("nav_hist", data)
}
