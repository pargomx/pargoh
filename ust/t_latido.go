package ust

import (
	"time"

	"github.com/pargomx/gecko/gko"
	"github.com/pargomx/gecko/gkt"
)

// Latido corresponde a un elemento de la tabla 'latidos'.
type Latido struct {
	Timestamp string // `latidos.timestamp`
	PersonaID int    // `latidos.persona_id`
	Segundos  int    // `latidos.segundos`
}

// Cantidad de minutos transcurridos desde las 6am del día de trabajo.
func (i *Latido) MinutosSince6am() int {
	return MinutosSince6am(i.Timestamp)
}

func (i *Latido) Time() time.Time {
	tm, err := gkt.ToFechaHora(i.Timestamp)
	if err != nil {
		gko.Err(err).Op("latido.Time").Log()
		return time.Time{}
	}
	return tm
}
