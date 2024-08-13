package ust

// NodoPersona corresponde a una consulta de solo lectura.
type NodoPersona struct {
	//  `per.persona_id`
	PersonaID int
	//  `per.proyecto_id`
	ProyectoID string
	//  `per.nombre`
	Nombre string
	//  `per.descripcion`
	Descripcion string
	//  `nod.padre_id`
	PadreID int
	//  `nod.padre_tbl`
	PadreTbl string
	//  `nod.nivel`
	Nivel int
	//  `nod.posicion`
	Posicion int
}
