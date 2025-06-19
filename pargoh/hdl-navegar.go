package main

import (
	"monorepo/dhistorias"

	"github.com/pargomx/gecko"
)

func (s *readhdl) navDesdeRoot(c *gecko.Context) error {
	proyectos, err := s.repoOld.ListProyectos()
	if err != nil {
		return err
	}
	data := map[string]any{
		"Proyectos": proyectos,
	}
	return c.RenderOk("nav_root", data)
}

func (s *readhdl) navDesdeProyecto(c *gecko.Context) error {
	proyecto, err := s.repoOld.GetProyecto(c.PathVal("proyecto_id"))
	if err != nil {
		return err
	}
	personas, err := s.repoOld.ListNodosPersonas(proyecto.ProyectoID)
	if err != nil {
		return err
	}
	data := map[string]any{
		"Proyecto": proyecto,
		"Personas": personas,
	}
	return c.RenderOk("nav_proy", data)
}

func (s *readhdl) navDesdePersona(c *gecko.Context) error {
	persona, err := s.repoOld.GetPersona(c.PathInt("persona_id"))
	if err != nil {
		return err
	}
	proyecto, err := s.repoOld.GetProyecto(persona.ProyectoID)
	if err != nil {
		return err
	}
	historias, err := s.repoOld.ListNodoHistoriasByPadreID(persona.PersonaID)
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

func (s *readhdl) navDesdeHistoria(c *gecko.Context) error {
	agg, err := dhistorias.GetHistoria(c.PathInt("historia_id"), dhistorias.GetDescendientes, s.repoOld)
	if err != nil {
		return err
	}
	data := map[string]any{
		"Agregado": agg,
	}
	return c.RenderOk("nav_hist", data)
}
