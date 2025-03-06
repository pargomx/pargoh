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

	// Limitar acceso a proyectos...
	ses, ok := c.Sesion.(*Sesion)
	if !ok {
		return gko.ErrDatoInvalido().Msg("Sesión inválida")
	}
	if ses.Usuario != s.cfg.adminUser {
		pry, err := s.repo.GetProyecto(ses.Usuario)
		if err != nil {
			gko.Err(err).Strf("usuario '%v' no correspone a ningún proyecto", ses.Usuario).ErrNoAutorizado().Log()
			return c.RedirFull("/logout")
		}
		Proyectos = []ust.Proyecto{*pry}
	}

	res := make([]Pry, len(Proyectos))
	for i, p := range Proyectos {
		res[i].Proyecto = p
		res[i].Personas, err = s.repo.ListNodosPersonas(p.ProyectoID)
		if err != nil {
			return err
		}
	}
	Bugs, err := s.repo.ListTareasBugs()
	if err != nil {
		return err
	}
	TareasEnCurso, err := s.repo.ListTareasEnCurso()
	if err != nil {
		return err
	}
	data := map[string]any{
		"Titulo":        "Pargo",
		"Proyectos":     res,
		"Bugs":          Bugs,
		"TareasEnCurso": TareasEnCurso,
	}
	return c.RenderOk("proyectos", data)
}

func (s *servidor) postProyecto(c *gecko.Context) error {
	err := dhistorias.NuevoProyecto(c.FormVal("clave"), c.FormVal("titulo"), c.FormVal("descripcion"), s.repo)
	if err != nil {
		return err
	}
	return c.RedirOtro("/")
}

func (s *servidor) updateProyecto(c *gecko.Context) error {
	err := dhistorias.ModificarProyecto(c.PathVal("proyecto_id"), ust.Proyecto{
		ProyectoID:  c.FormVal("clave"),
		Posicion:    c.FormInt("posicion"),
		Titulo:      c.FormVal("titulo"),
		Color:       c.FormVal("color"),
		Descripcion: c.FormVal("descripcion"),
	}, s.repo)
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
	return c.AskedFor()
}

func (s *servidor) patchProyecto(c *gecko.Context) error {
	err := dhistorias.ParcharProyecto(
		c.PathVal("proyecto_id"),
		c.PathVal("param"),
		c.FormValue("value"),
		s.repo,
	)
	if err != nil {
		return err
	}
	return c.RefreshHTMX()
}

func (s *servidor) deleteProyecto(c *gecko.Context) error {
	err := dhistorias.EliminarProyecto(c.PathVal("proyecto_id"), s.repo)
	if err != nil {
		return err
	}
	return c.RedirOtro("/")
}

func (s *servidor) deleteProyectoPorCompleto(c *gecko.Context) error {
	pry, err := dhistorias.GetProyectoExport(c.PathVal("proyecto_id"), s.repo)
	if err != nil {
		return err
	}
	if c.PromptVal() != "eliminar_"+pry.Proyecto.ProyectoID {
		return gko.ErrDatoInvalido().Msg("No se confirmó la eliminación")
	}
	err = pry.EliminarPorCompleto(s.repo)
	if err != nil {
		return err
	}
	return c.RedirOtro("/")
}

func (s *servidor) postTimeGestion(c *gecko.Context) error {
	err := s.timeTracker.AddTimeSpent(c.PathInt("persona_id"), c.PathInt("seg"))
	if err != nil {
		return err
	}
	return c.StringOk("ok")
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
	TareasEnCurso, err := s.repo.ListTareasEnCurso()
	if err != nil {
		return err
	}
	data := map[string]any{
		"Titulo":        Proyecto.Titulo,
		"Proyecto":      Proyecto,
		"Personas":      Personas,
		"Proyectos":     Proyectos, // Para cambiar de proyecto a una persona.
		"TareasEnCurso": TareasEnCurso,
	}
	return c.RenderOk("proyecto", data)
}

func (s *servidor) getDocumentacionProyecto(c *gecko.Context) error {
	Proyecto, err := s.repo.GetProyecto(c.PathVal("proyecto_id"))
	if err != nil {
		return err
	}
	Personas, err := s.repo.ListNodosPersonas(Proyecto.ProyectoID)
	if err != nil {
		return err
	}
	type Personaje struct {
		Persona   ust.NodoPersona
		Historias []ust.Historia
	}
	Personajes := make([]Personaje, len(Personas))
	for i, p := range Personas {
		hists, err := s.repo.ListHistoriasByPadreID(p.PersonaID)
		if err != nil {
			return err
		}
		Personajes[i] = Personaje{
			Persona:   Personas[i],
			Historias: hists,
		}
	}
	data := map[string]any{
		"Titulo":     Proyecto.Titulo,
		"Proyecto":   Proyecto,
		"Personajes": Personajes,
	}
	return c.RenderOk("proyecto_doc", data)
}
