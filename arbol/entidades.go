package arbol

// Grupo de proyectos para organizarlos en carpetas.
//
// Por default se generan dos grupos:
//   - "Activos" para los proyectos que se muestran en la página de inicio.
//   - "Archivados" para los que se ocultan de la página de inicio.
type Grupo struct {
	GrupoID  int
	PadreID  int
	Posicion int

	Nombre string
}

// Proyecto representa el esfuerzo de desarrollar una o varias apps.
type Proyecto struct {
	ProyectoID int // nodo "PRY"
	PadreID    int
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
type Persona struct {
	PersonaID int // nodo "PER"
	PadreID   int // antes ProyectoID
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
type HistoriaDeUsuario struct {
	HistoriaID int // `historias.historia_id`
	PadreID    int // `historias.padre_id`  Historia padre para el árbol
	Posicion   int // `historias.posicion`  Posición consecutiva con respecto a sus nodos hermanos

	Titulo      string // `historias.titulo`
	Objetivo    string // `historias.objetivo`
	Descripcion string // `historias.descripcion`  Descripción  de la historia en infinitivo para que la lea el usuario en la documentación.
	Notas       string // `historias.notas`  Apuntes técnicos sobre la implementación de la historia.

	Prioridad           int  // `historias.prioridad`
	Completada          bool // `historias.completada`
	SegundosPresupuesto int  // `historias.segundos_presupuesto`  Tiempo estimado en segundos para implementar la historia de usuario en su totalidad
	// Centavos            int

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
type HistoriaTécnica struct {
	HistoriaID int // `historias.historia_id`
	PadreID    int // `historias.padre_id`  Historia padre para el árbol
	Posicion   int // `historias.posicion`  Posición consecutiva con respecto a sus nodos hermanos

	Titulo      string // `historias.titulo`
	Objetivo    string // `historias.objetivo`
	Descripcion string // `historias.descripcion`  Descripción  de la historia en infinitivo para que la lea el usuario en la documentación.
	Notas       string // `historias.notas`  Apuntes técnicos sobre la implementación de la historia.

	Prioridad  int  // `historias.prioridad`
	Completada bool // `historias.completada`

	SegundosPresupuesto   int // `historias.segundos_presupuesto`  Tiempo estimado en segundos para implementar la historia de usuario en su totalidad
	SegundosDocumentacion int // `historias.segundos_documentacion`
	SegundosUtilizado     int // `historias.segundos_utilizado`
}

// Actividades de gestión parte del ciclo de vida de desarrollo y mantenimiento
// de la aplicación:  documentación, juntas, soporte técnico, proceso de venta, etc.
type ActividadDeGestión struct {
	HistoriaID int // `historias.historia_id`
	PadreID    int // `historias.padre_id`  Historia padre para el árbol
	Posicion   int // `historias.posicion`  Posición consecutiva con respecto a sus nodos hermanos

	Titulo      string // `historias.titulo`
	Objetivo    string // `historias.objetivo`
	Descripcion string // `historias.descripcion`  Descripción  de la historia en infinitivo para que la lea el usuario en la documentación.
	Notas       string // `historias.notas`  Apuntes técnicos sobre la implementación de la historia.

	Prioridad  int  // `historias.prioridad`
	Completada bool // `historias.completada`

	SegundosPresupuesto   int // `historias.segundos_presupuesto`  Tiempo estimado en segundos para implementar la historia de usuario en su totalidad
	SegundosDocumentacion int // `historias.segundos_documentacion`
	SegundosUtilizado     int // `historias.segundos_utilizado`
}

// Regla de negocio
type Regla struct {
	NodoID     int
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

	Prioridad int // Importancia
	Estatus   int
	Segundos  int
	Centavos  int

	// Legacy
	SegundosEstimado  int
	SegundosUtilizado int
}

type Tramo struct {
	NodoID     int
	HistoriaID int // PadreID
	Posicion   int

	Texto  string
	Imagen string
}
