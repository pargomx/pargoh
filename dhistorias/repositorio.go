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
	ListHistorias() ([]ust.Historia, error)
	ListNodoHistorias(PadreID int) ([]ust.NodoHistoria, error)

	// Tareas
	InsertTarea(tar ust.Tarea) error
	UpdateTarea(tar ust.Tarea) error
	DeleteTarea(tareaID int) error
	DeleteAllTareas(HistoriaID int) error
	GetTarea(tareaID int) (*ust.Tarea, error)
	ListTareas() ([]ust.Tarea, error)
	ListTareasByHistoriaID(historiaID int) ([]ust.Tarea, error)

	// Intervalos
	InsertIntervalo(interv ust.Intervalo) error
	UpdateIntervalo(interv ust.Intervalo) error
	DeleteIntervalo(TareaID int, Inicio string) error
	ListIntervalosByTareaID(TareaID int) ([]ust.Intervalo, error)

	// Viajes
	InsertTramo(tra ust.Tramo) error
	UpdateTramo(tra ust.Tramo) error
	ExisteTramo(HistoriaID int, Posicion int) error
	DeleteTramo(HistoriaID int, Posicion int) error
	DeleteAllTramos(HistoriaID int) error
	GetTramo(HistoriaID int, Posicion int) (*ust.Tramo, error)
	ListTramosByHistoriaID(HistoriaID int) ([]ust.Tramo, error)
	ReordenarTramo(HistoriaID, oldPos, newPos int) error
	MoverTramo(historiaID int, posicion int, newHistoriaID int) error

	// Reglas
	InsertRegla(reg ust.Regla) error
	UpdateRegla(reg ust.Regla) error
	ExisteRegla(HistoriaID int, Posicion int) error
	DeleteRegla(HistoriaID int, Posicion int) error
	DeleteAllReglas(HistoriaID int) error
	GetRegla(HistoriaID int, Posicion int) (*ust.Regla, error)
	ListReglasByHistoriaID(HistoriaID int) ([]ust.Regla, error)
	ReordenarRegla(HistoriaID, oldPos, newPos int) error
}
