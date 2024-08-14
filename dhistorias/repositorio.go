package dhistorias

import (
	"monorepo/ust"
)

type Repo interface {

	// Proyectos
	InsertProyecto(pro ust.Proyecto) error
	GetProyecto(ProyectoID string) (*ust.Proyecto, error)
	UpdateProyecto(pro ust.Proyecto) error
	ExisteProyecto(ProyectoID string) error
	DeleteProyecto(ProyectoID string) error
	ListProyectos() ([]ust.Proyecto, error)

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
	ListNodosPersonas(ProyectoID string) ([]ust.NodoPersona, error)

	// Historias
	ExisteHistoria(HistoriaID int) error
	InsertHistoria(his ust.Historia) error
	UpdateHistoria(ust.Historia) error
	DeleteHistoria(historiaID int) error
	GetHistoria(historiaID int) (*ust.Historia, error)
	GetNodoHistoria(nodoID int) (*ust.NodoHistoria, error)
	ListNodoHistorias(PadreID int) ([]ust.NodoHistoria, error)

	// Tareas
	InsertTarea(tar ust.Tarea) error
	UpdateTarea(tar ust.Tarea) error
	GetTarea(tareaID int) (*ust.Tarea, error)
	ListTareasByHistoriaID(historiaID int) ([]ust.Tarea, error)

	// Intervalos
	InsertIntervalo(interv ust.Intervalo) error
	UpdateIntervalo(interv ust.Intervalo) error
	ListIntervalosByTareaID(TareaID int) ([]ust.Intervalo, error)

	// Viajes
	InsertTramo(tra ust.Tramo) error
	UpdateTramo(tra ust.Tramo) error
	ExisteTramo(HistoriaID int, Posicion int) error
	DeleteTramo(HistoriaID int, Posicion int) error
	GetTramo(HistoriaID int, Posicion int) (*ust.Tramo, error)
	ListTramosByHistoriaID(HistoriaID int) ([]ust.Tramo, error)
}
