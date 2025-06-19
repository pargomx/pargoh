package arbol

// Latido corresponde a un elemento de la tabla 'latidos'.
type Latido struct {
	TsLatido string // `latidos.ts_latido`
	NodoID   int    // `latidos.nodo_id`
	Segundos int    // `latidos.segundos`
}
