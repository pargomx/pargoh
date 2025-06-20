package arbol

type ReadRepo interface {
	GetProyecto(proyectoID int) (*Proyecto, error)
	GetPersona(personaID int) (*Persona, error)
	GetHistoria(historiaID int) (*HistoriaDeUsuario, error)

	GetNodo(NodoID int) (*Nodo, error)
	ListNodosByPadreID(PadreID int) ([]Nodo, error)
	ListNodosByPadreIDTipo(PadreID int, Tipo string) ([]Nodo, error)

	ExisteNodo(NodoID int) error

	ListLatidos(desde, hasta string) ([]Latido, error)
	ListIntervalosByNodoID(NodoID int) ([]Intervalo, error)
	ExisteIntervalo(NodoID int, TsIni string) error
	GetIntervalo(NodoID int, TsIni string) (*Intervalo, error)
}

type Repo interface {
	ReadRepo

	InsertNodo(nod Nodo) error
	UpdateNodo(NodoID int, nod Nodo) error
	DeleteNodo(NodoID int) error
	DeleteHijos(NodoID int) error
	ReordenarNodo(nod Nodo, newPosicion int) error
	MoverNodo(nod Nodo, newPadreID int) error

	InsertLatido(lat Latido) error
	InsertIntervalo(itv Intervalo) error
	UpdateIntervalo(NodoID int, TsIni string, itv Intervalo) error
	DeleteIntervalo(NodoID int, TsIni string) error

	InsertReferencia(ref Referencia) error
	ExisteReferencia(NodoID int, RefNodoID int) error
	DeleteReferencia(NodoID int, RefNodoID int) error
}

type timeTrackerRepo interface {
	GetNodo(NodoID int) (*Nodo, error)
	UpdateNodo(NodoID int, nod Nodo) error
	InsertLatido(lat Latido) error
}
