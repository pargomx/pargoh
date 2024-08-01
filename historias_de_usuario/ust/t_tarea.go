package ust

import (
	"errors"

	"github.com/pargomx/gecko/gko"
)

// Tarea corresponde a un elemento de la tabla 'tareas'.
type Tarea struct {
	TareaID        int       // `tareas.tarea_id`
	HistoriaID     int       // `tareas.historia_id`
	Tipo           TipoTarea // `tareas.tipo`
	Descripcion    string    // `tareas.descripcion`
	Impedimentos   string    // `tareas.impedimentos`
	TiempoEstimado int       // `tareas.tiempo_estimado`
	TiempoReal     int       // `tareas.tiempo_real`
	Estatus        int       // `tareas.estatus`  0:no_iniciada 1:en_curso 2:en_pausa 3:finalizada
}

var (
	ErrTareaNotFound      error = errors.New("la tarea no se encuentra")
	ErrTareaAlreadyExists error = errors.New("la tarea ya existe")
)

func (tar *Tarea) Validar() error {

	if tar.Tipo.EsTodos() {
		return errors.New("ust.Tarea no admite propiedad TipoTodos")
	}

	return nil
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
