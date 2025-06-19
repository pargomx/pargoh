package main

import (
	"flag"
	"os"

	"monorepo/arbol"
	"monorepo/assets"
	"monorepo/dhistorias"
	"monorepo/exportdocx"
	"monorepo/htmltmpl"
	"monorepo/migraciones"
	"monorepo/sqlitearbol"
	"monorepo/sqlitepuente"

	"github.com/pargomx/gecko"
	"github.com/pargomx/gecko/gko"
	"github.com/pargomx/gecko/plantillas"
	"github.com/pargomx/gecko/sqlitedb"
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
	sourceDir    string // Directorio raíz para leer assets y plantillas (shadow embed)
	imagesDir    string // Directorio para guardar imágenes
	exportDir    string // Directorio para guardar exports
	unidocApiKey string // API Key para unidoc

	adminUser string
	adminPass string
	debug     debugConfig
}

type debugConfig struct {
	logDB      bool // Log de consultas a la base de datos
	writeDelay int  // ms de delay para solicitudes POST, PUT, PATCH, DELETE
	readDelay  int  // ms de delay para solicitudes GET
}

type servidor struct {
	cfg   configs
	gecko *gecko.Gecko
	db    *sqlitedb.SqliteDB

	repoOld dhistorias.Repo
	repo    arbol.ReadRepo

	auth *authService

	app *arbol.Servicio

	reloader reloader // websocket.go

	timeTracker *dhistorias.GestionTimeTracker

	noContinuar bool // feature flag
}

type readhdl struct {
	db      *sqlitedb.SqliteDB
	repo    arbol.ReadRepo
	repoOld dhistorias.Repo
}

type writehdl struct {
	db          *sqlitedb.SqliteDB
	app         *arbol.Servicio
	reloader    reloader // websocket.go
	timeTracker *dhistorias.GestionTimeTracker
}

