package ust

import (
	"time"

	"github.com/pargomx/gecko/gko"
	"github.com/pargomx/gecko/gkt"
)

// El día empieza a las 6am y termina a las 5:59 del día siguiente.
func MinutosSince6am(fechaHora string) int {
	fecha, err := gkt.ToFechaHora(fechaHora)
	if err != nil {
		gko.LogError(err)
		return 0
	}
	sixAM := time.Date(fecha.Year(), fecha.Month(), fecha.Day(), 6, 0, 0, 0, gkt.TzMexico)
	if fecha.Hour() < 6 {
		sixAM = sixAM.AddDate(0, 0, -1) // antes 6am, contar desde 6:00am del día anterior.
	}
	return int(fecha.Sub(sixAM).Minutes())
}
