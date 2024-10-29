package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"monorepo/assets"
	"monorepo/dhistorias"
	"monorepo/htmltmpl"
	"monorepo/migraciones"
	"monorepo/sqlitedb"
	"monorepo/sqliteust"

	"github.com/pargomx/gecko"
	"github.com/pargomx/gecko/gko"
	"github.com/pargomx/gecko/plantillas"
)

// Información de compilación establecida con:
//
//	BUILD_INFO="$(date -I):$(git log --format="%H" -n 1)"
//	go BUILD_INFO -ldflags "-X main.BUILD_INFO=$BUILD_INFO -X main.ambiente=dev"
var BUILD_INFO string // Información de compilación [ fecha:commit_hash ]
var AMBIENTE string   // Ambiente de ejecución [ DEV / PROD ]

type configs struct {
	puerto       int    // Puerto TCP del servidor
	directorio   string // Default: directorio actual
	databasePath string // Default: _pargo/pargo.sqlite
	logDB        bool   // Log de consultas a la base de datos
	sourceDir    string // Directorio raíz para leer assets y plantillas (shadow embed)
	imagesDir    string // Directorio para guardar imágenes
}

type servidor struct {
	cfg   configs
	gecko *gecko.Gecko
	db    *sqlitedb.SqliteDB
	repo  *sqliteust.Repositorio
	auth  *authService

	reloader reloader // websocket.go

	timeTracker *dhistorias.GestionTimeTracker
}

