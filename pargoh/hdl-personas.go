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

func (s *servidor) getMétricas(c *gecko.Context) error {
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

	Intervalos, err := s.repo.ListIntervalosEnDias()
	if err != nil {
		return err
	}
	DiasTrabajoMapHoras := make(map[string]float64)
	for _, dia := range Intervalos {
		DiasTrabajoMapHoras[dia.Fecha] += float64(dia.Segundos()) / 60 / 60
	}

	IntervalosMap := make(map[string][]ust.IntervaloEnDia)
	for _, interv := range Intervalos {
		IntervalosMap[interv.Fecha] = append(IntervalosMap[interv.Fecha], interv)
	}

	Dias, err := s.repo.ListDias()
	if err != nil {
		return err
	}
	type DiaTrabajo struct {
		Fecha    string
		Segundos int
		Horas    float64
		Tareas   map[int]ust.Tarea
	}
	DiasTrabajo := make([]DiaTrabajo, len(Dias))
	for i, dia := range Dias {
		DiasTrabajo[i].Fecha = dia
		for _, interv := range IntervalosMap[dia] {
			DiasTrabajo[i].Segundos += interv.Segundos()
			if DiasTrabajo[i].Tareas == nil {
				DiasTrabajo[i].Tareas = make(map[int]ust.Tarea)
			}
			if tar, ok := DiasTrabajo[i].Tareas[interv.TareaID]; !ok {
				tarea, err := s.repo.GetTarea(interv.TareaID)
				if err != nil {
					return err
				}
				tarea.TiempoReal = interv.Segundos()
				DiasTrabajo[i].Tareas[interv.TareaID] = *tarea
			} else {
				tar.TiempoReal += interv.Segundos()
				DiasTrabajo[i].Tareas[interv.TareaID] = tar
			}
		}
		DiasTrabajo[i].Horas = float64(DiasTrabajo[i].Segundos) / 60 / 60
	}

	data := map[string]any{
		"Titulo":        "👤 " + Persona.Nombre + " - " + Proyecto.Titulo,
		"Persona":       Persona,
		"Proyecto":      Proyecto,
		"Historias":     Historias,
		"TareasEnCurso": TareasEnCurso,
		"PersonaCosto":  PersonaCosto,

		"DiasTrabajoMapHoras": DiasTrabajoMapHoras,
		"DiasTrabajo":         DiasTrabajo,
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
