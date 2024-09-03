package main

import (
	"monorepo/dhistorias"
	"monorepo/ust"
	"strings"

	"github.com/pargomx/gecko"
	"github.com/pargomx/gecko/gko"
)

func (s *servidor) listaProyectos(c *gecko.Context) error {
	type Pry struct {
		ust.Proyecto
		Personas []ust.NodoPersona
	}
	Proyectos, err := s.repo.ListProyectos()
	if err != nil {
		return err
	}
	res := make([]Pry, len(Proyectos))
	for i, p := range Proyectos {
		res[i].Proyecto = p
		res[i].Personas, err = s.repo.ListNodosPersonas(p.ProyectoID)
		if err != nil {
			return err
		}
	}
	data := map[string]any{
		"Proyectos": res,
	}
	return c.RenderOk("proyectos", data)
}

func (s *servidor) postProyecto(c *gecko.Context) error {
	err := dhistorias.NuevoProyecto(c.FormVal("clave"), c.FormVal("titulo"), c.FormVal("descripcion"), s.repo)
	if err != nil {
		return err
	}
	return c.Redir("/")
}

func (s *servidor) updateProyecto(c *gecko.Context) error {
	err := dhistorias.ModificarProyecto(c.PathVal("proyecto_id"), c.FormVal("clave"), c.FormVal("titulo"), c.FormVal("descripcion"), s.repo)
	if err != nil {
		return err
	}
	hdr, err := c.FormFile("imagen")
	if err == nil {
		file, err := hdr.Open()
		if err != nil {
			return err
		}
		defer file.Close()
		gko.LogDebugf("Imagen recibida: %v\t Tamaño: %v\t MIME:%v", hdr.Filename, hdr.Size, hdr.Header.Get("Content-Type"))
		err = dhistorias.SetImagenProyecto(c.PathVal("proyecto_id"), strings.TrimPrefix(hdr.Header.Get("Content-Type"), "image/"), file, s.cfg.imagesDir, s.repo)
		if err != nil {
			return err
		}
	}
	return c.Redir("/")
}

func (s *servidor) deleteProyecto(c *gecko.Context) error {
	err := dhistorias.QuitarProyecto(c.PathVal("proyecto_id"), s.repo)
	if err != nil {
		return err
	}
	return c.Redir("/")
}

func (s *servidor) getProyecto(c *gecko.Context) error {
	Proyecto, err := s.repo.GetProyecto(c.PathVal("proyecto_id"))
	if err != nil {
		return err
	}
	Personas, err := s.repo.ListNodosPersonas(Proyecto.ProyectoID)
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
		"Proyectos": Proyectos, // Para cambiar de proyecto a una persona.
	}
	return c.RenderOk("proyecto", data)
}
