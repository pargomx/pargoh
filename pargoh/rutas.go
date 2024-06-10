package main

import (
	"embed"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"monorepo/assets"
	"monorepo/gecko"
	"monorepo/gecko/plantillas"
	"monorepo/historias_de_usuario/sqliteust"
	"monorepo/pargoh/migraciones"
	"monorepo/sqlitedb"
)

//go:embed plantillas
var plantillasFS embed.FS

var puerto = ""       // Default: 5050
var directorio = ""   // Default: directorio actual
var databasePath = "" // Default: _pargo/pargo.sqlite

type servidor struct {
	db   *sqlitedb.SqliteDB
	repo *sqliteust.Repositorio
}

func main() {

	// Setup
	gecko.MostrarMensajeEnErrores = true
	gecko.PrintLogTimestamps = false

	// Parámetros de ejecución
	flag.StringVar(&directorio, "dir", "", "directorio raíz del proyecto")
	flag.StringVar(&databasePath, "db", "_pargo/historias.db", "ubicación de la db sqlite")
	flag.StringVar(&puerto, "p", "5050", "el servidor escuchará en este puerto")
	flag.Parse()
	if directorio != "" {
		err := os.Chdir(directorio)
		if err != nil {
			fatal("directorio de proyecto inválido: " + err.Error())
		}
	}

	// Repositorio
	sqliteDB, err := sqlitedb.NuevoRepositorio(databasePath, migraciones.MigracionesFS)
	if err != nil {
		fatal(err.Error())
	}

	// Servicios
	repos := sqliteust.NuevoRepositorio(sqliteDB)
	srv := &servidor{db: sqliteDB, repo: repos}

	tpls, err := plantillas.NuevoServicioPlantillasEmbebidas(plantillasFS, "plantillas")
	if err != nil {
		fatal(err.Error())
	}
	g := gecko.New()
	g.Renderer = tpls
	g.TmplBaseLayout = "app/layout"

	g.StaticFS("/assets", assets.AssetsFS)
	g.FileFS("/favicon.ico", "img/favicon.ico", assets.AssetsFS)

	// ================================================================ //
	// ================================================================ //

	g.GET("/", getInicio)

	// ================================================================ //

	g.GET("/personas", srv.getPersonas)
	g.POS("/personas", srv.postPersona)
	g.PCH("/personas/{persona_id}", srv.patchPersona)
	g.DEL("/personas/{persona_id}", srv.deletePersona)

	g.GET("/arbol", srv.getArbolCompleto)
	g.GET("/lista/{nodo_id}", srv.getHistoriasLista)
	g.GET("/tablero/{nodo_id}", srv.getHistoriasTablero)
	g.GET("/prioritarias", srv.getHistoriasPrioritarias)
	g.GET("/historias/{historia_id}", srv.getTareasDeHistoria)

	g.PCH("/historias/{historia_id}", srv.patchHistoria)
	g.DEL("/historias/{historia_id}", srv.deleteHistoria)
	g.GET("/historias/{historia_id}/form", srv.formHistoria)
	g.GET("/historias/{historia_id}/mover", srv.moverHistoriaForm)
	g.POS("/historias/{historia_id}/mover", srv.moverHistoria)
	g.POS("/historias/{historia_id}/priorizar", srv.priorizarHistoria)
	g.POS("/historias/{historia_id}/marcar", srv.marcarHistoria)
	g.POS("/historias/{historia_id}/tareas", srv.postTarea)

	g.POS("/nodos/{nodo_id}", srv.postHistoria)
	g.POS("/nodos/{nodo_id}/reordenar", srv.reordenarNodo)

	g.PCH("/tareas/{tarea_id}", srv.modificarTarea)
	g.POS("/tareas/{tarea_id}/iniciar", srv.iniciarTarea)
	g.POS("/tareas/{tarea_id}/pausar", srv.pausarTarea)
	g.POS("/tareas/{tarea_id}/terminar", srv.terminarTarea)

	g.GET("/intervalos", srv.getIntervalos)

	// LOG SQLITE
	g.GET("/log", func(c *gecko.Context) error { sqliteDB.ToggleLog(); return c.StatusOk("Log toggled") })
	// sqliteDB.ToggleLog()

	// ================================================================ //
	// ================================================================ //

	// Handle interrupt.
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for sig := range ch {
			err = sqliteDB.Close()
			if err != nil {
				fmt.Println("sqliteDB.Close: ", err.Error())
			}
			fmt.Println("")
			gecko.LogInfof("servidor terminado: %v", sig.String())
			os.Exit(0)
		}
	}()

	// Listen and serve
	serv := http.Server{
		Addr:    ":" + puerto,
		Handler: g,
	}
	gecko.LogInfof("pargo escuchando en :%v", puerto)
	err = serv.ListenAndServe()
	if err != nil {
		fatal(err.Error())
	}
}

// ================================================================ //
// ========== UTILS =============================================== //

func fatal(msg string) {
	fmt.Println("[FATAL] " + msg)
	os.Exit(1)
}
