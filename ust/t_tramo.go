package ust

import "errors"

// Tramo corresponde a un elemento de la tabla 'tramos'.
type Tramo struct {
	HistoriaID int    // `tramos.historia_id`
	Posicion   int    // `tramos.posicion`  Posición consecutiva con respecto a sus nodos hermanos
	Texto      string // `tramos.texto`
	Imagen     string // `tramos.imagen`
}

var (
	ErrTramoNotFound      error = errors.New("el tramo no se encuentra")
	ErrTramoAlreadyExists error = errors.New("el tramo ya existe")
)

func (tra *Tramo) Validar() error {

	return nil
}
