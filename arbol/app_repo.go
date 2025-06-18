package arbol

type ReadRepo interface {
	GetProyecto(proyectoID int) (*Proyecto, error)
	GetPersona(personaID int) (*Persona, error)
	GetHistoria(historiaID int) (*HistoriaDeUsuario, error)

	GetNodo(NodoID int) (*Nodo, error)
}

type Repo interface {
	ReadRepo

	InsertNodo(nod Nodo) error
	UpdateNodo(NodoID int, nod Nodo) error
	DeleteNodo(NodoID int) error
	ReordenarNodo(nod Nodo, newPosicion int) error
}
