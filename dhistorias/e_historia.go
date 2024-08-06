package dhistorias

import "monorepo/ust"

type HistoriaConNietos struct {
	Persona   ust.Persona        // Siempre hay persona
	Ancestros []ust.NodoHistoria // Lista de ancestros desde el más grande al más pequeño
	Abuelo    *ust.NodoHistoria  // No siempre hay padre
	Padres    []HistoriaConHijos // Puede haber o no hijos
	Tareas    []ust.Tarea        // Puede haber o no tareas
}

type HistoriaConHijos struct {
	Padre  ust.NodoHistoria
	Hijos  []ust.NodoHistoria
	Tareas []ust.Tarea
}

type Arbol struct {
	Persona   ust.NodoPersona
	Historias []HistoriaRecursiva
}
type HistoriaRecursiva struct {
	Padre  ust.NodoHistoria
	Hijos  []HistoriaRecursiva
	Tareas []ust.Tarea
}
