package ust

import "errors"

// Historia corresponde a un elemento de la tabla 'historias'.
type Historia struct {
	HistoriaID      int    // `historias.historia_id`
	Titulo          string // `historias.titulo`
	Objetivo        string // `historias.objetivo`
	Prioridad       int    // `historias.prioridad`
	Completada      bool   // `historias.completada`
	PersonaID       int    // `historias.persona_id`  Para índice
	ProyectoID      string // `historias.proyecto_id`  Para índice
	MinutosEstimado int    // `historias.minutos_estimado`  Tiempo estimado en minutos para implementar la historia de usuario en su totalidad
}

var (
	ErrHistoriaNotFound      error = errors.New("la historia de usuario no se encuentra")
	ErrHistoriaAlreadyExists error = errors.New("la historia de usuario ya existe")
)

func (his *Historia) Validar() error {

	return nil
}
