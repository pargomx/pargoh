package dhistorias

import "monorepo/ust"

type HistoriaConNietos struct {
	Persona   ust.Persona        // Siempre hay persona
	Proyecto  ust.Proyecto       // Siempre hay proyecto
	Ancestros []ust.NodoHistoria // Lista de ancestros desde el más grande al más pequeño
	Abuelo    *ust.NodoHistoria  // No siempre hay padre
	Padres    []HistoriaConHijos // Puede haber o no hijos
	Tareas    []ust.Tarea        // Puede haber o no tareas
	Tramos    []ust.Tramo        // Puede haber o no tramos
}

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

type HistoriaConHijos struct {
	Padre  ust.NodoHistoria
	Hijos  []ust.NodoHistoria
	Tareas []ust.Tarea
}

type HistoriaRecursiva struct {
	Historia      ust.NodoHistoria
	Descendientes []HistoriaRecursiva
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
