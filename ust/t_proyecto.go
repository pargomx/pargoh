package ust

import "errors"

// Proyecto corresponde a un elemento de la tabla 'proyectos'.
type Proyecto struct {
	ProyectoID  string // `proyectos.proyecto_id`
	Titulo      string // `proyectos.titulo`
	Imagen      string // `proyectos.imagen`
	Descripcion string // `proyectos.descripcion`
}

var (
	ErrProyectoNotFound      error = errors.New("el proyecto no se encuentra")
	ErrProyectoAlreadyExists error = errors.New("el proyecto ya existe")
)

func (pro *Proyecto) Validar() error {

	return nil
}
