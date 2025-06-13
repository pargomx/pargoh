package arbol

// Nodo corresponde a un elemento de la tabla 'nodos'.
type Nodo struct {
	NodoID   int    // `nodos.nodo_id`
	PadreID  int    // `nodos.padre_id`
	Tipo     string // `nodos.tipo`  Tipo de nodo
	Posicion int    // `nodos.posicion`  Posición consecutiva con respecto a sus nodos hermanos

	Título      string // `nodos.titulo`  Nombre o título de la entidad.
	Descripcion string // `nodos.descripcion`
	Objetivo    string // `nodos.objetivo`
	Notas       string // `nodos.notas`

	Color  string // `nodos.color`
	Imagen string // `nodos.imagen`

	Prioridad int // `nodos.prioridad`
	Estatus   int // `nodos.estatus`  Completada, implementada, etc.
	Segundos  int // `nodos.segundos`  Tiempo para guardar presupuestos, estimados, etc.
	Centavos  int // `nodos.centavos`  Centavos para guardar dinero de presupuesto, inversión, etc.
}
