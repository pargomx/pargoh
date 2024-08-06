package ust

import "errors"

// Historia corresponde a un elemento de la tabla 'historias'.
type Historia struct {
	HistoriaID int    // `historias.historia_id`
	Titulo     string // `historias.titulo`
	Objetivo   string // `historias.objetivo`
	Prioridad  int    // `historias.prioridad`
	Completada bool   // `historias.completada`
}

var (
	ErrHistoriaNotFound      error = errors.New("la historia de usuario no se encuentra")
	ErrHistoriaAlreadyExists error = errors.New("la historia de usuario ya existe")
)

func (his *Historia) Validar() error {

	return nil
}
