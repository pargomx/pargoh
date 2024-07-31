package dhistorias

import "monorepo/historias_de_usuario/ust"

type Repo interface {

	// Nodos
	InsertNodo(nod ust.Nodo) error
	EliminarNodo(nodoID int) error
	MoverNodo(nodoID int, nuevoPadreID int) error
	ReordenarNodo(nodoID int, oldPosicion int, newPosicion int) error
	GetNodo(nodoID int) (*ust.Nodo, error)
	ListNodosByPadreID(PadreID int) ([]ust.Nodo, error)

	// Personas
	InsertPersona(per ust.Persona) error
	GetPersona(personaID int) (*ust.Persona, error)
	UpdatePersona(per ust.Persona) error
	DeletePersona(personaID int) error
	ListNodosPersonas() ([]ust.NodoPersona, error)

	// Historias
	InsertHistoria(his ust.Historia) error
	UpdateHistoria(ust.Historia) error
	DeleteHistoria(historiaID int) error
	GetHistoria(historiaID int) (*ust.Historia, error)
	GetNodoHistoria(nodoID int) (*ust.NodoHistoria, error)
	ListNodoHistoriasByPadreID(PadreID int) ([]ust.NodoHistoria, error)

	// Tareas
	InsertTarea(tar ust.Tarea) error
	UpdateTarea(tar ust.Tarea) error
	GetTarea(tareaID int) (*ust.Tarea, error)
	ListTareasByHistoriaID(historiaID int) ([]ust.Tarea, error)

	// Intervalos
	InsertIntervalo(interv ust.Intervalo) error
	UpdateIntervalo(interv ust.Intervalo) error
	ListIntervalosByTareaID(TareaID int) ([]ust.Intervalo, error)
}
