package ust

// IntervaloReciente corresponde a una consulta de solo lectura.
type IntervaloReciente struct {
	//  `tar.historia_id`
	HistoriaID int
	//  `interv.tarea_id`
	TareaID int
	//  `interv.inicio`
	Inicio string
	//  `interv.fin`
	Fin string
	//  `tar.tipo`
	Tipo TipoTarea
	//  `tar.descripcion`
	Descripcion string
	//  `tar.impedimentos`
	Impedimentos string
	//  `tar.tiempo_estimado`
	TiempoEstimado int
	//  `tar.tiempo_real`
	TiempoReal int
	//  `tar.estatus`
	Estatus int
	//  `his.titulo`
	Titulo string
	//  `his.objetivo`
	Objetivo string
	//  `his.completada`
	Completada bool
	//  `his.prioridad`
	Prioridad int
}
