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

func (h *Historia) SegundosEstimadoTareas() int {
	total := 0
	for _, t := range h.Tareas {
		total += t.SegundosEstimado
	}
	return total
}

func (h *Historia) SegundosRealTareas() int {
	total := 0
	for _, t := range h.Tareas {
		total += t.SegundosReal
	}
	return total
}

func (h *HistoriaRecursiva) SegundosEstimadoTareas() int {
	total := h.Segundos
	for _, h := range h.Descendientes {
		total += h.SegundosEstimadoTareas()
	}
	return total
}

// ================================================================ //

func (h *HistoriaRecursiva) SegundosTranscTotal() int {
	total := h.Segundos
	for _, h := range h.Descendientes {
		total += h.SegundosTranscTotal()
	}
	return total
}
func (h *HistoriaRecursiva) SegundosAvanceTeorico() int {
	return h.AvancePorcentual() * h.SegundosPresupuesto / 100
}
func (h *HistoriaRecursiva) AvancePorcentual() int {
	return 75
}

func (h *HistoriaRecursiva) HorasEstimado() float64 {
	return math.Round(float64(h.SegundosPresupuesto)/3600*100) / 100
}
func (h *HistoriaRecursiva) HorasTranscTotal() float64 {
	return math.Round(float64(h.SegundosTranscTotal())/3600*100) / 100
}
func (h *HistoriaRecursiva) HorasAvanceTeorico() float64 {
	return math.Round(float64(h.SegundosAvanceTeorico())/3600*100) / 100
}

// ================================================================ //

func (h *Historia) ValorPonderadoTotal() int {
	total := 0
	for _, t := range h.Tareas {
		total += t.ValorPonderado()
	}
	return total
}

func (h *Historia) AvancePonderadoTotal() int {
	total := 0
	for _, t := range h.Tareas {
		total += t.AvancePonderado()
	}
	return total
}

func (h *Historia) ValorPorcentual(tareaID int) float64 {
	for _, t := range h.Tareas {
		if t.TareaID == tareaID {
			valor := float64(t.ValorPonderado())
			total := float64(h.ValorPonderadoTotal())
			if total == 0 {
				return 0
			}
			return math.Round(valor/total*100*100) / 100
		}
	}
	return 0
}

func (h *Historia) AvancePorcentual() float64 {
	valor := float64(h.ValorPonderadoTotal())
	avance := float64(h.AvancePonderadoTotal())
	if valor == 0 {
		return 0
	}
	return math.Round(avance/valor*100*100) / 100
}
