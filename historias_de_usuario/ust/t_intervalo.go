package ust

import "errors"

// Intervalo corresponde a un elemento de la tabla 'intervalos'.
type Intervalo struct {
	TareaID int    // `intervalos.tarea_id`
	Inicio  string // `intervalos.inicio`
	Fin     string // `intervalos.fin`
}

var (
	ErrIntervaloNotFound      error = errors.New("el Intervalo no se encuentra")
	ErrIntervaloAlreadyExists error = errors.New("el Intervalo ya existe")
)

func (interv *Intervalo) Validar() error {

	return nil
}
