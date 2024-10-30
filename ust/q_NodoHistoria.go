package ust

// NodoHistoria corresponde a una consulta de solo lectura.
type NodoHistoria struct {
	//  `his.historia_id`
	HistoriaID int
	//  `his.proyecto_id`
	ProyectoID string
	//  `his.persona_id`
	PersonaID int
	//  `his.titulo`
	Titulo string
	//  `his.objetivo`
	Objetivo string
	//  `his.prioridad`
	Prioridad int
	//  `his.completada`
	Completada bool
	//  `nod.padre_id`
	PadreID int
	//  `nod.padre_tbl`
	PadreTbl string
	//  `nod.nivel`
	Nivel int
	//  `nod.posicion`
	Posicion int
	//  `his.segundos_presupuesto`
	SegundosPresupuesto int
	//  `his.descripcion`
	Descripcion string
	//  `(SELECT COUNT(nodo_id) FROM nodos WHERE padre_id = his.historia_id)`
	NumHistorias int
	//  `(SELECT COUNT(tarea_id) FROM tareas WHERE historia_id = his.historia_id)`
	NumTareas int
	//  `(SELECT SUM(segundos_estimado) FROM tareas WHERE historia_id = his.historia_id)`
	SegundosEstimado int
	//  `(SELECT SUM(unixepoch(coalesce(nullif(interv.fin,''),datetime('now','-6 hours'))) - unixepoch(interv.inicio)) FROM intervalos interv JOIN tareas tar ON tar.tarea_id = interv.tarea_id WHERE tar.historia_id = his.historia_id GROUP BY tar.historia_id)`
	SegundosUtilizado int
}
