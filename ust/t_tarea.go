package ust

import (
	"errors"

	"github.com/pargomx/gecko/gko"
)

// Tarea corresponde a un elemento de la tabla 'tareas'.
type Tarea struct {
	TareaID           int              // `tareas.tarea_id`
	HistoriaID        int              // `tareas.historia_id`
	Descripcion       string           // `tareas.descripcion`
	Importancia       ImportanciaTarea // `tareas.importancia`
	Tipo              TipoTarea        // `tareas.tipo`
	Estatus           int              // `tareas.estatus`  0:no_iniciada 1:en_curso 2:en_pausa 3:finalizada
	Impedimentos      string           // `tareas.impedimentos`
	SegundosEstimado  int              // `tareas.segundos_estimado`
	SegundosUtilizado int              // `tareas.segundos_real`  Dato materializado que debe sumar el tiempo para todos los intervalos de la tarea.
}

var (
	ErrTareaNotFound      error = errors.New("la tarea no se encuentra")
	ErrTareaAlreadyExists error = errors.New("la tarea ya existe")
)

func (tar *Tarea) Validar() error {

	if tar.Importancia.EsTodos() {
		return errors.New("ust.Tarea no admite propiedad ImportanciaTodos")
	}

	if tar.Tipo.EsTodos() {
		return errors.New("ust.Tarea no admite propiedad TipoTodos")
	}

	return nil
}

//  ================================================================  //
//  ========== Importancia =========================================  //

// Enumeración
type ImportanciaTarea struct {
	ID          int
	String      string
	Filtro      string
	Etiqueta    string
	Descripcion string
}

var (
	// ImportanciaTareaTodos solo se utiliza como filtro.
	ImportanciaTareaTodos = ImportanciaTarea{
		ID:          -1,
		String:      "",
		Filtro:      "todos",
		Etiqueta:    "Todos",
		Descripcion: "Todos los valores posibles para Importancia",
	}
	// Indica explícitamente que la propiedad está indefinida.
	ImportanciaTareaIndefinido = ImportanciaTarea{
		ID:          0,
		String:      "",
		Filtro:      "sin_importancia",
		Etiqueta:    "Indefinido",
		Descripcion: "Indefinido",
	}

	// Idea: Algo que aún no está bien planteado o que no aporta suficiente valor y por lo tanto no cuenta mucho para el progreso de la historia de usuario
	ImportanciaTareaIdea = ImportanciaTarea{
		ID:          1,
		String:      "IDEA",
		Filtro:      "idea",
		Etiqueta:    "Idea",
		Descripcion: "Algo que aún no está bien planteado o que no aporta suficiente valor y por lo tanto no cuenta mucho para el progreso de la historia de usuario",
	}
	// Mejora: Una mejora a la UI, UX, optimizaciones, etc.
	ImportanciaTareaMejora = ImportanciaTarea{
		ID:          2,
		String:      "MEJORA",
		Filtro:      "mejora",
		Etiqueta:    "Mejora",
		Descripcion: "Una mejora a la UI, UX, optimizaciones, etc.",
	}
	// Necesaria: Indispensable para la implementación de la historia de usuario
	ImportanciaTareaNecesaria = ImportanciaTarea{
		ID:          3,
		String:      "NECESARIA",
		Filtro:      "necesaria",
		Etiqueta:    "Necesaria",
		Descripcion: "Indispensable para la implementación de la historia de usuario",
	}
)

// Enumeración excluyendo ImportanciaTareaTodos
var ListaImportanciaTarea = []ImportanciaTarea{
	ImportanciaTareaIndefinido,

	ImportanciaTareaIdea,
	ImportanciaTareaMejora,
	ImportanciaTareaNecesaria,
}

// Enumeración incluyendo ImportanciaTareaTodos
var ListaFiltroImportanciaTarea = []ImportanciaTarea{
	ImportanciaTareaTodos,
	ImportanciaTareaIndefinido,

	ImportanciaTareaIdea,
	ImportanciaTareaMejora,
	ImportanciaTareaNecesaria,
}

// Comparar un Importancia con otro.
func (a ImportanciaTarea) Es(e ImportanciaTarea) bool {
	return a.ID == e.ID
}

func (e ImportanciaTarea) EsTodos() bool {
	return e.ID == ImportanciaTareaTodos.ID
}
func (e ImportanciaTarea) EsIndefinido() bool {
	return e.ID == ImportanciaTareaIndefinido.ID
}
func (e ImportanciaTarea) EsIdea() bool {
	return e.ID == ImportanciaTareaIdea.ID
}
func (e ImportanciaTarea) EsMejora() bool {
	return e.ID == ImportanciaTareaMejora.ID
}
func (e ImportanciaTarea) EsNecesaria() bool {
	return e.ID == ImportanciaTareaNecesaria.ID
}

