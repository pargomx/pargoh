package main

import (
	"monorepo/arbol"
	"monorepo/dhistorias"
	"monorepo/ust"

	"github.com/pargomx/gecko"
)

func (s *readhdl) getPersona(c *gecko.Context) error {
	Per, err := s.repo.GetPersona(c.PathInt("persona_id"))
	if err != nil {
		return err
	}

	Persona, err := s.repoOld.GetPersona(c.PathInt("persona_id"))
	if err != nil {
		return err
	}
	Proyecto, err := s.repoOld.GetProyecto(Persona.ProyectoID)
	if err != nil {
		return err
	}
	// Historias, err := dhistorias.GetHistoriasDescendientes(Persona.PersonaID, 0, s.repo)
	hists, err := s.repoOld.ListHistoriasByPadreID(Persona.PersonaID)
	if err != nil {
		return err
	}
	Historias := make(dhistorias.HistoriaAgregadoList, len(hists))
	for i, h := range hists {
		agg, err := dhistorias.GetHistoria(h.HistoriaID, dhistorias.GetDescendientes, s.repoOld)
		if err != nil {
			return err
		}
		Historias[i] = *agg
	}
	TareasEnCurso, err := s.repoOld.ListTareasEnCurso()
	if err != nil {
		return err
	}
	data := map[string]any{
		"Titulo":        Persona.Nombre + " - " + Proyecto.Titulo,
		"Persona":       Per,
		"Proyecto":      Proyecto,
		"Historias":     Historias,
		"TareasEnCurso": TareasEnCurso,
	}
	return c.RenderOk("persona", data)
}

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

func (s *writehdl) postPersona(c *gecko.Context, tx *handlerTx) error {
	args := arbol.ArgsAgregarHoja{
		Tipo:    "PER",
		NodoID:  ust.NewRandomID(),
		PadreID: c.FormInt("proyecto_id"),
		Titulo:  c.FormVal("nombre"),
	}
	err := tx.app.AgregarHoja(args)
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
	err := dhistorias.ActualizarPersona(persona, s.repoOld)
	if err != nil {
		return err
	}
	return c.RefreshHTMX()
}

func (s *writehdl) patchPersona(c *gecko.Context, tx *handlerTx) error {
	err := tx.app.ParcharNodo(arbol.ArgsParcharNodo{
		NodoID: c.PathInt("persona_id"),
		Campo:  c.PathVal("param"),
		NewVal: c.FormValue("value"),
	})
	if err != nil {
		return err
	}
	return c.RefreshHTMX()
}
