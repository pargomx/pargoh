package ust

// IntervaloEnDia corresponde a una consulta de solo lectura.
type IntervaloEnDia struct {
	//  `his.historia_id`
	HistoriaID int
	//  `interv.tarea_id`
	TareaID int
	//  `date(interv.inicio,'-5 hours')`
	Fecha string
	//  `interv.inicio`
	Inicio string
	//  `interv.fin`
	Fin string
}