// Recibe la forma .String
func SetImportanciaTareaDB(str string) ImportanciaTarea {
	for _, e := range ListaImportanciaTarea {
		if e.String == str {
			return e
		}
	}
	if str == ImportanciaTareaTodos.String {
		gko.LogWarn("ust.SetImportanciaTarea: ImportanciaTareaTodos es inválido en DB")
		return ImportanciaTareaIndefinido
	}
	gko.LogWarnf("ust.SetImportanciaTarea inválido: '%v'", str)
	return ImportanciaTareaIndefinido
}

// Recibe la forma .Filtro
func SetImportanciaTareaFiltro(str string) ImportanciaTarea {
	if str == "" || str == ImportanciaTareaTodos.Filtro {
		return ImportanciaTareaTodos
	}
	for _, e := range ListaImportanciaTarea {
		if e.Filtro == str {
			return e
		}
	}
	gko.LogWarnf("ust.SetImportanciaTarea inválido: '%v'", str)
	return ImportanciaTareaIndefinido
}

// Recibe la forma .String o .Filtro
func (i *Tarea) SetImportancia(str string) {
	for _, e := range ListaImportanciaTarea {
		if e.String == str {
			i.Importancia = e
			return
		}
		if e.Filtro == str {
			i.Importancia = e
			return
		}
	}
	if str == ImportanciaTareaTodos.String {
		gko.LogWarn("ust.SetImportanciaTarea: ImportanciaTareaTodos es inválido en DB")
		i.Importancia = ImportanciaTareaIndefinido
	}
	gko.LogWarnf("ust.SetImportanciaTarea inválido: '%v'", str)
	i.Importancia = ImportanciaTareaIndefinido
}

//  ================================================================  //
//  ========== Tipo de tarea =======================================  //

// Enumeración
type TipoTarea struct {
	ID          int
	String      string
	Filtro      string
	Etiqueta    string
	Descripcion string
}

var (
	// TipoTareaTodos solo se utiliza como filtro.
	TipoTareaTodos = TipoTarea{
		ID:          -1,
		String:      "",
		Filtro:      "todos",
		Etiqueta:    "Todos",
		Descripcion: "Todos los valores posibles para Tipo de tarea",
	}
	// Indica explícitamente que la propiedad está indefinida.
	TipoTareaIndefinido = TipoTarea{
		ID:          0,
		String:      "",
		Filtro:      "sin_tipo",
		Etiqueta:    "Indefinido",
		Descripcion: "Indefinido",
	}

	// Planeación: Planeación y estimación de tareas del sprint
	TipoTareaPlan = TipoTarea{
		ID:          1,
		String:      "PLAN",
		Filtro:      "plan",
		Etiqueta:    "Planeación",
		Descripcion: "Planeación y estimación de tareas del sprint",
	}
	// Conocer dominio: Construcción de conocimiento de dominio (entrevistas, etc)
	TipoTareaAsk = TipoTarea{
		ID:          2,
		String:      "ASK",
		Filtro:      "ask",
		Etiqueta:    "Conocer dominio",
		Descripcion: "Construcción de conocimiento de dominio (entrevistas, etc)",
	}
	// Diseño del dominio: Diseño de dominio (entidades, agregados, value objects, comandos, consultas, eventos)
	TipoTareaDominio = TipoTarea{
		ID:          3,
		String:      "DOMINIO",
		Filtro:      "dominio",
		Etiqueta:    "Diseño del dominio",
		Descripcion: "Diseño de dominio (entidades, agregados, value objects, comandos, consultas, eventos)",
	}
	// Base de datos: Diseño de esquema de base de datos
	TipoTareaDb = TipoTarea{
		ID:          4,
		String:      "DB",
		Filtro:      "db",
		Etiqueta:    "Base de datos",
		Descripcion: "Diseño de esquema de base de datos",
	}
	// Configuración: Configuración de ambientes
	TipoTareaConf = TipoTarea{
		ID:          5,
		String:      "CONF",
		Filtro:      "conf",
		Etiqueta:    "Configuración",
		Descripcion: "Configuración de ambientes",
	}
	// Repositorios: Codificar adaptadores de repositorio
	TipoTareaRepo = TipoTarea{
		ID:          6,
		String:      "REPO",
		Filtro:      "repo",
		Etiqueta:    "Repositorios",
		Descripcion: "Codificar adaptadores de repositorio",
	}
	// Adaptadores: Codificar adaptadores para servicios externos (APIs, FS, etc)
	TipoTareaAdapt = TipoTarea{
		ID:          7,
		String:      "ADAPT",
		Filtro:      "adapt",
		Etiqueta:    "Adaptadores",
		Descripcion: "Codificar adaptadores para servicios externos (APIs, FS, etc)",
	}
	// Handlers: Codificar controladores handlers para solicitudes HTTP
	TipoTareaHandlr = TipoTarea{
		ID:          8,
		String:      "HANDLR",
		Filtro:      "handlr",
		Etiqueta:    "Handlers",
		Descripcion: "Codificar controladores handlers para solicitudes HTTP",
	}
	// Interfaz web: Codificar interfaz web (HTML, CSS, JS)
	TipoTareaWebUi = TipoTarea{
		ID:          9,
		String:      "WEB_UI",
		Filtro:      "web_ui",
		Etiqueta:    "Interfaz web",
		Descripcion: "Codificar interfaz web (HTML, CSS, JS)",
	}
	// Bug: Arreglar defecto en el programa
	TipoTareaBug = TipoTarea{
		ID:          10,
		String:      "BUG",
		Filtro:      "bug",
		Etiqueta:    "Bug",
		Descripcion: "Arreglar defecto en el programa",
	}
)