func main() {
	gko.LogInfof("Versión:%s:%s", BUILD_INFO, AMBIENTE)

	s := servidor{}

	// Parámetros de ejecución
	flag.StringVar(&s.cfg.directorio, "dir", "", "directorio raíz de la aplicación")
	flag.StringVar(&s.cfg.databasePath, "db", "historias.db", "ubicación de la db sqlite")
	flag.IntVar(&s.cfg.puerto, "p", 5050, "el servidor escuchará en este puerto")
	flag.StringVar(&s.cfg.sourceDir, "src", "", "directorio con assets y htmltmpl para no usar embeded")
	flag.StringVar(&s.cfg.imagesDir, "img", "imagenes", "directorio con las imágenes de historias y proyectos")
	flag.StringVar(&s.cfg.exportDir, "exp", "exports", "directorio con los archivos exportados")
	flag.StringVar(&s.cfg.unidocApiKey, "unidoc", "", "api key para exportar docx con unidoc")
	flag.StringVar(&s.cfg.adminUser, "auser", "tulio", "usuario del administrador")
	flag.StringVar(&s.cfg.adminPass, "apass", "flores99leetcode", "contraseña del administrador")
	flag.BoolVar(&s.cfg.debug.logDB, "logdb", false, "log de consultas a la base de datos")
	flag.IntVar(&s.cfg.debug.writeDelay, "wdelay", 300, "ms de delay para solicitudes POST, PUT, PATCH, DELETE")
	flag.IntVar(&s.cfg.debug.readDelay, "rdelay", 300, "ms de delay para solicitudes GET")

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
	if s.cfg.debug.logDB {
		s.db.ToggleLog()
	}
	s.repoOld = sqlitepuente.NuevoRepo(s.db)
	s.repo = sqlitearbol.NuevoRepo(s.db)
	s.timeTracker = dhistorias.NewGestionTimeTracker(s.repoOld, 0)

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

	if err = exportdocx.VerificarDirectorioExports(s.cfg.exportDir); err != nil {
		gko.FatalError(err)
	}

	s.gecko.TmplBaseLayout = "app/layout"

	s.auth = NewAuthService(s.cfg.adminUser, s.cfg.adminPass)
	s.auth.RecuperarSesiones()

	s.app, err = arbol.NuevoServicio(arbol.Config{})
	if err != nil {
		gko.FatalError(err)
	}

	r := readhdl{
		db:      s.db,
		repo:    sqlitearbol.NuevoRepo(s.db),
		repoOld: sqlitepuente.NuevoRepo(s.db),
	}
	w := writehdl{
		db:          s.db,
		app:         s.app,
		reloader:    s.reloader,
		timeTracker: s.timeTracker,
	}

	// ================================================================ //

	if s.cfg.sourceDir != "" {
		s.gecko.StaticAbs("/assets", s.cfg.sourceDir+"/assets")
		s.gecko.FileAbs("/favicon.ico", s.cfg.sourceDir+"/assets/img/favicon.ico")
		s.gecko.FileAbs("/service-worker.js", s.cfg.sourceDir+"/assets/js/service-worker.js")
		s.gecko.FileAbs("/pargo.webmanifest", s.cfg.sourceDir+"/assets/manifest.json")
	} else {
		s.gecko.StaticFS("/assets", assets.AssetsFS)
		s.gecko.FileFS("/favicon.ico", "img/favicon.ico", assets.AssetsFS)
		s.gecko.FileFS("/service-worker.js", "js/service-worker.js", assets.AssetsFS)
		s.gecko.FileFS("/pargo.webmanifest", "manifest.json", assets.AssetsFS)
	}
	s.gecko.GET("/assets/js/htmx.js", s.gecko.ServirHtmxMinJS())
	s.gecko.GET("/assets/js/gecko.js", s.gecko.ServirGeckoJS())

	// Sesiones
	s.GET("/", s.auth.getLogin)
	s.POS("/login", s.auth.postLogin)
	s.GET("/logout", s.auth.logout)
	s.GET("/sesiones", s.auth.printSesiones)

	s.GET("/buscar", r.buscar)

	s.GET("/continuar", s.continuar)
	s.GET("/offline", s.offline)

	// Proyectos
	s.GET("/proyectos", r.listaProyectos)
	s.GET("/proyectos/{proyecto_id}", r.getProyecto)
	s.GET("/proyectos/{proyecto_id}/doc", r.getDocumentacionProyecto)
	s.DEL("/proyectos/{proyecto_id}", s.deleteProyecto)
	s.DEL("/proyectos/{proyecto_id}/definitivo", s.deleteProyectoPorCompleto)
	s.PUT("/proyectos/{proyecto_id}", s.updateProyecto)
	s.PCH("/proyectos/{proyecto_id}/{param}", s.patchProyecto)

	// Personas
	s.GET("/personas/{persona_id}", r.getPersona)
	s.GET("/personas/{persona_id}/doc", r.getPersonaDoc)
	s.GET("/personas/{persona_id}/debug", r.getPersonaDebug)
	s.DEL("/personas/{persona_id}", s.deletePersona)
	s.PUT("/personas/{persona_id}", s.updatePersona)
	s.PCH("/personas/{persona_id}/{param}", s.patchPersona)
	s.POS("/personas/{persona_id}/time/{seg}", s.postTimeGestion)

	// Historias
	s.GET("/historias/{historia_id}", r.getHistoria)
	s.GET("/historias/{historia_id}/tablero", r.getHistoriaTablero)

	s.DEL("/historias/{historia_id}", s.deleteHistoria)
	s.PUT("/historias/{historia_id}", s.updateHistoria)
	s.PCH("/historias/{historia_id}/{param}", s.patchHistoria)
	s.POS("/historias/{historia_id}/priorizar", s.priorizarHistoria)
	s.POS("/historias/{historia_id}/priorizar/{prioridad}", s.priorizarHistoriaNuevo)
	s.POS("/historias/{historia_id}/marcar", s.marcarHistoria)
	s.POS("/historias/{historia_id}/marcar/{completada}", s.marcarHistoriaNueva)

	// Navegador del árbol de historias
	s.GET("/nav", r.navDesdeRoot)
	s.GET("/nav/proy/{proyecto_id}", r.navDesdeProyecto)
	s.GET("/nav/pers/{persona_id}", r.navDesdePersona)
	s.GET("/nav/hist/{historia_id}", r.navDesdeHistoria)

	s.GET("/historias/{historia_id}/mover", r.moverHistoriaForm)

	// MOVER
	s.POS("/historias/{historia_id}/mover", w.inTx(w.moverHistoria))
	s.POS("/mover/tramo", w.inTx(w.moverTramo))
	s.POS("/mover/tarea", w.inTx(w.moverTarea))
	s.POS("/mover/historia", w.inTx(w.moverHistoria))

	// Tareas técnicas

	s.GET("/tareas/{tarea_id}", r.getTarea)
	s.PCH("/tareas/{tarea_id}", s.modificarTarea)
	s.DEL("/tareas/{tarea_id}", s.eliminarTarea)
	s.PCH("/tareas/{tarea_id}/estimado", s.cambiarEstimadoTarea)
	s.POS("/tareas/{tarea_id}/importancia", s.ciclarImportanciaTarea)
	s.POS("/tareas/{tarea_id}/iniciar", s.iniciarTarea)
	s.POS("/tareas/{tarea_id}/pausar", s.pausarTarea)
	s.POS("/tareas/{tarea_id}/terminar", s.terminarTarea)

	s.GET("/intervalos", r.getIntervalos)
	s.PCH("/tareas/{tarea_id}/intervalos/{inicio}", s.patchIntervalo)

	// Quick tasks
	s.GET("/tareas", r.getQuickTasks)

	// Viaje de usuario
	s.DEL("/historias/{historia_id}/viaje/{posicion}", s.deleteTramoDeViaje)
	s.PCH("/historias/{historia_id}/viaje/{posicion}", s.patchTramoDeViaje)

	s.gecko.StaticSub("/imagenes", s.cfg.imagesDir)
	s.POS("/imagenes", s.setImagenTramo)
	s.DEL("/imagenes/{historia_id}/{posicion}", s.deleteImagenTramo)

	// Reglas de negocio
	s.DEL("/historias/{historia_id}/reglas/{posicion}", s.deleteRegla)
	s.PCH("/historias/{historia_id}/reglas/{posicion}", s.patchRegla)
	s.PCH("/historias/{historia_id}/reglas/{posicion}/marcar", s.marcarRegla)

	// Referencias
	s.POS("/historias/{historia_id}/referencias", s.postReferencia)
	s.DEL("/historias/{historia_id}/referencias/{ref_historia_id}", s.deleteReferencia)

	// AGREGAR HOJA
	s.POS("/proyectos", w.inTx(w.postProyecto))
	s.POS("/personas", w.inTx(w.postPersona))
	s.POS("/personas/{persona_id}", w.inTx(w.postHistoriaDePersona))
	s.POS("/historias/{historia_id}", w.inTx(w.postHistoriaDeHistoria))
	s.POS("/historias/{historia_id}/padre", w.inTx(w.postPadreParaHistoria))
	s.POS("/historias/{historia_id}/reglas", w.inTx(w.postRegla))
	s.POS("/historias/{historia_id}/viaje", w.inTx(w.postTramoDeViaje))
	s.POS("/historias/{historia_id}/tareas", w.inTx(w.postTarea))
	s.POS("/tareas", w.inTx(w.postQuickTask))

	// REORDENAR
	s.POS("/reordenar-persona", w.inTx(w.reordenarPersona))
	s.POS("/reordenar-historia", w.inTx(w.reordenarHistoria))
	s.POS("/reordenar-tramo", w.inTx(w.reordenarTramo))
	s.POS("/reordenar-regla", w.inTx(w.reordenarRegla))

	// Exportar e importar
	/*
		s.GET("/arbol", s.exportarArbolTXT)
		s.GET("/fake", func(c *gecko.Context) error { return dhistorias.ImportarFake(s.repo) })
		s.POS("/proyectos/importar", s.importarJSON)
		s.GET("/proyectos/{proyecto_id}/exportar.json", s.exportarJSON)
		s.GET("/proyectos/{proyecto_id}/exportar.md", s.exportarMarkdown)
			s.GET("/proyectos/{proyecto_id}/exportar.docx", s.exportarProyectoDocx)
			s.GET("/proyectos/{proyecto_id}/exportar.tex", s.exportarProyectoTeX)
			s.GET("/proyectos/{proyecto_id}/exportar.pdf", s.exportarPDF)
			s.GET("/personas/{persona_id}/exportar.pdf", s.exportarPersonaPDF)
			s.POS("/personas/{persona_id}/docx", s.exportarPersonaDocx(s.cfg.unidocApiKey))
			s.gecko.StaticSub("/exports", s.cfg.exportDir) // TODO: autenticar
	*/

	// General
	s.GET("/metricas", r.getMétricas)
	// s.GET("/metricas2", s.getMétricas2)
	// s.GET("/materializar-tiempos", s.materializarTiemposTareas)
	// s.GET("/materializar-historias", s.materializarHistorias)

	s.GET("/reload", s.brodcastReload)
	s.GET("/historias/{historia_id}/ws", s.reloader.nuevoWS)

	// Mantenimiento
	s.GET("/log", func(c *gecko.Context) error { s.db.ToggleLog(); return c.StatusOk("Log toggled") })
	s.GET("/clear", func(c *gecko.Context) error {
		c.Response().Header().Set("Clear-Site-Data", `"cache", "cache", "clientHints", "storage", "executionContexts"`)
		return c.StringOk("Datos del sitio limpiados. Ok.")
	})

	// ================================================================ //
	// ================================================================ //

	s.gecko.CleanupFunc = func() {
		err = s.db.Close()
		if err != nil {
			gko.Op("ShutdownDB").Err(err).Log()
		}
		s.auth.PersistirSesiones()
	}

	err = s.gecko.IniciarEnPuerto(s.cfg.puerto)
	if err != nil {
		gko.FatalError(err)
	}
}
