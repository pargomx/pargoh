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

	reloader reloader // websocket.go
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

	// ================================================================ //

	if s.cfg.sourceDir != "" {
		s.gecko.StaticAbs("/assets", s.cfg.sourceDir+"/assets")
		s.gecko.File("/favicon.ico", s.cfg.sourceDir+"/assets/img/favicon.ico") // Todo: agregar FileAbs
	} else {
		s.gecko.StaticFS("/assets", assets.AssetsFS)
		s.gecko.FileFS("/favicon.ico", "img/favicon.ico", assets.AssetsFS)
	}

	s.GET("/reload", s.brodcastReload)
	s.GET("/historias/{historia_id}/ws", s.reloader.nuevoWS)

	s.GET("/", s.listaProyectos)
	s.GET("/proyectos/{proyecto_id}", s.getProyecto)
	s.GET("/personas/{persona_id}", s.getPersona)
	s.GET("/historias/{historia_id}", s.getHistoria)

	s.POS("/proyectos", s.postProyecto)
	s.PUT("/proyectos/{proyecto_id}", s.updateProyecto)
	s.DEL("/proyectos/{proyecto_id}", s.deleteProyecto)

	s.POS("/personas", s.postPersona)
	s.PCH("/personas/{persona_id}", s.patchPersona)
	s.DEL("/personas/{persona_id}", s.deletePersona)

	s.GET("/lista/{nodo_id}", s.getHistoriasLista)
	s.GET("/tablero/{nodo_id}", s.getHistoriasTablero)
	s.GET("/prioritarias", s.getHistoriasPrioritarias)

	s.POS("/historias/{historia_id}", s.postHistoriaQuick)
	s.PUT("/historias/{historia_id}", s.putHistoria)
	s.PCH("/historias/{historia_id}/{param}", s.patchHistoria)
	s.DEL("/historias/{historia_id}", s.deleteHistoria)
	s.GET("/historias/{historia_id}/form", s.formHistoria)
	s.GET("/historias/{historia_id}/mover", s.moverHistoriaForm)
	s.POS("/historias/{historia_id}/mover", s.moverHistoria)
	s.POS("/historias/{historia_id}/priorizar", s.priorizarHistoria)
	s.POS("/historias/{historia_id}/marcar", s.marcarHistoria)

	s.GET("/historias/{historia_id}/tareas", s.getTareasDeHistoria)
	s.POS("/historias/{historia_id}/tareas", s.postTarea)

	s.POS("/historias/{historia_id}/viaje", s.postTramoDeViaje)
	s.DEL("/historias/{historia_id}/viaje/{posicion}", s.deleteTramoDeViaje)
	s.PCH("/historias/{historia_id}/viaje/{posicion}", s.patchTramoDeViaje)
	s.POS("/reordenar-tramo", s.reordenarTramo)

	s.POS("/historias/{historia_id}/reglas", s.postRegla)
	s.DEL("/historias/{historia_id}/reglas/{posicion}", s.deleteRegla)
	s.PCH("/historias/{historia_id}/reglas/{posicion}", s.patchRegla)
	s.POS("/reordenar-regla", s.reordenarRegla)

	s.gecko.StaticSub("/imagenes", s.cfg.imagesDir)
	s.POS("/imagenes", s.setImagenTramo)
	s.DEL("/imagenes/{historia_id}/{posicion}", s.deleteImagenTramo)

	s.POS("/nodos/{nodo_id}", s.postHistoria)
	s.POS("/nodos/{nodo_id}/reordenar", s.reordenarNodo)

	s.PCH("/tareas/{tarea_id}", s.modificarTarea)
	s.DEL("/tareas/{tarea_id}", s.eliminarTarea)
	s.POS("/tareas/{tarea_id}/iniciar", s.iniciarTarea)
	s.POS("/tareas/{tarea_id}/pausar", s.pausarTarea)
	s.POS("/tareas/{tarea_id}/terminar", s.terminarTarea)

	s.GET("/intervalos", s.getIntervalos)

	// Exportar e importar
	s.GET("/arbol", s.exportarArbolTXT)
	s.GET("/fake", func(c *gecko.Context) error { return dhistorias.ImportarFake(s.repo) })
	s.POS("/proyectos/importar", s.importarJSON)
	s.GET("/proyectos/{proyecto_id}/exportar.json", s.exportarJSON)
	s.GET("/proyectos/{proyecto_id}/exportar.md", s.exportarMarkdown)
	s.GET("/proyectos/{proyecto_id}/exportar.docx", s.exportarFile)

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
