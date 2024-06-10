package ust

// NodoHistoria corresponde a una consulta de solo lectura.
type NodoHistoria struct {
	//  `his.historia_id`
	HistoriaID int
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
	//  `(SELECT COUNT(nodo_id) FROM nodos WHERE padre_id = his.historia_id)`
	NumHistorias int
	//  `(SELECT COUNT(tarea_id) FROM tareas WHERE historia_id = his.historia_id)`
	NumTareas int
}
