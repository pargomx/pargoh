package ust

import "github.com/pargomx/gecko/gko"

func (t *Tarea) FactorImportancia() int {
	switch {
	case t.Importancia.EsIdea(), t.Importancia.EsIndefinido():
		return 1
	case t.Importancia.EsMejora():
		return 32
	case t.Importancia.EsNecesaria():
		return 128
	default:
		gko.LogWarnf("Importancia inválida: %v para tarea %v", t.Importancia, t.TareaID)
		return 0
	}
}

const SEGUNDOS_PONDERACION_DEFAULT = 3600 // 1h

// Puntaje basado en la importancia y esfuerzo en tiempo de la tarea.
func (t *Tarea) ValorPonderado() int {
	segundos := t.SegundosEstimado
	if t.Finalizada() ||
		(!t.Finalizada() && t.SegundosUtilizado > t.SegundosEstimado) {
		segundos = t.SegundosUtilizado
	}
	if segundos == 0 {
		segundos = SEGUNDOS_PONDERACION_DEFAULT
	}
	return segundos * t.FactorImportancia()
}

func (t *Tarea) AvancePonderado() int {
	pond := t.ValorPonderado()
	if t.Finalizada() {
		return pond
	}
	if t.SegundosUtilizado > t.SegundosEstimado {
		return pond * 90 / 100 // 90% para tareas no completadas que superaron estimado.
	}
	return t.SegundosUtilizado * t.FactorImportancia()
}
