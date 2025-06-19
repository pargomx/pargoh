package arbol

const NodoProyectosActivos = 2 // NodoID del grupo de Proyectos activos.

func esTipoValido(tipo string) bool {
	switch tipo {
	case "GRP", "PRY", "PER", "HIS", "TEC", "GES", "REG", "TAR", "VIA":
		return true
	default:
		return false
	}
}

func (nod Nodo) EsGrupo() bool {
	return nod.Tipo == "GRP"
}
func (nod Nodo) EsProyecto() bool {
	return nod.Tipo == "PRY"
}
func (nod Nodo) EsPersona() bool {
	return nod.Tipo == "PER"
}
func (nod Nodo) EsHistoriaDeUsuario() bool {
	return nod.Tipo == "HIS"
}
func (nod Nodo) EsHistoriaTecnica() bool {
	return nod.Tipo == "TEC"
}
func (nod Nodo) EsActividadDeGestión() bool {
	return nod.Tipo == "GES"
}
func (nod Nodo) EsRegla() bool {
	return nod.Tipo == "REG"
}
func (nod Nodo) EsTarea() bool {
	return nod.Tipo == "TAR"
}
func (nod Nodo) EsTramo() bool {
	return nod.Tipo == "VIA"
}
