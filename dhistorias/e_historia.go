package dhistorias

import (
	"math"
	"monorepo/ust"
)

type Historia struct {
	Historia ust.NodoHistoria
	Persona  ust.Persona
	Proyecto ust.Proyecto
	Tareas   []ust.Tarea
	Tramos   []ust.Tramo
	Reglas   []ust.Regla

	Ancestros     []ust.NodoHistoria
	Descendientes []HistoriaRecursiva
}

type HistoriaRecursiva struct {
	ust.NodoHistoria
	Descendientes []HistoriaRecursiva
}

// ================================================================ //

const prioridadInvalidaMsg = "La prioridad debe estar entre 0 y 3"

func prioridadValida(prioridad int) bool {
	return prioridad >= 0 && prioridad <= 3
}

// ================================================================ //

func (h *Historia) TiempoEstimado() int {
	total := 0
	for _, t := range h.Tareas {
		total += t.TiempoEstimado
	}
	return total
}

func (h *Historia) TiempoReal() int {
	total := 0
	for _, t := range h.Tareas {
		total += t.TiempoReal
	}
	return total
}

func (h *HistoriaRecursiva) TiempoEstimado() int {
	total := h.Segundos
	for _, h := range h.Descendientes {
		total += h.TiempoEstimado()
	}
	return total
}

// ================================================================ //

func (h *HistoriaRecursiva) SegundosEstimado() int {
	return h.MinutosEstimado * 60
}
func (h *HistoriaRecursiva) SegundosTranscTotal() int {
	total := h.Segundos
	for _, h := range h.Descendientes {
		total += h.SegundosTranscTotal()
	}
	return total
}
func (h *HistoriaRecursiva) SegundosAvanceTeorico() int {
	return h.AvancePorcentual() * h.SegundosEstimado() / 100
}
func (h *HistoriaRecursiva) AvancePorcentual() int {
	return 75
}

func (h *HistoriaRecursiva) HorasEstimado() float64 {
	return math.Round(float64(h.SegundosEstimado())/3600*100) / 100
}
func (h *HistoriaRecursiva) HorasTranscTotal() float64 {
	return math.Round(float64(h.SegundosTranscTotal())/3600*100) / 100
}
func (h *HistoriaRecursiva) HorasAvanceTeorico() float64 {
	return math.Round(float64(h.SegundosAvanceTeorico())/3600*100) / 100
}
