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

// Relación entre ValorPonderado y AvancePonderado
// obtenido de las tareas en la historia raíz & sus hijos.
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
