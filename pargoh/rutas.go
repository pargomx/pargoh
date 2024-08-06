package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"monorepo/assets"
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
var AMBIENTE string   // Ambiente de ejecución [ dev / prod ]

type configs struct {
	puerto       int    // Puerto TCP del servidor
	directorio   string // Default: directorio actual
	databasePath string // Default: _pargo/pargo.sqlite
}

type servidor struct {
	cfg   configs
	gecko *gecko.Gecko
	db    *sqlitedb.SqliteDB
	repo  *sqliteust.Repositorio
}

func main() {
	gko.LogInfof("Versión:%s:%s", BUILD_INFO, AMBIENTE)

	s := servidor{
		gecko: gecko.New(),
	}

	// Parámetros de ejecución
	flag.StringVar(&s.cfg.directorio, "dir", "", "directorio raíz de la aplicación")
	flag.StringVar(&s.cfg.databasePath, "db", "historias.db", "ubicación de la db sqlite")
	flag.IntVar(&s.cfg.puerto, "p", 5050, "el servidor escuchará en este puerto")
	flag.Parse()
	if s.cfg.directorio != "" {
		err := os.Chdir(s.cfg.directorio)
		if err != nil {
			gko.FatalExit("directorio de proyecto inválido: " + err.Error())
		}
	}
	var err error

	// Repositorio
	s.db, err = sqlitedb.NuevoRepositorio(s.cfg.databasePath, migraciones.MigracionesFS)
	if err != nil {
		gko.FatalError(err)
	}
	s.repo = sqliteust.NuevoRepo(s.db)

	tpls, err := plantillas.NuevoServicioPlantillasEmbebidas(htmltmpl.PlantillasFS, "plantillas")
	if err != nil {
		gko.FatalError(err)
	}
	s.gecko.Renderer = tpls
	s.gecko.TmplBaseLayout = "app/layout"

	// s.exportarFile()

	// ================================================================ //

	s.gecko.StaticFS("/assets", assets.AssetsFS)
	s.gecko.FileFS("/favicon.ico", "img/favicon.ico", assets.AssetsFS)

	s.GET("/", s.getPersonas)
	s.GET("/personas", s.getPersonas)
	s.POS("/personas", s.postPersona)
	s.PCH("/personas/{persona_id}", s.patchPersona)
	s.DEL("/personas/{persona_id}", s.deletePersona)

	s.GET("/arbol", s.getArbolCompleto)
	s.GET("/lista/{nodo_id}", s.getHistoriasLista)
	s.GET("/tablero/{nodo_id}", s.getHistoriasTablero)
	s.GET("/prioritarias", s.getHistoriasPrioritarias)
	s.GET("/historias/{historia_id}", s.getTareasDeHistoria)

	s.PCH("/historias/{historia_id}", s.patchHistoria)
	s.DEL("/historias/{historia_id}", s.deleteHistoria)
	s.GET("/historias/{historia_id}/form", s.formHistoria)
	s.GET("/historias/{historia_id}/mover", s.moverHistoriaForm)
	s.POS("/historias/{historia_id}/mover", s.moverHistoria)
	s.POS("/historias/{historia_id}/priorizar", s.priorizarHistoria)
	s.POS("/historias/{historia_id}/marcar", s.marcarHistoria)
	s.POS("/historias/{historia_id}/tareas", s.postTarea)

	s.POS("/nodos/{nodo_id}", s.postHistoria)
	s.POS("/nodos/{nodo_id}/reordenar", s.reordenarNodo)

	s.PCH("/tareas/{tarea_id}", s.modificarTarea)
	s.POS("/tareas/{tarea_id}/iniciar", s.iniciarTarea)
	s.POS("/tareas/{tarea_id}/pausar", s.pausarTarea)
	s.POS("/tareas/{tarea_id}/terminar", s.terminarTarea)

	s.GET("/intervalos", s.getIntervalos)

	s.GET("/export.md", s.exportarMarkdown)
	s.GET("/export.docx", s.exportarFile)

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
