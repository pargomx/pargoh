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
	s.GET("/metricas", s.r.getMétricas)

	s.GET("/continuar", s.continuar)
	s.GET("/offline", s.offline)
	s.GET("/reload", s.brodcastReload)

	// Ver el árbol
	s.GET("/h", s.r.getProyectosActivos)
	s.GET("/h/{nodo_id}", s.r.getNodoCualquiera)
	s.GET("/h/{nodo_id}/raw", s.r.getRawNodoEditor)
	s.GET("/h/{nodo_id}/doc", s.r.getDocumentacion)
	s.GET("/h/{nodo_id}/tablero", s.r.getNodoTablero)

	// Agregar al árbol
	s.POS("/h/{nodo_id}", s.w.postNodo)
	s.POS("/h/{nodo_id}/padre", s.w.postNodoPadre)
	s.POS("/h/{nodo_id}/{tipo}", s.w.postNodoDeTipo)
	s.POS("/h/{nodo_id}/tareas", s.w.postTarea)

	// Eliminar del árbol
	s.DEL("/h/{nodo_id}", s.w.eliminarNodo)
	s.DEL("/h/{nodo_id}/definitivo", s.w.eliminarRama)

	// Mover el árbol
	s.GET("/nav", s.r.navDesdeRoot)
	s.GET("/nav/{nodo_id}", s.r.navDesdeNodo)
	s.POS("/mover", s.w.moverNodo)
	s.POS("/reordenar", s.w.reordenarNodo)
	// Referencias
	s.POS("/h/{nodo_id}/referencias", s.w.postReferencia)
	s.DEL("/h/{nodo_id}/referencias/{ref_nodo_id}", s.w.deleteReferencia)

	// Modificar nodos
	s.PCH("/h/{nodo_id}", s.modificarTarea)
	s.PCH("/h/{nodo_id}/{param}", s.w.parcharNodo)
	s.PCH("/h/{nodo_id}/estimado", s.w.cambiarEstimadoPrompt)
	s.POS("/h/{nodo_id}/priorizar", s.w.priorizarHistoria)
	s.POS("/h/{nodo_id}/priorizar/{prioridad}", s.w.priorizarHistoria)
	s.POS("/h/{nodo_id}/marcar", s.w.marcarHistoria)
	s.POS("/h/{nodo_id}/marcar/{completada}", s.w.marcarHistoria)

	// Time tracker
	s.GET("/intervalos", s.r.getIntervalos)
	s.GET("/h/{nodo_id}/ws", s.reloader.nuevoWS)
	s.POS("/h/{nodo_id}/time/{seg}", s.w.postAppTime)
	s.PCH("/h/{nodo_id}/intervalos/{ts_id}/{cambiar}", s.w.patchIntervalo)

	// Quick tasks
	s.GET("/tareas", s.r.getQuickTasks)
	s.POS("/tareas", s.w.postQuickTask)

	// Imágenes
	s.gecko.StaticSub("/imagenes", s.cfg.ImagesDir)
	s.POS("/imagenes", s.setImagenTramo)
	s.PUT("/h/{proyecto_id}/imagen", s.setImagenProyecto)
	s.DEL("/imagenes/{historia_id}/{posicion}", s.eliminarImagen)

	// Mantenimiento
	s.GET("/log", func(c *gecko.Context) error { s.db.ToggleLog(); return c.StatusOk("Log toggled") })
	s.GET("/clear", func(c *gecko.Context) error {
		c.Response().Header().Set("Clear-Site-Data", `"cache", "cache", "clientHints", "storage", "executionContexts"`)
		return c.StringOk("Datos del sitio limpiados. Ok.")
	})
}
