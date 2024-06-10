package ust

import "errors"

// Persona corresponde a un elemento de la tabla 'personas'.
type Persona struct {
	PersonaID   int    // `personas.persona_id`
	Nombre      string // `personas.nombre`
	Descripcion string // `personas.descripcion`
}

var (
	ErrPersonaNotFound      error = errors.New("la persona del dominio no se encuentra")
	ErrPersonaAlreadyExists error = errors.New("la persona del dominio ya existe")
)

func (per *Persona) Validar() error {

	return nil
}
