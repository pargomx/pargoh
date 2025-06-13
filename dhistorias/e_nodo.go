package dhistorias

import "monorepo/ust"

type Nodo struct {
	Nodo     ust.Nodo
	Persona  *ust.Persona
	Historia *ust.NodoHistoria
	Tarea    *ust.Tarea
}
