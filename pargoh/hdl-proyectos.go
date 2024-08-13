package main

import (
	"monorepo/dhistorias"

	"github.com/pargomx/gecko"
)

func (s *servidor) listaProyectos(c *gecko.Context) error {
	Proyectos, err := s.repo.ListProyectos()
	if err != nil {
		return err
	}
	data := map[string]any{
		"Proyectos": Proyectos,
	}
	return c.RenderOk("proyectos", data)
}

func (s *servidor) postProyecto(c *gecko.Context) error {
	err := dhistorias.NuevoProyecto(c.FormVal("clave"), c.FormVal("titulo"), c.FormVal("descripcion"), s.repo)
	if err != nil {
		return err
	}
	return c.RefreshHTMX()
}

func (s *servidor) updateProyecto(c *gecko.Context) error {
	err := dhistorias.ModificarProyecto(c.PathVal("proyecto_id"), c.FormVal("clave"), c.FormVal("titulo"), c.FormVal("descripcion"), s.repo)
	if err != nil {
		return err
	}
	return c.RefreshHTMX()
}

func (s *servidor) deleteProyecto(c *gecko.Context) error {
	err := dhistorias.QuitarProyecto(c.PathVal("proyecto_id"), s.repo)
	if err != nil {
		return err
	}
	return c.RefreshHTMX()
}

func (s *servidor) getProyecto(c *gecko.Context) error {
	Proyecto, err := s.repo.GetProyecto(c.PathVal("proyecto_id"))
	if err != nil {
		return err
	}
	Personas, err := s.repo.ListNodosPersonasByProyecto(Proyecto.ProyectoID)
	if err != nil {
		return err
	}
	Proyectos, err := s.repo.ListProyectos()
	if err != nil {
		return err
	}
	data := map[string]any{
		"Titulo":    Proyecto.Titulo,
		"Proyecto":  Proyecto,
		"Personas":  Personas,
		"Proyectos": Proyectos,
	}
	return c.RenderOk("personas", data)
}
