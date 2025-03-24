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
	Relacionadas  []ust.NodoHistoria

	avance *avanceEscalado
}

// ================================================================ //

// Presupuesto agregando todo el árbol de historias.
func (h *HistoriaAgregado) SegundosPresupuesto() int {
	total := h.Historia.SegundosPresupuesto
	for _, d := range h.Descendientes {
		total += d.SegundosPresupuestoMust()
	}
	return total
}

// Tiempo estimado agregando todo el árbol de historias.
func (h *HistoriaAgregado) SegundosEstimado() int {
	total := h.Tareas.SegundosEstimado()
	for _, d := range h.Descendientes {
		total += d.SegundosEstimadoMust()
	}
	return total
}

// Tiempo utilizado agregando todo el árbol de historias.
func (h *HistoriaAgregado) SegundosUtilizado() int {
	total := h.Tareas.SegundosUtilizado()
	for _, d := range h.Descendientes {
		total += d.SegundosUtilizadoMust()
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

// Porcentaje utilizado del presupuesto.
func (h *HistoriaAgregado) DesviacionPresupuestal() float64 {
	if h.SegundosPresupuesto() == 0 {
		return 0
	}
	return math.Round(float64(h.SegundosUtilizado())/float64(h.SegundosPresupuesto())*100*10) / 10
}

// Tiempo que debería haberse gastado del persupuesto según el avance obtenido.
func (h *HistoriaAgregado) SegundosExpectativaAvancePresupuesto() int {
	if h.AvancePorcentual() == 0 {
		return 0
	}
	return int(math.Round(float64(h.SegundosPresupuesto()) * h.AvancePorcentual() / 100))
}

// ================================================================ //

// Para valor en gráfico de barras.
func (h *HistoriaAgregado) HorasPresupuesto() float64 {
	return math.Round(float64(h.SegundosPresupuesto())/3600*100) / 100
}
func (h *HistoriaAgregado) HorasEstimado() float64 {
	return math.Round(float64(h.SegundosEstimado())/3600*100) / 100
}
func (h *HistoriaAgregado) HorasUtilizado() float64 {
	return math.Round(float64(h.SegundosUtilizado())/3600*100) / 100
}
func (h *HistoriaAgregado) HorasExpectativaAvancePresupuesto() float64 {
	return math.Round(float64(h.SegundosExpectativaAvancePresupuesto())/3600*100) / 100
}

// ================================================================ //

type HistoriaAgregadoList []HistoriaAgregado

func (h HistoriaAgregadoList) SegundosPresupuesto() int {
	total := 0
	for _, d := range h {
		total += d.SegundosPresupuesto()
	}
	return total
}
func (h HistoriaAgregadoList) SegundosEstimado() int {
	total := 0
	for _, d := range h {
		total += d.SegundosEstimado()
	}
	return total
}
func (h HistoriaAgregadoList) SegundosUtilizado() int {
	total := 0
	for _, d := range h {
		total += d.SegundosUtilizado()
	}
	return total
}
