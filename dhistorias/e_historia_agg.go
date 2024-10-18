package dhistorias

import (
	"math"
	"monorepo/ust"
)

type HistoriaAgregado struct {
	Historia ust.NodoHistoria // Historia raíz
	Persona  ust.Persona
	Proyecto ust.Proyecto
	Tareas   TareasList
	Tramos   []ust.Tramo
	Reglas   []ust.Regla

	Ancestros     []ust.NodoHistoria
	Descendientes []HistoriaRecursiva
}

// ================================================================ //

// Presupuesto agregando todo el árbol de historias.
func (h *HistoriaAgregado) SegundosPresupuesto() int {
	total := h.Historia.SegundosPresupuesto
	for _, d := range h.Descendientes {
		total += d.SegundosPresupuesto
	}
	return total
}

// Tiempo estimado agregando todo el árbol de historias.
func (h *HistoriaAgregado) SegundosEstimado() int {
	total := h.Tareas.SegundosEstimado()
	for _, d := range h.Descendientes {
		total += d.SegundosEstimado
	}
	return total
}

// Tiempo utilizado agregando todo el árbol de historias.
func (h *HistoriaAgregado) SegundosUtilizado() int {
	total := h.Tareas.SegundosUtilizado()
	for _, d := range h.Descendientes {
		total += d.SegundosUtilizado
	}
	return total
}

// Progreso agregando la ponderación y avance de todo el árbol de historias.
func (h *HistoriaAgregado) AvancePorcentual() float64 {
	avance := h.Tareas.AvancePonderado()
	valor := h.Tareas.ValorPonderado()
	for _, d := range h.Descendientes {
		avance += d.AvancePonderado()
		valor += d.ValorPonderado()
	}
	if valor == 0 {
		return 0
	}
	return math.Round(float64(avance)/float64(valor)*10*100) / 10
}
