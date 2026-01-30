package arbol

const (
	TipoGrupo    = "GRP"
	TipoProyecto = "PRY"
	TipoPersona  = "PER"
	TipoHistoria = "HIS"
	TipoTecnica  = "TEC"
	TipoGestion  = "GES"
	TipoRegla    = "REG"
	TipoTarea    = "TAR"
	TipoViaje    = "VIA"
	TipoRaiz     = "ROOT"
)

type Raiz struct {
	Grupos        []Grupo
	Proyectos     []Proyecto
	TareasEnCurso []Tarea
}

// ================================================================ //

// Grupo de proyectos para organizarlos en carpetas.
//
// Por default se generan dos grupos:
//   - "Activos" para los proyectos que se muestran en la página de inicio.
//   - "Archivados" para los que se ocultan de la página de inicio.
//
// El padre siempre es el nodo ROOT.
type Grupo struct {
	GrupoID  int
	Posicion int

	Nombre string

	Grupos    []Grupo
	Proyectos []Proyecto
}

func (nod Nodo) ToGrupo() Grupo {
	return Grupo{
		GrupoID:  nod.NodoID,
		Posicion: nod.Posicion,

		Nombre: nod.Titulo,
	}
}

// ================================================================ //

// Proyecto representa el esfuerzo de desarrollar una o varias apps.
//
// El padre siempre es un Grupo.
type Proyecto struct {
	ProyectoID int // nodo "PRY"
	GrupoID    int
	Posicion   int

	Titulo      string
	Descripcion string
	// Objetivo string
	// Notas    string

	Color  string
	Imagen string

	// Prioridad int
	// Estatus   int
	// Segundos  int
	// Centavos  int

	Personas    []Persona
	HisUsuario  []HistoriaDeUsuario
	HisTecnicas []HistoriaTecnica
	HisGestion  []ActividadDeGestión
}

func (nod Nodo) ToProyecto() Proyecto {
	return Proyecto{
		ProyectoID: nod.NodoID,
		GrupoID:    nod.PadreID,
		Posicion:   nod.Posicion,

		Titulo:      nod.Titulo,
		Descripcion: nod.Descripcion,
		Color:       nod.Color,
		Imagen:      nod.Imagen,
	}
}

// ================================================================ //

// Persona de las historias de usuario descendientes. Para hacer mapa de
// empatía y escribir historias desde la perspectiva del usuario.
//
// El padre puede ser un proyecto o historia.
type Persona struct {
	PersonaID int // nodo "PER"
	PadreID   int // Proyecto o Historia.
	Posicion  int

	Nombre      string
	Descripcion string
	// Objetivo string
	// Notas    string

	// Color  string
	// Imagen string

	// Prioridad int
	// Estatus   int
	// Segundos  int
	// Centavos  int

	Personas []Persona

	Historias   []HistoriaDeUsuario
	HisTecnicas []HistoriaTecnica
	HisGestion  []ActividadDeGestión

	// LEGACY
	ProyectoID      string
	SegundosGestion int
	HorasGestion    int
}

func (nod Nodo) ToPersona() Persona {
	return Persona{
		PersonaID: nod.NodoID,
		PadreID:   nod.PadreID,
		Posicion:  nod.Posicion,

		Nombre:      nod.Titulo,
		Descripcion: nod.Descripcion,
	}
}

// ================================================================ //

// Historia de usuario que representan funcionalidad que aporta valor a quien
// utiliza la aplicación o el software. Se pueden descomponer en historias más
// pequeñas hasta hacerlas unidades discretas de trabajo.
//
// El padre puede ser un proyecto, persona o historia.
type HistoriaDeUsuario struct {
	HistoriaID int
	PadreID    int // Proyecto, Persona, Historia.
	Posicion   int

	Titulo      string
	Descripcion string
	Objetivo    string
	Notas       string

	Imagen string

	Prioridad           int
	Completada          bool
	SegundosPresupuesto int
	Centavos            int

	Proyecto  Proyecto // Siempre debe tener este ancestro
	Ancestros []Nodo   // Personas u historias

	Personas    []Persona
	HisUsuario  []HistoriaDeUsuario
	HisTecnicas []HistoriaTecnica
	HisGestion  []ActividadDeGestión

	Reglas []Regla
	Tareas TareasList
	Tramos []Tramo

	Relacionadas []HistoriaDeUsuario

	// Legacy
	Persona                 Persona
	PersonaID               int
	ProyectoID              string
	SegundosPresupuestoMust int
	SegundosEstimadoMust    int
	SegundosEstimado        int
	AvancePorcentual        float32
	SegundosUtilizado       int
	Utilizado               int
}

func (his HistoriaDeUsuario) Descendientes() []HistoriaDeUsuario {
	return his.HisUsuario
}

func (nod Nodo) ToHistoriaDeUsuario() HistoriaDeUsuario {
	return HistoriaDeUsuario{
		HistoriaID: nod.NodoID,
		PadreID:    nod.PadreID,
		Posicion:   nod.Posicion,

		Titulo:      nod.Titulo,
		Descripcion: nod.Descripcion,
		Objetivo:    nod.Objetivo,
		Notas:       nod.Notas,

		Imagen: nod.Imagen,

		Prioridad:           nod.Prioridad,
		Completada:          nod.Estatus > 0,
		SegundosPresupuesto: nod.Segundos,
		Centavos:            nod.Centavos,
	}
}

// ================================================================ //

