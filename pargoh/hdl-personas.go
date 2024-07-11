package main

import (
	"monorepo/historias_de_usuario/dhistorias"
	"monorepo/historias_de_usuario/ust"

	"github.com/pargomx/gecko"
)

func (s *servidor) getPersonas(c *gecko.Context) error {
	Personas, err := s.repo.ListNodosPersonas()
	if err != nil {
		return err
	}
	data := map[string]any{
		"Titulo":   "Pargo - Personas",
		"Personas": Personas,
	}
	return c.RenderOk("personas", data)
}

func (s *servidor) postPersona(c *gecko.Context) error {
	persona := ust.Persona{
		PersonaID:   ust.NewPersonaID(),
		Nombre:      c.FormVal("nombre"),
		Descripcion: c.FormVal("descripcion"),
	}
	err := dhistorias.InsertarPersona(persona, s.repo)
	if err != nil {
		return err
	}
	return c.RefreshHTMX()
}

func (s *servidor) patchPersona(c *gecko.Context) error {
	persona := ust.Persona{
		PersonaID:   c.PathInt("persona_id"),
		Nombre:      c.FormVal("nombre"),
		Descripcion: c.FormVal("descripcion"),
	}
	err := dhistorias.ActualizarPersona(persona, s.repo)
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