func main() {
	gko.LogInfof("Versión:%s:%s", BUILD_INFO, AMBIENTE)

	s := servidor{}

	// Parámetros de ejecución
	flag.StringVar(&s.cfg.directorio, "dir", "", "directorio raíz de la aplicación")
	flag.StringVar(&s.cfg.databasePath, "db", "historias.db", "ubicación de la db sqlite")
	flag.IntVar(&s.cfg.puerto, "p", 5050, "el servidor escuchará en este puerto")
	flag.BoolVar(&s.cfg.logDB, "logdb", false, "log de consultas a la base de datos")
	flag.StringVar(&s.cfg.sourceDir, "src", "", "directorio con assets y htmltmpl para no usar embeded")
	flag.StringVar(&s.cfg.imagesDir, "img", "imagenes", "directorio con las imágenes de historias y proyectos")

	flag.Parse()
	if s.cfg.directorio != "" {
		err := os.Chdir(s.cfg.directorio)
		if err != nil {
			gko.FatalExit("directorio de proyecto inválido: " + err.Error())
		}
	}
	s.gecko = gecko.New()
	var err error

	// Repositorio
	s.db, err = sqlitedb.NuevoRepositorio(s.cfg.databasePath, migraciones.MigracionesFS)
	if err != nil {
		gko.FatalError(err)
	}
	if s.cfg.logDB {
		s.db.ToggleLog()
	}
	s.repo = sqliteust.NuevoRepo(s.db)
	s.timeTracker = dhistorias.NewGestionTimeTracker(s.repo, 0)

	if s.cfg.sourceDir != "" {
		gko.LogInfo("Usando plantillas y assets " + s.cfg.sourceDir)
		s.gecko.Renderer, err = plantillas.NuevoServicioPlantillas(s.cfg.sourceDir+"/htmltmpl", AMBIENTE == "DEV")
	} else {
		s.gecko.Renderer, err = plantillas.NuevoServicioPlantillasEmbebidas(htmltmpl.PlantillasFS, "plantillas")
	}
	if err != nil {
		gko.FatalError(err)
	}

	if err = s.verificarDirectorioImagenes(); err != nil {
		gko.FatalError(err)
	}

	s.gecko.TmplBaseLayout = "app/layout"

	s.auth = NewAuthService()

	// ================================================================ //

	if s.cfg.sourceDir != "" {
		s.gecko.StaticAbs("/assets", s.cfg.sourceDir+"/assets")
		s.gecko.FileAbs("/favicon.ico", s.cfg.sourceDir+"/assets/img/favicon.ico")
		s.gecko.FileAbs("/htmx.js", s.cfg.sourceDir+"/assets/js/htmx.js")
	} else {
		s.gecko.StaticFS("/assets", assets.AssetsFS)
		s.gecko.FileFS("/favicon.ico", "img/favicon.ico", assets.AssetsFS)
		s.gecko.FileFS("/htmx.js", "js/htmx.min.js", assets.AssetsFS)
	}

	// Sesiones
	s.GET("/", s.auth.getLogin)
	s.POS("/login", s.auth.postLogin)
	s.GET("/logout", s.auth.logout)
	s.GET("/sesiones", s.auth.printSesiones)

	// Proyectos
	s.GET("/proyectos", s.listaProyectos)
	s.POS("/proyectos", s.postProyecto)
	s.GET("/proyectos/{proyecto_id}", s.getProyecto)
	s.GET("/proyectos/{proyecto_id}/metricas", s.getMetricasProyecto)
	s.DEL("/proyectos/{proyecto_id}", s.deleteProyecto)
	s.DEL("/proyectos/{proyecto_id}/definitivo", s.deleteProyectoPorCompleto)
	s.PUT("/proyectos/{proyecto_id}", s.updateProyecto)
	s.PCH("/proyectos/{proyecto_id}/{param}", s.patchProyecto)

	// Personas
	s.POS("/personas", s.postPersona)
	s.GET("/personas/{persona_id}", s.getPersona)
	s.GET("/personas/{persona_id}/doc", s.getPersonaDoc)
	s.GET("/personas/{persona_id}/debug", s.getPersonaDebug)
	s.POS("/personas/{persona_id}", s.postHistoriaDePersona)
	s.GET("/personas/{persona_id}/metricas", s.getMétricasPersona)
	s.DEL("/personas/{persona_id}", s.deletePersona)
	s.PUT("/personas/{persona_id}", s.updatePersona)
	s.PCH("/personas/{persona_id}/{param}", s.patchPersona)
	s.POS("/personas/{persona_id}/time/{seg}", s.postTimeGestion)
	s.POS("/reordenar-persona", s.reordenarPersona)

	// Historias
	s.GET("/historias/{historia_id}", s.getHistoria)
	s.GET("/historias/{historia_id}/tablero", s.getHistoriaTablero)
	s.POS("/historias/{historia_id}", s.postHistoriaDeHistoria)
	s.DEL("/historias/{historia_id}", s.deleteHistoria)
	s.PUT("/historias/{historia_id}", s.updateHistoria)
	s.PCH("/historias/{historia_id}/{param}", s.patchHistoria)
	s.POS("/historias/{historia_id}/priorizar", s.priorizarHistoria)
	s.POS("/historias/{historia_id}/priorizar/{prioridad}", s.priorizarHistoriaNuevo)
	s.POS("/historias/{historia_id}/marcar", s.marcarHistoria)
	s.POS("/historias/{historia_id}/marcar/{completada}", s.marcarHistoriaNueva)
	s.POS("/reordenar-historia", s.reordenarHistoria)

	s.GET("/historias/{historia_id}/mover", s.moverHistoriaForm)
	s.POS("/historias/{historia_id}/mover", s.moverHistoria)

	// Navegador del árbol de historias
	s.GET("/nav", s.navDesdeRoot)
	s.GET("/nav/proy/{proyecto_id}", s.navDesdeProyecto)
	s.GET("/nav/pers/{persona_id}", s.navDesdePersona)
	s.GET("/nav/hist/{historia_id}", s.navDesdeHistoria)

	s.POS("/mover/tramo", s.moverTramo)
	s.POS("/mover/tarea", s.moverTarea)
	s.POS("/mover/historia", s.moverHistoria)

	// Tareas técnicas
	s.POS("/historias/{historia_id}/tareas", s.postTarea)

	s.GET("/tareas/{tarea_id}", s.getTarea)
	s.PCH("/tareas/{tarea_id}", s.modificarTarea)
	s.DEL("/tareas/{tarea_id}", s.eliminarTarea)
	s.PCH("/tareas/{tarea_id}/estimado", s.cambiarEstimadoTarea)
	s.POS("/tareas/{tarea_id}/importancia", s.ciclarImportanciaTarea)
	s.POS("/tareas/{tarea_id}/iniciar", s.iniciarTarea)
	s.POS("/tareas/{tarea_id}/pausar", s.pausarTarea)
	s.POS("/tareas/{tarea_id}/terminar", s.terminarTarea)

	s.GET("/intervalos", s.getIntervalos)
	s.PCH("/tareas/{tarea_id}/intervalos/{inicio}", s.patchIntervalo)

	// Viaje de usuario
	s.POS("/historias/{historia_id}/viaje", s.postTramoDeViaje)
	s.DEL("/historias/{historia_id}/viaje/{posicion}", s.deleteTramoDeViaje)
	s.PCH("/historias/{historia_id}/viaje/{posicion}", s.patchTramoDeViaje)
	s.POS("/reordenar-tramo", s.reordenarTramo)

	s.gecko.StaticSub("/imagenes", s.cfg.imagesDir)
	s.POS("/imagenes", s.setImagenTramo)
	s.DEL("/imagenes/{historia_id}/{posicion}", s.deleteImagenTramo)

	// Reglas de negocio
	s.POS("/historias/{historia_id}/reglas", s.postRegla)
	s.DEL("/historias/{historia_id}/reglas/{posicion}", s.deleteRegla)
	s.PCH("/historias/{historia_id}/reglas/{posicion}", s.patchRegla)
	s.POS("/reordenar-regla", s.reordenarRegla)

	// Exportar e importar
	s.GET("/arbol", s.exportarArbolTXT)
	s.GET("/fake", func(c *gecko.Context) error { return dhistorias.ImportarFake(s.repo) })
	s.POS("/proyectos/importar", s.importarJSON)
	s.GET("/proyectos/{proyecto_id}/exportar.json", s.exportarJSON)
	s.GET("/proyectos/{proyecto_id}/exportar.md", s.exportarMarkdown)
	s.GET("/proyectos/{proyecto_id}/exportar.docx", s.exportarDocx)
	s.GET("/proyectos/{proyecto_id}/exportar.tex", s.exportarProyectoTeX)
	s.GET("/proyectos/{proyecto_id}/exportar.pdf", s.exportarPDF)
	s.GET("/personas/{persona_id}/exportar.pdf", s.exportarPersonaPDF)

	// General
	s.GET("/metricas", s.getMétricas1)
	s.GET("/metricas2", s.getMétricas2)
	s.GET("/materializar-tiempos", s.materializarTiemposTareas)
	s.GET("/materializar-historias", s.materializarHistorias)

	s.GET("/reload", s.brodcastReload)
	s.GET("/historias/{historia_id}/ws", s.reloader.nuevoWS)

	// LOG SQLITE
	s.GET("/log", func(c *gecko.Context) error { s.db.ToggleLog(); return c.StatusOk("Log toggled") })

	// ================================================================ //
	// ================================================================ //

	// Handle interrupt.
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for sig := range ch {
			err = s.db.Close()
			if err != nil {
				fmt.Println("sqliteDB.Close: ", err.Error())
			}
			fmt.Println("")
			gko.LogInfof("servidor terminado: %v", sig.String())
			os.Exit(0)
		}
	}()

	err = s.gecko.IniciarEnPuerto(s.cfg.puerto)
	if err != nil {
		gko.FatalError(err)
	}
}
