package ust

// Latido corresponde a un elemento de la tabla 'latidos'.
type Latido struct {
	Timestamp string // `latidos.timestamp`
	PersonaID int    // `latidos.persona_id`
	Segundos  int    // `latidos.segundos`
}

// Cantidad de minutos transcurridos desde las 6am del día de trabajo.
func (i *Latido) MinutosSince6am() int {
	return MinutosSince6am(i.Timestamp)
}
