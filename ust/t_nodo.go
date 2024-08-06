package ust

import "errors"

// Nodo corresponde a un elemento de la tabla 'nodos'.
type Nodo struct {
	NodoID   int    // `nodos.nodo_id`  ID del nodo
	NodoTbl  string // `nodos.nodo_tbl`  Tabla de donde obtener el nodo
	PadreID  int    // `nodos.padre_id`  ID del nodo padre
	PadreTbl string // `nodos.padre_tbl`  Tabla de donde obtener el padre
	Nivel    int    // `nodos.nivel`  Nivel de profundidad en el árbol
	Posicion int    // `nodos.posicion`  Posición consecutiva con respecto a sus nodos hermanos
}

var (
	ErrNodoNotFound      error = errors.New("el nodo no se encuentra")
	ErrNodoAlreadyExists error = errors.New("el nodo ya existe")
	ErrNodoIDInvalido    error = errors.New("nodo_id debe ser mayor a 0")
	ErrNivelInvalido     error = errors.New("nivel debe ser mayor a 0")
)

func (nod *Nodo) Validar() error {
	if nod.NodoID < 1 {
		return ErrNodoIDInvalido
	}
	if nod.Nivel < 0 {
		return ErrNivelInvalido
	}
	return nil
}

// ================================================================ //
// ================================================================ //

const (
	TipoNodoPersona  = "per"
	TipoNodoHistoria = "his"
	TipoNodoTarea    = "tar"
	RootNodoID       = 0
)

func (n Nodo) EsPersona() bool {
	return n.NodoTbl == TipoNodoPersona
}

func (n Nodo) EsHistoria() bool {
	return n.NodoTbl == TipoNodoHistoria
}

func (n Nodo) EsTarea() bool {
	return n.NodoTbl == TipoNodoTarea
}
