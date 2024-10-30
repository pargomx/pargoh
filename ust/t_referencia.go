package ust

import "errors"

// Referencia corresponde a un elemento de la tabla 'referencias'.
type Referencia struct {
	HistoriaID    int // `referencias.historia_id`
	RefHistoriaID int // `referencias.ref_historia_id`
}

var (
	ErrReferenciaNotFound      error = errors.New("la referencia no se encuentra")
	ErrReferenciaAlreadyExists error = errors.New("la referencia ya existe")
)

func (ref *Referencia) Validar() error {

	return nil
}
