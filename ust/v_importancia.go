package ust

import "github.com/pargomx/gecko/gko"

// Importancia de una tarea.

func NuevaImportancia(valor int) (int, error) {
	if valor < 0 {
		return 0, gko.ErrDatoInvalido().Msg("Importancia debe ser mayor o igual a 0")
	}
	return 0, nil
}

const (
	TareaIdea      = 0
	TareaMejora    = 1
	TareaNecesaria = 2
)

var Importancias = []string{"Idea", "Mejora", "Necesaria"} // 0, 1, 2

func (t *Tarea) ImportanciaString() string {
	switch t.Importancia {
	case TareaIdea:
		return "Idea"
	case TareaMejora:
		return "Mejora"
	case TareaNecesaria:
		return "Necesaria"
	default:
		gko.LogWarnf("Importancia inválida: %v para tarea %v", t.Importancia, t.TareaID)
		return "???"
	}
}

func (t *Tarea) EsIdea() bool {
	return t.Importancia == TareaIdea
}
func (t *Tarea) EsMejora() bool {
	return t.Importancia == TareaMejora
}
func (t *Tarea) EsNecesaria() bool {
	return t.Importancia == TareaNecesaria
}

func (t *Tarea) PonderacionImportancia() int {
	switch t.Importancia {
	case TareaIdea:
		return 1
	case TareaMejora:
		return 32
	case TareaNecesaria:
		return 128
	default:
		gko.LogWarnf("Importancia inválida: %v para tarea %v", t.Importancia, t.TareaID)
		return 0
	}
}

// TODO: todos los campos deberían usar segundos para ser consistentes.
// TODO: usar campo STRING para importancia, de todas formas no se usan los números de la DB en la lógica y así es expandible a otras ponderaciones.

const SEGUNDOS_CALCULO_DEFAULT = 3600 // 1h

func (t *Tarea) ValorPonderado() int {
	segundos := t.SegundosEstimado
	if t.Finalizada() ||
		(!t.Finalizada() && t.SegundosReal >= t.SegundosEstimado) {
		segundos = t.SegundosReal
	}
	if segundos == 0 {
		segundos = SEGUNDOS_CALCULO_DEFAULT
	}
	return segundos * t.PonderacionImportancia()
}

func (t *Tarea) AvancePonderado() int {
	pond := t.ValorPonderado()
	if t.Finalizada() {
		return pond
	}
	if t.SegundosReal >= t.SegundosEstimado*60 {
		return pond * 90 / 100 // 90%
	}
	return t.PonderacionImportancia() * t.SegundosReal
}
