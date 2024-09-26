package ust

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
