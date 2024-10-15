package dhistorias

import (
	"math"
	"monorepo/ust"
)

type HistoriaRecursiva struct {
	ust.NodoHistoria
	Descendientes []HistoriaRecursiva
}

// ================================================================ //

// Tiempo estimado para todas las tareas de la historia raíz y todas sus descendientes.
func (h *HistoriaRecursiva) SegundosEstimadoTareas() int {
	total := h.SegundosEstimado
	for _, h := range h.Descendientes {
		total += h.SegundosEstimadoTareas()
	}
	return total
}

// Tiempo real para todas las tareas de la historia raíz y todas sus descendientes.
func (h *HistoriaRecursiva) SegundosTranscTotal() int {
	total := h.SegundosReal
	for _, h := range h.Descendientes {
		total += h.SegundosTranscTotal()
	}
	return total
}

func (h *HistoriaRecursiva) SegundosAvanceTeorico() int {
	return h.AvancePorcentual() * h.SegundosPresupuesto / 100
}

func (h *HistoriaRecursiva) AvancePorcentual() int {
	if h.SegundosEstimado == 0 {
		return 0
	}
	return h.SegundosReal * 100 / h.SegundosEstimado
}

// ================================================================ //

func (h *HistoriaRecursiva) HorasEstimado() float64 {
	return math.Round(float64(h.SegundosPresupuesto)/3600*100) / 100
}

func (h *HistoriaRecursiva) HorasTranscTotal() float64 {
	return math.Round(float64(h.SegundosTranscTotal())/3600*100) / 100
}

func (h *HistoriaRecursiva) HorasAvanceTeorico() float64 {
	return math.Round(float64(h.SegundosAvanceTeorico())/3600*100) / 100
}
