package arbol

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
}

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
}

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

	// LEGACY
	ProyectoID      string
	SegundosGestion int
}

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

	// Legacy
	PersonaID  int
	ProyectoID string
}

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
}

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
}

// Regla de negocio.
//
// El padre siempre es una historia.
type Regla struct {
	ReglaID    int // NodoID
	HistoriaID int // PadreID
	Posicion   int

	Texto string
	// Descripcion string
	// Objetivo    string
	// Notas       string

	// Color  string
	// Imagen string

	// Prioridad int
	Estatus int
	// Segundos int
	// Centavos int

	// Legacy
	Implementada bool // Estatus
	Probada      bool // Estatus
}

// Tarea de desarrollo
//
// El padre siempre es una historia.
type Tarea struct {
	TareaID    int // NodoID
	HistoriaID int // PadreID
	Posicion   int

	Descripcion  string
	Impedimentos string
	// Objetivo    string // TipoTarea
	// Notas       string // SegundosEstimado

	// 	Color  string
	// Imagen string

	Prioridad        int // Importancia
	Estatus          int
	SegundosEstimado int
	// Centavos         int

	// Legacy
	SegundosUtilizado int
}

// Tramo del viaje de usuario
//
// El padre siempre es una historia.
type Tramo struct {
	TramoID    int // NodoID
	HistoriaID int // PadreID
	Posicion   int

	Texto  string
	Imagen string
}
