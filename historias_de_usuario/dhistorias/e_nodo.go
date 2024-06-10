package dhistorias

import "monorepo/historias_de_usuario/ust"

type Nodo struct {
	Nodo     ust.Nodo
	Persona  *ust.Persona
	Historia *ust.NodoHistoria
	Tarea    *ust.Tarea
}

type NodoConHijos struct {
	Nodo     ust.Nodo
	Persona  *ust.Persona
	Historia *ust.NodoHistoria
	Tarea    *ust.Tarea

	Padre *Nodo
	Hijos []Nodo
}
