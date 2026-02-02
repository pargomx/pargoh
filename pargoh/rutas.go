package main

import (
	"monorepo/assets"

	"github.com/pargomx/gecko"
)

func (s *servidor) registrarRutas() {
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

	// Auth
	s.gecko.GET("/", s.auth.getLogin)
	s.gecko.GET("/login", s.auth.getLogin)
	s.gecko.POS("/login", s.auth.postLogin)
	s.GET("/logout", s.auth.logout)
	s.GET("/sesiones", s.auth.printSesiones)

	// General
	s.GET("/buscar", s.r.buscar)

	s.GET("/continuar", s.continuar)
	s.GET("/offline", s.offline)

	s.POS("/h", s.w.postNodo)
	s.GET("/h", s.r.getProyectosActivos)
	s.GET("/h/{nodo_id}", s.r.getNodoCualquiera)
	s.GET("/h/{nodo_id}/tablero", s.r.getHistoriaTablero)
	s.GET("/h/{nodo_id}/raw", s.r.getRawNodoEditor)
	s.GET("/h/{nodo_id}/doc", s.r.getDocumentacionProyecto)

	s.PCH("/h/{nodo_id}/{param}", s.w.patchRawNodo)
	// s.PCH("/h/{nodo_id}/{param}", s.w.patchHistoria)

	s.POS("/h/{nodo_id}/tareas", s.w.postTarea)
	s.GET("/h/{nodo_id}/mover", s.r.moverHistoriaForm)

	s.POS("/h/{nodo_id}", s.w.postHistoriaDeHistoria)
	s.POS("/h/{nodo_id}/padre", s.w.postPadreParaHistoria)
	s.POS("/h/{nodo_id}/reglas", s.w.postRegla)
	s.POS("/h/{nodo_id}/viaje", s.w.postTramoDeViaje)

	// TIME TRACKER
	s.POS("/h/{nodo_id}/time/{seg}", s.w.postAppTime)
	s.GET("/h/{nodo_id}/ws", s.reloader.nuevoWS)
	s.GET("/reload", s.brodcastReload)

	s.GET("/intervalos", s.r.getIntervalos)

	// Quick tasks
	s.GET("/tareas", s.r.getQuickTasks)
	s.POS("/tareas", s.w.postQuickTask)

	// Navegador del árbol de historias
	s.GET("/nav", s.r.navDesdeRoot)
	s.GET("/nav/proy/{proyecto_id}", s.r.navDesdeProyecto)
	s.GET("/nav/pers/{persona_id}", s.r.navDesdePersona)
	s.GET("/nav/hist/{historia_id}", s.r.navDesdeHistoria)

	// MOVER
	s.POS("/historias/{historia_id}/mover", s.w.moverHistoria)
	s.POS("/mover/tramo", s.w.moverTramo)
	s.POS("/mover/tarea", s.w.moverTarea)
	s.POS("/mover/historia", s.w.moverHistoria)

	// AGREGAR HOJA
	s.POS("/proyectos", s.w.postProyecto)
	s.POS("/personas", s.w.postPersona)
	s.POS("/personas/{persona_id}", s.w.postHistoriaDePersona)

	// REORDENAR
	s.POS("/reordenar-persona", s.w.reordenarPersona)
	s.POS("/reordenar-historia", s.w.reordenarHistoria)
	s.POS("/reordenar-tramo", s.w.reordenarTramo)
	s.POS("/reordenar-regla", s.w.reordenarRegla)

	// PARCHAR
	s.PCH("/proyectos/{proyecto_id}/{param}", s.w.patchProyecto)
	s.PCH("/personas/{persona_id}/{param}", s.w.patchPersona)
	s.PCH("/h/{nodo_id}/reglas/{posicion}", s.w.patchRegla)
	s.PCH("/h/{nodo_id}/viaje/{posicion}", s.w.patchTramoDeViaje)

	s.PCH("/h/{nodo_id}/reglas/{posicion}/marcar", s.w.marcarRegla)
	s.POS("/h/{nodo_id}/priorizar", s.w.priorizarHistoria)
	s.POS("/h/{nodo_id}/priorizar/{prioridad}", s.w.priorizarHistoria)
	s.POS("/h/{nodo_id}/marcar", s.w.marcarHistoria)
	s.POS("/h/{nodo_id}/marcar/{completada}", s.w.marcarHistoria)

	s.PCH("/tareas/{tarea_id}", s.modificarTarea)
	s.PCH("/h/{nodo_id}/estimado", s.w.cambiarEstimadoTarea)
	s.POS("/h/{nodo_id}/importancia", s.w.ciclarImportanciaTarea)
	s.POS("/h/{nodo_id}/iniciar", s.w.iniciarTarea)
	s.POS("/h/{nodo_id}/pausar", s.w.pausarTarea)
	s.POS("/h/{nodo_id}/terminar", s.w.terminarTarea)
	s.PCH("/h/{nodo_id}/intervalos/{ts_id}/{cambiar}", s.w.patchIntervalo)

	// ELIMINAR
	s.DEL("/h/{nodo_id}/definitivo", s.w.eliminarRama)
	s.DEL("/h/{nodo_id}", s.w.eliminarNodo)

	// IMAGENES
	s.gecko.StaticSub("/imagenes", s.cfg.ImagesDir)
	s.POS("/imagenes", s.setImagenTramo)
	s.PUT("/proyectos/{proyecto_id}", s.setImagenProyecto)
	s.DEL("/imagenes/{historia_id}/{posicion}", s.eliminarImagen)

	// Referencias
	s.POS("/h/{nodo_id}/referencias", s.w.postReferencia)
	s.DEL("/h/{nodo_id}/referencias/{ref_nodo_id}", s.w.deleteReferencia)

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
	s.GET("/metricas", s.r.getMétricas)
	// s.GET("/metricas2", s.getMétricas2)
	// s.GET("/materializar-tiempos", s.materializarTiemposTareas)
	// s.GET("/materializar-historias", s.materializarHistorias)

	// Mantenimiento
	s.GET("/log", func(c *gecko.Context) error { s.db.ToggleLog(); return c.StatusOk("Log toggled") })
	s.GET("/clear", func(c *gecko.Context) error {
		c.Response().Header().Set("Clear-Site-Data", `"cache", "cache", "clientHints", "storage", "executionContexts"`)
		return c.StringOk("Datos del sitio limpiados. Ok.")
	})
}
