package main

import (
	"monorepo/dhistorias"
	"monorepo/ust"

	"github.com/pargomx/gecko"
)

func (s *servidor) getPersona(c *gecko.Context) error {
	Persona, err := s.repo.GetPersona(c.PathInt("persona_id"))
	if err != nil {
		return err
	}
	Proyecto, err := s.repo.GetProyecto(Persona.ProyectoID)
	if err != nil {
		return err
	}
	Historias, err := s.repo.ListNodoHistorias(Persona.PersonaID)
	if err != nil {
		return err
	}
	TareasEnCurso, err := s.repo.ListTareasEnCurso()
	if err != nil {
		return err
	}
	data := map[string]any{
		"Titulo":        "👤 " + Persona.Nombre + " - " + Proyecto.Titulo,
		"Persona":       Persona,
		"Proyecto":      Proyecto,
		"Historias":     Historias,
		"TareasEnCurso": TareasEnCurso,
	}
	return c.RenderOk("persona", data)
}

func (s *servidor) getPersonaCosto(c *gecko.Context) error {
	Persona, err := s.repo.GetPersona(c.PathInt("persona_id"))
	if err != nil {
		return err
	}
	Proyecto, err := s.repo.GetProyecto(Persona.ProyectoID)
	if err != nil {
		return err
	}
	Historias, err := s.repo.ListNodoHistorias(Persona.PersonaID)
	if err != nil {
		return err
	}
	TareasEnCurso, err := s.repo.ListTareasEnCurso()
	if err != nil {
		return err
	}
	HistoriasCosto, err := s.repo.ListHistoriasCosto(Persona.PersonaID)
	if err != nil {
		return err
	}
	PersonaCosto := ust.PersonaCosto{
		Persona:   *Persona,
		Historias: HistoriasCosto,
	}
	DiasTrabajo, err := s.repo.ListIntervalosEnDias()
	if err != nil {
		return err
	}
	MapaDiasSegundos := make(map[string]float64)
	for _, dia := range DiasTrabajo {
		MapaDiasSegundos[dia.Fecha] += float64(dia.Segundos()) / 60 / 60
		if MapaDiasSegundos[dia.Fecha] > 24 {
			MapaDiasSegundos[dia.Fecha] = 20
		}
	}

	data := map[string]any{
		"Titulo":           "👤 " + Persona.Nombre + " - " + Proyecto.Titulo,
		"Persona":          Persona,
		"Proyecto":         Proyecto,
		"Historias":        Historias,
		"TareasEnCurso":    TareasEnCurso,
		"PersonaCosto":     PersonaCosto,
		"DiasTrabajo":      DiasTrabajo,
		"MapaDiasSegundos": MapaDiasSegundos,
	}
	return c.RenderOk("persona", data)
}

func (s *servidor) postPersona(c *gecko.Context) error {
	persona := ust.Persona{
		PersonaID:   ust.NewPersonaID(),
		ProyectoID:  c.FormVal("proyecto_id"),
		Nombre:      c.FormVal("nombre"),
		Descripcion: c.FormVal("descripcion"),
	}
	err := dhistorias.InsertarPersona(persona, s.repo)
	if err != nil {
		return err
	}
	return c.RefreshHTMX()
}

func (s *servidor) updatePersona(c *gecko.Context) error {
	persona := ust.Persona{
		PersonaID:   c.PathInt("persona_id"),
		ProyectoID:  c.FormVal("proyecto_id"),
		Nombre:      c.FormVal("nombre"),
		Descripcion: c.FormVal("descripcion"),
	}
	err := dhistorias.ActualizarPersona(persona, s.repo)
	if err != nil {
		return err
	}
	return c.RefreshHTMX()
}

func (s *servidor) patchPersona(c *gecko.Context) error {
	err := dhistorias.ParcharPersona(
		c.PathInt("persona_id"),
		c.PathVal("param"),
		c.FormValue("value"),
		s.repo,
	)
	if err != nil {
		return err
	}
	return c.RefreshHTMX()
}

func (s *servidor) deletePersona(c *gecko.Context) error {
	err := dhistorias.EliminarPersona(c.PathInt("persona_id"), s.repo)
	if err != nil {
		return err
	}
	return c.RefreshHTMX()
}
