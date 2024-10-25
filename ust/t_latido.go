package ust

import "errors"

// Latido corresponde a un elemento de la tabla 'latidos'.
type Latido struct {
	Timestamp string // `latidos.timestamp`
	PersonaID int    // `latidos.persona_id`
	Segundos  int    // `latidos.segundos`
}

var (
	ErrLatidoNotFound      error = errors.New("el latido gestión no se encuentra")
	ErrLatidoAlreadyExists error = errors.New("el latido gestión ya existe")
)

func (lat *Latido) Validar() error {

	return nil
}
