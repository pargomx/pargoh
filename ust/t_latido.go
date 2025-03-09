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
	ini, err := gkt.ToFechaHora(i.Timestamp)
	if err != nil {
		gko.LogError(err)
		return 0
	}
	sixAM := time.Date(ini.Year(), ini.Month(), ini.Day(), 6, 0, 0, 0, gkt.TzMexico)
	return int(ini.Sub(sixAM).Minutes())
}
