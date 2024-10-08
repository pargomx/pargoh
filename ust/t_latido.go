package ust

import "errors"

// Latido corresponde a un elemento de la tabla 'latidos'.
type Latido struct {
	Timestamp  string // `latidos.timestamp`
	Segundos   int    // `latidos.segundos`
	ProyectoID string // `latidos.proyecto_id`
	HistoriaID *int   // `latidos.historia_id`
}

var (
	ErrLatidoNotFound      error = errors.New("el latido gestión no se encuentra")
	ErrLatidoAlreadyExists error = errors.New("el latido gestión ya existe")
)

func (lat *Latido) Validar() error {

	return nil
}
