package dhistorias

import "monorepo/ust"

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

func (h *HistoriaRecursiva) TiempoReal() int {
	total := h.Segundos
	for _, h := range h.Descendientes {
		total += h.TiempoReal()
	}
	return total
}
