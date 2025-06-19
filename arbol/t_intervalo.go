package arbol

// Intervalo corresponde a un elemento de la tabla 'intervalos'.
type Intervalo struct {
	NodoID int    // `intervalos.nodo_id`  ID del nodo
	TsIni  string // `intervalos.ts_ini`  Inicio del intervalo en hora local
	TsFin  string // `intervalos.ts_fin`  Fin del intervalo en hora local
}
