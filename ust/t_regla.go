package ust

import "errors"

// Regla corresponde a un elemento de la tabla 'reglas'.
type Regla struct {
	HistoriaID   int    // `reglas.historia_id`
	Posicion     int    // `reglas.posicion`  Posición consecutiva con respecto a sus hermanos
	Texto        string // `reglas.texto`
	Implementada bool   // `reglas.implementada`
	Probada      bool   // `reglas.probada`
}

var (
	ErrReglaNotFound      error = errors.New("la regla no se encuentra")
	ErrReglaAlreadyExists error = errors.New("la regla ya existe")
)

func (reg *Regla) Validar() error {

	return nil
}
