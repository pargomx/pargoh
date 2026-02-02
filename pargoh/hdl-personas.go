package main

import (
	"monorepo/dhistorias"

	"github.com/pargomx/gecko"
)

func (s *readhdl) getPersonaDoc(c *gecko.Context) error {
	Persona, err := s.repoOld.GetPersona(c.PathInt("persona_id"))
	if err != nil {
		return err
	}
	Proyecto, err := s.repoOld.GetProyecto(Persona.ProyectoID)
	if err != nil {
		return err
	}
	hists, err := s.repoOld.ListHistoriasByPadreID(Persona.PersonaID)
	if err != nil {
		return err
	}
	Historias := make(dhistorias.HistoriaAgregadoList, len(hists))
	for i, h := range hists {
		agg, err := dhistorias.GetHistoria(h.HistoriaID, dhistorias.GetDescendientes|dhistorias.GetTramos|dhistorias.GetReglas|dhistorias.GetRelacionadas, s.repoOld)
		if err != nil {
			return err
		}
		Historias[i] = *agg
	}
	data := map[string]any{
		"Titulo":    Persona.Nombre + " - " + Proyecto.Titulo,
		"Persona":   Persona,
		"Proyecto":  Proyecto,
		"Historias": Historias,
	}
	return c.Render(200, "persona_doc", data)
}

func (s *readhdl) getPersonaDebug(c *gecko.Context) error {
	Persona, err := s.repoOld.GetPersona(c.PathInt("persona_id"))
	if err != nil {
		return err
	}
	Proyecto, err := s.repoOld.GetProyecto(Persona.ProyectoID)
	if err != nil {
		return err
	}
	type HistoriaDebug struct {
		Agg dhistorias.HistoriaAgregado
		Rec dhistorias.HistoriaRecursiva
	}
	HistoriasRec, err := dhistorias.GetHistoriasDescendientes(Persona.PersonaID, 0, dhistorias.GetReglas|dhistorias.GetTareas, s.repoOld)
	if err != nil {
		return err
	}
	Historias := make([]HistoriaDebug, len(HistoriasRec))
	for i, h := range HistoriasRec {
		agg, err := dhistorias.GetHistoria(h.HistoriaID, dhistorias.GetDescendientes, s.repoOld)
		if err != nil {
			return err
		}
		Historias[i] = HistoriaDebug{
			Agg: *agg,
			Rec: h,
		}
	}
	data := map[string]any{
		"Titulo":    Persona.Nombre + " - Debug historias",
		"Persona":   Persona,
		"Proyecto":  Proyecto,
		"Historias": Historias,
	}
	return c.RenderOk("persona_debug", data)
}
