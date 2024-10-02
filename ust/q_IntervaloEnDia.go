package ust

// IntervaloEnDia corresponde a una consulta de solo lectura.
type IntervaloEnDia struct {
	//  `his.proyecto_id`
	ProyectoID string
	//  `his.persona_id`
	PersonaID int
	//  `his.historia_id`
	HistoriaID int
	//  `interv.tarea_id`
	TareaID int
	//  `interv.inicio`
	Inicio string
	//  `interv.fin`
	Fin string
	//  `date(interv.inicio,'-5 hours')`
	Fecha string
	//  `unixepoch(coalesce(nullif(interv.fin,''),datetime('now','-6 hours'))) - unixepoch(interv.inicio)`
	Segundos int
}
