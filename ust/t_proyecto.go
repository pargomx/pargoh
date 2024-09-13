package ust

import "errors"

// Proyecto corresponde a un elemento de la tabla 'proyectos'.
type Proyecto struct {
	ProyectoID    string // `proyectos.proyecto_id`
	Titulo        string // `proyectos.titulo`
	Descripcion   string // `proyectos.descripcion`
	Imagen        string // `proyectos.imagen`
	TiempoGestion int    // `proyectos.tiempo_gestion`  Número de segundos que se ha trabajado en la gestión del proyecto dentro de Pargo
}

var (
	ErrProyectoNotFound      error = errors.New("el proyecto no se encuentra")
	ErrProyectoAlreadyExists error = errors.New("el proyecto ya existe")
)

func (pro *Proyecto) Validar() error {

	return nil
}
