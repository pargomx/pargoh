package dhistorias

import (
	"math"
	"monorepo/ust"
)

type HistoriaRecursiva struct {
	ust.NodoHistoria
	Tareas        TareasList
	Descendientes []HistoriaRecursiva
}

// ================================================================ //

// Tiempo estimado para todas las tareas de la historia raíz y todas sus descendientes.
func (h *HistoriaRecursiva) SegundosEstimadoMust() (total int) {
	if h.EsPrioridadMust() || h.Completada {
		total = h.SegundosEstimado
	}
	for _, h := range h.Descendientes {
		total += h.SegundosEstimadoMust()
	}
	return total
}

// Tiempo real para todas las tareas de la historia raíz y todas sus descendientes.
func (h *HistoriaRecursiva) SegundosUtilizadoMust() (total int) {
	if h.EsPrioridadMust() || h.Completada {
		total = h.SegundosUtilizado
	}
	for _, h := range h.Descendientes {
		total += h.SegundosUtilizadoMust()
	}
	return total
}

// Puntaje obtenido de las tareas en todas las historias del árbol.
func (h *HistoriaRecursiva) ValorPonderado() (total int) {
	if h.EsPrioridadMust() || h.Completada {
		for _, t := range h.Tareas {
			total += t.ValorPonderado()
		}
	}
	for _, h := range h.Descendientes {
		total += h.ValorPonderado()
	}
	return total
}

// Puntaje obtenido de las tareas en todas las historias del árbol.
func (h *HistoriaRecursiva) AvancePonderado() (total int) {
	if h.EsPrioridadMust() || h.Completada {
		for _, t := range h.Tareas {
			total += t.AvancePonderado()
		}
	}
	for _, h := range h.Descendientes {
		total += h.AvancePonderado()
	}
	return total
}

// Relación entre ValorPonderado y AvancePonderado
func (h *HistoriaRecursiva) AvancePorcentual() float64 {
	if h.ValorPonderado() == 0 {
		return 0
	}
	return math.Round(
		float64(h.AvancePonderado())/
			float64(h.ValorPonderado())*
			10*100) / 10
}

// ================================================================ //

func (h *HistoriaRecursiva) SegundosAvanceTeorico() int {
	return int(h.AvancePorcentual() * float64(h.SegundosPresupuesto) / 100)
}

// ================================================================ //

func (h *HistoriaRecursiva) HorasEstimado() float64 {
	return math.Round(float64(h.SegundosPresupuesto)/3600*100) / 100
}

func (h *HistoriaRecursiva) HorasUtilizado() float64 {
	return math.Round(float64(h.SegundosUtilizadoMust())/3600*100) / 100
}

func (h *HistoriaRecursiva) HorasAvanceTeorico() float64 {
	return math.Round(float64(h.SegundosAvanceTeorico())/3600*100) / 100
}
