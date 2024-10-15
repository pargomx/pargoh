package dhistorias

import (
	"math"
	"monorepo/ust"
)

type HistoriaAgregado struct {
	Historia ust.NodoHistoria // Historia raíz
	Persona  ust.Persona
	Proyecto ust.Proyecto
	Tareas   []ust.Tarea
	Tramos   []ust.Tramo
	Reglas   []ust.Regla

	Ancestros     []ust.NodoHistoria
	Descendientes []HistoriaRecursiva
}

// ================================================================ //

// Tiempo estimado solo para las tareas de la historia raíz.
func (h *HistoriaAgregado) SegundosEstimadoTareas() (total int) {
	for _, t := range h.Tareas {
		total += t.SegundosEstimado
	}
	return total
}

// Tiempo real solo para las tareas de la historia raíz.
func (h *HistoriaAgregado) SegundosRealTareas() (total int) {
	for _, t := range h.Tareas {
		total += t.SegundosReal
	}
	return total
}

// Puntaje obtenido de las tareas en la historia raíz.
func (h *HistoriaAgregado) ValorPonderadoTotal() (total int) {
	for _, t := range h.Tareas {
		total += t.ValorPonderado()
	}
	return total
}

// Puntaje obtenido de las tareas en la historia raíz.
func (h *HistoriaAgregado) AvancePonderadoTotal() (total int) {
	for _, t := range h.Tareas {
		total += t.AvancePonderado()
	}
	return total
}

// Relación entre ValorPonderado y AvancePonderado
// obtenido de las tareas en la historia raíz.
func (h *HistoriaAgregado) AvancePorcentual() float64 {
	if h.ValorPonderadoTotal() == 0 {
		return 0
	}
	return math.Round(
		float64(h.AvancePonderadoTotal())/
			float64(h.ValorPonderadoTotal())*
			10*100) / 10
}

// ================================================================ //

// Relación entre ValorPonderado de una tarea y el ValorPonderadoTotal
// de todas las tareas de la historia raíz.
func (h *HistoriaAgregado) ValorPorcentual(tareaID int) float64 {
	for _, t := range h.Tareas {
		if t.TareaID == tareaID {
			vIndividual := float64(t.ValorPonderado())
			vTodas := float64(h.ValorPonderadoTotal())
			if vTodas == 0 {
				return 0
			}
			return math.Round(vIndividual/vTodas*10*100) / 10
		}
	}
	return 0
}
