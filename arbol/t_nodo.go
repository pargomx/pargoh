package arbol

// Nodo corresponde a un elemento de la tabla 'nodos'.
type Nodo struct {
	NodoID         int    // `nodos.nodo_id`
	PadreID        int    // `nodos.padre_id`
	Posicion       int    // `nodos.posicion`  Posición consecutiva con respecto a sus nodos hermanos
	Tipo           string // `nodos.tipo`  Tipo de nodo
	Título         string // `nodos.titulo`  Nombre o título de la entidad.
	Descripcion    string // `nodos.descripcion`
	Objetivo       string // `nodos.objetivo`
	Notas          string // `nodos.notas`
	Color          string // `nodos.color`  Proyectos
	Imagen         string // `nodos.imagen`  Proyectos
	Prioridad      int    // `nodos.prioridad`  Historias
	Completada     bool   // `nodos.completada`  Historias
	PresupSegundos int    // `nodos.presup_seg`  Presupuesto en tiempo
	PresupCentavos int    // `nodos.presup_cent`  Presupuesto en dinero
}
