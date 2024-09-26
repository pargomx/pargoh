package ust

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
