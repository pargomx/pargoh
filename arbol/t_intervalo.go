package arbol

import (
	"time"

	"github.com/pargomx/gecko/gko"
	"github.com/pargomx/gecko/gkt"
)

// Intervalo corresponde a un elemento de la tabla 'intervalos'.
type Intervalo struct {
	NodoID int    // `intervalos.nodo_id`  ID del nodo
	TsIni  string // `intervalos.ts_ini`  Inicio del intervalo en hora local
	TsFin  string // `intervalos.ts_fin`  Fin del intervalo en hora local
}

func (itv Intervalo) TareaID() int {
	return itv.NodoID
}

// Segundos transcurridos entre el inicio y el fin del intervalo.
// Si aún no tiene fin, entonces entre el inicio y la hora actual.
func (itv Intervalo) Segundos() int {
	inicio, err := time.ParseInLocation(gkt.FormatoFechaHora, itv.TsIni, gkt.TzMexico)
	if err != nil {
		gko.Err(err).Op("Intervalo.ParseInicio").Ctx("string", itv.TsIni).Log()
	}
	var fin time.Time
	if itv.TsFin == "" {
		fin = time.Now().In(gkt.TzMexico)
	} else {
		fin, err = time.ParseInLocation(gkt.FormatoFechaHora, itv.TsFin, gkt.TzMexico)
		if err != nil {
			gko.Err(err).Op("Intervalo.ParseFin").Ctx("string", itv.TsFin).Log()
		}
	}
	return int(fin.Sub(inicio).Seconds())
}
