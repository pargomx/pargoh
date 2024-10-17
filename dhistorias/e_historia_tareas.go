package dhistorias

import (
	"math"
	"monorepo/ust"
)

type TareasList []ust.Tarea

// Tiempo estimado solo para las tareas de la historia raíz.
func (tareas TareasList) SegundosEstimado() (total int) {
	for _, t := range tareas {
		total += t.SegundosEstimado
	}
	return total
}

// Tiempo real solo para las tareas de la historia raíz.
func (tareas TareasList) SegundosReal() (total int) {
	for _, t := range tareas {
		total += t.SegundosReal
	}
	return total
}

// Puntaje obtenido de las tareas en la historia raíz.
func (tareas TareasList) ValorPonderado() (total int) {
	for _, t := range tareas {
		total += t.ValorPonderado()
	}
	return total
}

// Puntaje obtenido de las tareas en la historia raíz.
func (tareas TareasList) AvancePonderado() (total int) {
	for _, t := range tareas {
		total += t.AvancePonderado()
	}
	return total
}

// Relación entre ValorPonderado y AvancePonderado
// obtenido de las tareas en la historia raíz.
func (tareas TareasList) AvancePorcentual() float64 {
	if tareas.ValorPonderado() == 0 {
		return 0
	}
	return math.Round(
		float64(tareas.AvancePonderado())/
			float64(tareas.ValorPonderado())*
			10*100) / 10
}

// ================================================================ //

// Relación entre ValorPonderado de una tarea y el ValorPonderadoTotal
// de todas las tareas de la historia raíz.
func (tareas TareasList) ValorPorcentual(tareaID int) float64 {
	for _, t := range tareas {
		if t.TareaID == tareaID {
			vIndividual := float64(t.ValorPonderado())
			vTodas := float64(tareas.ValorPonderado())
			if vTodas == 0 {
				return 0
			}
			return math.Round(vIndividual/vTodas*10*100) / 10
		}
	}
	return 0
}