// Historia técnica que representa una mejora o trabajo necesario relacionado
// con el funcionamiento interno de la aplicación: optimizar, actualizar,
// refactorizar, configurar, soportar más plataformas, etc.
//
// Cuestiones técnicas o de configuración que son parte del proyecto en general
// pero no aportan valor funcional al software, sino que mejoran la seguridad,
// eficiencia o mantenibilidad del sistema.
//
// El padre puede ser un proyecto, persona o historia.
type HistoriaTecnica struct {
	HistoriaID int
	PadreID    int // Proyecto, Persona, Historia.
	Posicion   int

	Titulo      string
	Descripcion string
	Objetivo    string
	Notas       string

	Imagen string

	Prioridad           int
	Completada          bool
	SegundosPresupuesto int
	Centavos            int

	// Materializa
	SegundosDocumentacion int
	SegundosUtilizado     int

	Personas    []Persona
	HisUsuario  []HistoriaDeUsuario
	HisTecnicas []HistoriaTecnica
	HisGestion  []ActividadDeGestión

	Reglas []Regla
	Tareas []Tarea
	Tramos []Tramo
}

func (nod Nodo) ToHistoriaTecnica() HistoriaTecnica {
	return HistoriaTecnica{
		HistoriaID: nod.NodoID,
		PadreID:    nod.PadreID,
		Posicion:   nod.Posicion,

		Titulo:      nod.Titulo,
		Descripcion: nod.Descripcion,
		Objetivo:    nod.Objetivo,
		Notas:       nod.Notas,

		Imagen: nod.Imagen,

		Prioridad:           nod.Prioridad,
		Completada:          nod.Estatus > 0,
		SegundosPresupuesto: nod.Segundos,
		Centavos:            nod.Centavos,
	}
}

// ================================================================ //

// Actividades de gestión parte del ciclo de vida de desarrollo y mantenimiento
// de la aplicación:  documentación, juntas, soporte técnico, proceso de venta, etc.
//
// El padre puede ser un proyecto, persona o historia.
type ActividadDeGestión struct {
	HistoriaID int
	PadreID    int // Proyecto, Persona, Historia.
	Posicion   int

	Titulo      string
	Descripcion string
	Objetivo    string
	Notas       string

	Imagen string

	Prioridad           int
	Completada          bool
	SegundosPresupuesto int
	Centavos            int

	// Materializa
	SegundosDocumentacion int
	SegundosUtilizado     int

	Personas    []Persona
	HisUsuario  []HistoriaDeUsuario
	HisTecnicas []HistoriaTecnica
	HisGestion  []ActividadDeGestión

	Reglas []Regla
	Tareas []Tarea
	Tramos []Tramo
}

func (nod Nodo) ToActividadDeGestión() ActividadDeGestión {
	return ActividadDeGestión{
		HistoriaID: nod.NodoID,
		PadreID:    nod.PadreID,
		Posicion:   nod.Posicion,

		Titulo:      nod.Titulo,
		Descripcion: nod.Descripcion,
		Objetivo:    nod.Objetivo,
		Notas:       nod.Notas,

		Imagen: nod.Imagen,

		Prioridad:           nod.Prioridad,
		Completada:          nod.Estatus > 0,
		SegundosPresupuesto: nod.Segundos,
		Centavos:            nod.Centavos,
	}
}

// ================================================================ //

// Regla de negocio.
//
// El padre siempre es una historia.
type Regla struct {
	ReglaID    int // nod.NodoID
	HistoriaID int // nod.PadreID
	Posicion   int

	Texto string // nod.Titulo
	// Descripcion string
	// Objetivo    string
	// Notas       string

	// Color  string
	// Imagen string

	// Prioridad int

	Estatus int // 0:gris, 1:naranja, 2:verde

	// Segundos int
	// Centavos int

	// Legacy
	Implementada bool // Estatus
	Probada      bool // Estatus
}

func (nod Nodo) ToRegla() Regla {
	return Regla{
		ReglaID:    nod.NodoID,
		HistoriaID: nod.PadreID,
		Posicion:   nod.Posicion,

		Texto:   nod.Titulo,
		Estatus: nod.Estatus,
	}
}

// ================================================================ //

// Tarea de desarrollo
//
// El padre siempre es una historia.
type Tarea struct {
	TareaID    int // nod.NodoID
	HistoriaID int // nod.PadreID
	Posicion   int

	Descripcion  string // nod.Titulo
	Impedimentos string // nod.Descripcion
	// Objetivo    string // contiene backup TipoTarea
	// Notas       string // contiene backup SegundosEstimado

	// Color  string
	// Imagen string

	Prioridad        int // Importancia 0:indefinido, 1:idea, 2:mejora, 3:necesaria
	Estatus          int
	SegundosEstimado int // nod.Segundos
	// Centavos         int

	SegundosUtilizado int // Deprecated: Legacy

	Intervalos []Intervalo
}

func (nod Nodo) ToTarea() Tarea {
	return Tarea{
		TareaID:    nod.NodoID,
		HistoriaID: nod.PadreID,
		Posicion:   nod.Posicion,

		Descripcion:  nod.Titulo,
		Impedimentos: nod.Descripcion,

		Prioridad:        nod.Prioridad,
		Estatus:          nod.Estatus,
		SegundosEstimado: nod.Segundos,
	}
}

// ================================================================ //

// Tramo del viaje de usuario
//
// El padre siempre es una historia.
type Tramo struct {
	TramoID    int // nod.NodoID
	HistoriaID int // nod.PadreID
	Posicion   int

	Texto  string // nod.Titulo
	Imagen string
}

func (nod Nodo) ToTramo() Tramo {
	return Tramo{
		TramoID:    nod.NodoID,
		HistoriaID: nod.PadreID,
		Posicion:   nod.Posicion,

		Texto:  nod.Titulo,
		Imagen: nod.Imagen,
	}
}

// ================================================================ //