// Enumeración excluyendo TipoTareaTodos
var ListaTipoTarea = []TipoTarea{
	TipoTareaIndefinido,

	TipoTareaPlan,
	TipoTareaAsk,
	TipoTareaDominio,
	TipoTareaDb,
	TipoTareaConf,
	TipoTareaRepo,
	TipoTareaAdapt,
	TipoTareaHandlr,
	TipoTareaWebUi,
	TipoTareaBug,
}

// Enumeración incluyendo TipoTareaTodos
var ListaFiltroTipoTarea = []TipoTarea{
	TipoTareaTodos,
	TipoTareaIndefinido,

	TipoTareaPlan,
	TipoTareaAsk,
	TipoTareaDominio,
	TipoTareaDb,
	TipoTareaConf,
	TipoTareaRepo,
	TipoTareaAdapt,
	TipoTareaHandlr,
	TipoTareaWebUi,
	TipoTareaBug,
}

// Comparar un Tipo de tarea con otro.
func (a TipoTarea) Es(e TipoTarea) bool {
	return a.ID == e.ID
}

func (e TipoTarea) EsTodos() bool {
	return e.ID == TipoTareaTodos.ID
}
func (e TipoTarea) EsIndefinido() bool {
	return e.ID == TipoTareaIndefinido.ID
}
func (e TipoTarea) EsPlan() bool {
	return e.ID == TipoTareaPlan.ID
}
func (e TipoTarea) EsAsk() bool {
	return e.ID == TipoTareaAsk.ID
}
func (e TipoTarea) EsDominio() bool {
	return e.ID == TipoTareaDominio.ID
}
func (e TipoTarea) EsDb() bool {
	return e.ID == TipoTareaDb.ID
}
func (e TipoTarea) EsConf() bool {
	return e.ID == TipoTareaConf.ID
}
func (e TipoTarea) EsRepo() bool {
	return e.ID == TipoTareaRepo.ID
}
func (e TipoTarea) EsAdapt() bool {
	return e.ID == TipoTareaAdapt.ID
}
func (e TipoTarea) EsHandlr() bool {
	return e.ID == TipoTareaHandlr.ID
}
func (e TipoTarea) EsWebUi() bool {
	return e.ID == TipoTareaWebUi.ID
}
func (e TipoTarea) EsBug() bool {
	return e.ID == TipoTareaBug.ID
}

// Recibe la forma .String
func SetTipoTareaDB(str string) TipoTarea {
	for _, e := range ListaTipoTarea {
		if e.String == str {
			return e
		}
	}
	if str == TipoTareaTodos.String {
		gko.LogWarn("ust.SetTipoTarea: TipoTareaTodos es inválido en DB")
		return TipoTareaIndefinido
	}
	gko.LogWarnf("ust.SetTipoTarea inválido: '%v'", str)
	return TipoTareaIndefinido
}

// Recibe la forma .Filtro
func SetTipoTareaFiltro(str string) TipoTarea {
	if str == "" || str == TipoTareaTodos.Filtro {
		return TipoTareaTodos
	}
	for _, e := range ListaTipoTarea {
		if e.Filtro == str {
			return e
		}
	}
	gko.LogWarnf("ust.SetTipoTarea inválido: '%v'", str)
	return TipoTareaIndefinido
}

// Recibe la forma .String o .Filtro
func (i *Tarea) SetTipo(str string) {
	for _, e := range ListaTipoTarea {
		if e.String == str {
			i.Tipo = e
			return
		}
		if e.Filtro == str {
			i.Tipo = e
			return
		}
	}
	if str == TipoTareaTodos.String {
		gko.LogWarn("ust.SetTipoTarea: TipoTareaTodos es inválido en DB")
		i.Tipo = TipoTareaIndefinido
	}
	gko.LogWarnf("ust.SetTipoTarea inválido: '%v'", str)
	i.Tipo = TipoTareaIndefinido
}
