package arbol

type Repo interface {
	GetProyecto(proyectoID int) (*Proyecto, error)
	GetPersona(personaID int) (*Persona, error)
	GetHistoria(historiaID int) (*HistoriaDeUsuario, error)

	ReordenarNodo(nodoID int, newPosicion int) error
}
