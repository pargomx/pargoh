package ust

import (
	"time"

	"github.com/pargomx/gecko/gko"
)

// ================================================================ //
// ========== PERSONA ============================================= //

type PersonaCosto struct {
	Persona
	Historias []HistoriaCosto
}

func (p PersonaCosto) TiempoEstimado() string {
	suma := 0
	for _, h := range p.Historias {
		suma += h.MinutosEstimado
	}
	return MinutosToString(suma)
}

func (p PersonaCosto) TiempoReal() string {
	suma := 0
	for _, h := range p.Historias {
		suma += h.SegundosReal
	}
	return SegundosToString(suma)
}

func (p PersonaCosto) HistCompletadas() (res []HistoriaCosto) {
	for _, h := range p.Historias {
		if h.Completada {
			res = append(res, h)
		}
	}
	return res
}

func (p PersonaCosto) HistNoCompletadas() (res []HistoriaCosto) {
	for _, h := range p.Historias {
		if !h.Completada {
			res = append(res, h)
		}
	}
	return res
}

func (p PersonaCosto) NumHistCompletadas() (res int) {
	for _, h := range p.Historias {
		if h.Completada {
			res++
		}
	}
	return res
}

func (p PersonaCosto) NumHistNoCompletadas() (res int) {
	for _, h := range p.Historias {
		if !h.Completada && h.Prioridad > 0 {
			res++
		}
	}
	return res
}

// ================================================================ //
// ========== HISTORIA ============================================ //

type HistoriaCosto struct {
	HistoriaID      int
	PadreID         int
	Nivel           int
	Posicion        int
	Titulo          string
	Prioridad       int
	Completada      bool
	MinutosEstimado int
	SegundosReal    int
}

func (h *HistoriaCosto) TiempoEstimado() string {
	return MinutosToString(h.MinutosEstimado)
}

func (h *HistoriaCosto) TiempoReal() string {
	return SegundosToString(h.SegundosReal)
}

// ================================================================ //
// ========== INTERVALOS ========================================== //

// TODO: revisar error y quizá hacer configurable. También en dhistorias/e_tarea-w.go.
var locationMexicoCity, _ = time.LoadLocation("America/Mexico_City")

func (itv Intervalo) Segundos() int {
	inicio, err := time.ParseInLocation("2006-01-02 15:04:05", itv.Inicio, locationMexicoCity)
	if err != nil {
		gko.Err(err).Op("IntervaloEnDia.ParseInicio").Ctx("string", itv.Inicio).Log()
	}
	var fin time.Time
	if itv.Fin == "" {
		fin = time.Now().In(locationMexicoCity)
	} else {
		fin, err = time.ParseInLocation("2006-01-02 15:04:05", itv.Fin, locationMexicoCity)
		if err != nil {
			gko.Err(err).Op("IntervaloEnDia.ParseFin").Ctx("string", itv.Fin).Log()
		}
	}
	return int(fin.Sub(inicio).Seconds())
}

func (itv IntervaloEnDia) Segundos() int {
	inicio, err := time.ParseInLocation("2006-01-02 15:04:05", itv.Inicio, locationMexicoCity)
	if err != nil {
		gko.Err(err).Op("IntervaloEnDia.ParseInicio").Ctx("string", itv.Inicio).Log()
	}
	var fin time.Time
	if itv.Fin == "" {
		fin = time.Now().In(locationMexicoCity)
	} else {
		fin, err = time.ParseInLocation("2006-01-02 15:04:05", itv.Fin, locationMexicoCity)
		if err != nil {
			gko.Err(err).Op("IntervaloEnDia.ParseFin").Ctx("string", itv.Fin).Log()
		}
	}
	return int(fin.Sub(inicio).Seconds())
}

func (itv Intervalo) Duracion() string {
	return SegundosToString(itv.Segundos())
}
