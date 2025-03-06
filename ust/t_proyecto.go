package ust

// Proyecto corresponde a un elemento de la tabla 'proyectos'.
type Proyecto struct {
	ProyectoID    string // `proyectos.proyecto_id`
	Posicion      int    // `proyectos.posicion`  Posición consecutiva con respecto a sus nodos hermanos
	Titulo        string // `proyectos.titulo`
	Color         string // `proyectos.color`
	Imagen        string // `proyectos.imagen`
	Descripcion   string // `proyectos.descripcion`
	FechaRegistro string // `proyectos.fecha_registro`
}
