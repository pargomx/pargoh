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
