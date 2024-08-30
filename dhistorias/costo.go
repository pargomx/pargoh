package dhistorias

import "monorepo/ust"

// ================================================================ //
// ========== TIEMPO ============================================== //

func (h *Historia) TiempoEstimado() int {
	total := 0
	for _, t := range h.Tareas {
		total += t.TiempoEstimado
	}
	return total
}

func (h *Historia) TiempoEstimadoString() string {
	return ust.MinutosToString(h.TiempoEstimado())
}

func (h *Historia) TiempoReal() int {
	total := 0
	for _, t := range h.Tareas {
		total += t.TiempoReal
	}
	return total
}

func (h *Historia) TiempoRealString() string {
	return ust.SegundosToString(h.TiempoReal())
}
