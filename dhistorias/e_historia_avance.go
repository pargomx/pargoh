package dhistorias

import "github.com/pargomx/gecko/gko"

// Para gráfico SVG de avance.
type avanceEscalado struct {
	Presupuesto int
	Estimado    int
	Utilizado   int
	Expectativa int
	Separadores []separadorHora
}

type separadorHora struct {
	Hora          int  // número de hora
	Posicion      int  // a escala
	EsPresupuesto bool // si este separador es el presupuesto
	EsUltimo      bool // si es el último separador
}

// Medidas de avance, presupuesto, estimado y utilizado relativas a mil.
func (h *HistoriaAgregado) AvanceRelativoMil() avanceEscalado {
	if h.avance != nil {
		return *h.avance
	}
	maximo := h.SegundosMaxPresupuestoEstimadoUtilizado()
	const mil = 1000 // ancho esperado
	h.avance = &avanceEscalado{
		Presupuesto: h.SegundosPresupuesto() * mil / maximo,
		Estimado:    h.SegundosEstimado() * mil / maximo,
		Utilizado:   h.SegundosUtilizado() * mil / maximo,
		Expectativa: h.SegundosExpectativaAvancePresupuesto() * mil / maximo,
	}

	// Separadores por hora
	horas := maximo / 3600
	gko.LogInfo("Horas ", horas, " max ", maximo)
	if maximo%3600 > 150 { // si se pasa más de 150 segundos, redondear hacia siguiente hora
		horas = (maximo / 3600) + 1
		gko.LogInfo("Redondeo ", horas)
	}
	horasPresupuesto := h.SegundosPresupuesto() / 3600
	for hora := range horas + 1 {
		if hora == 0 {
			continue
		}
		// Siempre agregar hora de presupuesto
		if hora == horasPresupuesto {
			h.avance.Separadores = append(h.avance.Separadores, separadorHora{
				Hora:          hora,
				Posicion:      hora * 3600 * mil / maximo,
				EsPresupuesto: true,
			})
			continue
		}
		// mostrar cada 5, 3, o 2 horas si son muchas
		if horas >= 90 {
			if hora%10 != 0 {
				continue
			}
		} else if horas >= 60 {
			if hora%6 != 0 {
				continue
			}
		} else if horas >= 30 {
			if hora%5 != 0 {
				continue
			}
		} else if horas >= 21 {
			if hora%3 != 0 {
				continue
			}
		} else if horas >= 10 {
			if hora%2 != 0 {
				continue
			}
		}
		h.avance.Separadores = append(h.avance.Separadores, separadorHora{
			Hora:     hora,
			Posicion: hora * 3600 * mil / maximo,
		})
	}
	if len(h.avance.Separadores) > 0 {
		h.avance.Separadores[len(h.avance.Separadores)-1].EsUltimo = true
	}
	return *h.avance
}

func (h *HistoriaAgregado) SegundosMaxPresupuestoEstimadoUtilizado() int {
	max := h.SegundosPresupuesto()
	est := h.SegundosEstimado()
	uti := h.SegundosUtilizado()
	if est > max {
		max = est
	}
	if uti > max {
		max = uti
	}
	return max
}
