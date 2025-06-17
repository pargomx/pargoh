package arbol

import (
	"math"
	"monorepo/ust"
)

// Para gráfico SVG de avance.
type avanceEscalado struct {
	Presupuesto int
	Estimado    int
	Utilizado   int
	Expectativa int
	Separadores []separadorHora
}

type separadorHora struct {
	Hora          int  // número de hora
	Posicion      int  // a escala
	EsPresupuesto bool // si este separador es el presupuesto
	EsUltimo      bool // si es el último separador
}

// Medidas de avance, presupuesto, estimado y utilizado relativas a mil.
func (h *HistoriaDeUsuario) AvanceRelativoMil() avanceEscalado {
	return avanceEscalado{}
}

func (h *HistoriaDeUsuario) SegundosMaxPresupuestoEstimadoUtilizado() int {
	return 10
}

// Porcentaje utilizado del presupuesto.
func (h *HistoriaDeUsuario) DesviacionPresupuestal() float64 {
	return 0
}

// Tiempo que debería haberse gastado del persupuesto según el avance obtenido.
func (h *HistoriaDeUsuario) SegundosExpectativaAvancePresupuesto() int {
	return 0
}

// ================================================================ //

// Para valor en gráfico de barras.
func (h *HistoriaDeUsuario) HorasPresupuesto() float64 {
	return 0
}
func (h *HistoriaDeUsuario) HorasEstimado() float64 {
	return 0
}
func (h *HistoriaDeUsuario) HorasUtilizado() float64 {
	return 0
}
func (h *HistoriaDeUsuario) HorasExpectativaAvancePresupuesto() float64 {
	return 0
}

// ================================================================ //
// ================================================================ //

type TareasList []Tarea

// Suma del tiempo estimado para las tareas.
func (tareas TareasList) SegundosEstimado() (total int) {
	for _, t := range tareas {
		total += t.SegundosEstimado
	}
	return total
}

// Suma del tiempo utilizado para las tareas.
func (tareas TareasList) SegundosUtilizado() (total int) {
	for _, t := range tareas {
		total += t.SegundosUtilizado
	}
	return total
}

// Suma del valor ponderado para las tareas.
func (tareas TareasList) ValorPonderado() (total int) {
	return total
}

// Suma del avance ponderado para las tareas.
func (tareas TareasList) AvancePonderado() (total int) {
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

func (tareas TareasList) ValorPorcentual(x int) float64 {
	return 0
}

func (tar Tarea) Finalizada() bool {
	return tar.Estatus > 1
}

func (tar Tarea) Importancia() ust.ImportanciaTarea {
	switch tar.Prioridad {
	default:
		return ust.ImportanciaTareaIndefinido
	case 1:
		return ust.ImportanciaTareaIdea
	case 2:
		return ust.ImportanciaTareaMejora
	case 3:
		return ust.ImportanciaTareaNecesaria
	}
}

func (tar Tarea) Tipo() ust.TipoTarea {
	return ust.TipoTareaIndefinido
}

const (
	EstatusTareaNoIniciada = 0
	EstatusTareaEnCurso    = 1
	EstatusTareaEnPausa    = 2
	EstatusTareaFinalizada = 3
)

func (tar Tarea) NoIniciada() bool {
	return tar.Estatus == EstatusTareaNoIniciada
}
func (tar Tarea) EnCurso() bool {
	return tar.Estatus == EstatusTareaEnCurso
}
func (tar Tarea) EnPausa() bool {
	return tar.Estatus == EstatusTareaEnPausa
}

// func (tar Tarea) Finalizada() bool {
// 	return tar.Estatus == EstatusTareaFinalizada
// }

// func (tar IntervaloReciente) NoIniciada() bool {
// 	return tar.Estatus == EstatusTareaNoIniciada
// }
// func (tar IntervaloReciente) EnCurso() bool {
// 	return tar.Estatus == EstatusTareaEnCurso
// }
// func (tar IntervaloReciente) EnPausa() bool {
// 	return tar.Estatus == EstatusTareaEnPausa
// }
// func (tar IntervaloReciente) Finalizada() bool {
// 	return tar.Estatus == EstatusTareaFinalizada
// }

func (t *Tarea) FactorImportancia() int { return 1 }
func (t *Tarea) ValorPonderado() int    { return 1 }
func (t *Tarea) AvancePonderado() int   { return 1 }
