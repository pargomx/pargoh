package ust

// Historia corresponde a un elemento de la tabla 'historias'.
type Historia struct {
	HistoriaID          int    // `historias.historia_id`
	Titulo              string // `historias.titulo`
	Objetivo            string // `historias.objetivo`
	Prioridad           int    // `historias.prioridad`
	Completada          bool   // `historias.completada`
	PersonaID           int    // `historias.persona_id`  Para índice
	ProyectoID          string // `historias.proyecto_id`  Para índice
	SegundosPresupuesto int    // `historias.segundos_presupuesto`  Tiempo estimado en segundos para implementar la historia de usuario en su totalidad
	Descripcion         string // `historias.descripcion`  Descripción  de la historia en infinitivo para que la lea el usuario en la documentación.
	Notas               string // `historias.notas`  Apuntes técnicos sobre la implementación de la historia.
}
