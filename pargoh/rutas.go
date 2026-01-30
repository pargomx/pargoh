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

	// Imágenes de usuario
	s.gecko.StaticSub("/imagenes", s.cfg.ImagesDir)

	// Sesiones
	s.gecko.GET("/", s.auth.getLogin)
	s.gecko.GET("/login", s.auth.getLogin)
	s.gecko.POS("/login", s.auth.postLogin)
	s.GET("/logout", s.auth.logout)
	s.GET("/sesiones", s.auth.printSesiones)

	// General
	s.GET("/buscar", s.r.buscar)

	s.GET("/continuar", s.continuar)
	s.GET("/offline", s.offline)

	// Proyectos
	s.GET("/proyectos", s.r.listaProyectos)
	s.GET("/proyectos/{proyecto_id}", s.r.getProyecto)
	s.GET("/proyectos/{proyecto_id}/doc", s.r.getDocumentacionProyecto)

	// Personas
	s.GET("/personas/{persona_id}", s.r.getPersona)
	s.GET("/personas/{persona_id}/doc", s.r.getPersonaDoc)
	s.GET("/personas/{persona_id}/debug", s.r.getPersonaDebug)

	// Historias
	s.GET("/historias/{historia_id}", s.r.getHistoria)
	s.GET("/historias/{historia_id}/tablero", s.r.getHistoriaTablero)

	// Tareas técnicas
	s.GET("/tareas/{tarea_id}", s.r.getTarea)
	s.GET("/intervalos", s.r.getIntervalos)

	// Quick tasks
	s.GET("/tareas", s.r.getQuickTasks)

	// Navegador del árbol de historias
	s.GET("/nav", s.r.navDesdeRoot)
	s.GET("/nav/proy/{proyecto_id}", s.r.navDesdeProyecto)
	s.GET("/nav/pers/{persona_id}", s.r.navDesdePersona)
	s.GET("/nav/hist/{historia_id}", s.r.navDesdeHistoria)
	s.GET("/historias/{historia_id}/mover", s.r.moverHistoriaForm)

	// MOVER
	s.POS("/historias/{historia_id}/mover", s.w.moverHistoria)
	s.POS("/mover/tramo", s.w.moverTramo)
	s.POS("/mover/tarea", s.w.moverTarea)
	s.POS("/mover/historia", s.w.moverHistoria)

	// AGREGAR HOJA
	s.POS("/proyectos", s.w.postProyecto)
	s.POS("/personas", s.w.postPersona)
	s.POS("/personas/{persona_id}", s.w.postHistoriaDePersona)
	s.POS("/historias/{historia_id}", s.w.postHistoriaDeHistoria)
	s.POS("/historias/{historia_id}/padre", s.w.postPadreParaHistoria)
	s.POS("/historias/{historia_id}/reglas", s.w.postRegla)
	s.POS("/historias/{historia_id}/viaje", s.w.postTramoDeViaje)
	s.POS("/historias/{historia_id}/tareas", s.w.postTarea)
	s.POS("/tareas", s.w.postQuickTask)

	// REORDENAR
	s.POS("/reordenar-persona", s.w.reordenarPersona)
	s.POS("/reordenar-historia", s.w.reordenarHistoria)
	s.POS("/reordenar-tramo", s.w.reordenarTramo)
	s.POS("/reordenar-regla", s.w.reordenarRegla)

	// PARCHAR
	s.PCH("/proyectos/{proyecto_id}/{param}", s.w.patchProyecto)
	s.PCH("/personas/{persona_id}/{param}", s.w.patchPersona)
	s.PCH("/historias/{historia_id}/{param}", s.w.patchHistoria)
	s.PCH("/historias/{historia_id}/reglas/{posicion}", s.w.patchRegla)
	s.PCH("/historias/{historia_id}/viaje/{posicion}", s.w.patchTramoDeViaje)

	// OTROS
	s.PCH("/historias/{historia_id}/reglas/{posicion}/marcar", s.w.marcarRegla)
	s.POS("/historias/{historia_id}/priorizar", s.w.priorizarHistoria)
	s.POS("/historias/{historia_id}/priorizar/{prioridad}", s.w.priorizarHistoria)
	s.POS("/historias/{historia_id}/marcar", s.w.marcarHistoria)
	s.POS("/historias/{historia_id}/marcar/{completada}", s.w.marcarHistoria)

	s.PCH("/tareas/{tarea_id}", s.modificarTarea)
	s.PCH("/tareas/{tarea_id}/estimado", s.w.cambiarEstimadoTarea)
	s.POS("/tareas/{tarea_id}/importancia", s.w.ciclarImportanciaTarea)
	s.POS("/tareas/{tarea_id}/iniciar", s.w.iniciarTarea)
	s.POS("/tareas/{tarea_id}/pausar", s.w.pausarTarea)
	s.POS("/tareas/{tarea_id}/terminar", s.w.terminarTarea)
	s.PCH("/tareas/{tarea_id}/intervalos/{ts_id}/{cambiar}", s.w.patchIntervalo)

	// ELIMINAR
	s.DEL("/proyectos/{proyecto_id}", s.w.deleteProyecto)
	s.DEL("/proyectos/{proyecto_id}/definitivo", s.w.deleteProyectoPorCompleto)
	s.DEL("/personas/{persona_id}", s.w.deletePersona)
	s.DEL("/historias/{historia_id}", s.w.deleteHistoria)
	s.DEL("/tareas/{tarea_id}", s.w.eliminarTarea)
	s.DEL("/historias/{historia_id}/reglas/{posicion}", s.w.deleteRegla)
	s.DEL("/historias/{historia_id}/viaje/{posicion}", s.w.deleteTramoDeViaje)

	// Raw nodo
	s.GET("/nodos/{nodo_id}", s.r.getRawNodoEditor)
	s.PCH("/nodos/{nodo_id}/{param}", s.w.patchRawNodo)

	// TIME TRACKER
	s.POS("/nodos/{nodo_id}/time/{seg}", s.w.postAppTime)

	// IMAGENES
	s.POS("/imagenes", s.setImagenTramo)
	s.PUT("/proyectos/{proyecto_id}", s.setImagenProyecto)
	s.DEL("/imagenes/{historia_id}/{posicion}", s.deleteImagenTramo)

	// Referencias
	s.POS("/historias/{nodo_id}/referencias", s.w.postReferencia)
	s.DEL("/historias/{nodo_id}/referencias/{ref_nodo_id}", s.w.deleteReferencia)

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

	s.GET("/reload", s.brodcastReload)
	s.GET("/historias/{historia_id}/ws", s.reloader.nuevoWS)

	// Mantenimiento
	s.GET("/log", func(c *gecko.Context) error { s.db.ToggleLog(); return c.StatusOk("Log toggled") })
	s.GET("/clear", func(c *gecko.Context) error {
		c.Response().Header().Set("Clear-Site-Data", `"cache", "cache", "clientHints", "storage", "executionContexts"`)
		return c.StringOk("Datos del sitio limpiados. Ok.")
	})
}
