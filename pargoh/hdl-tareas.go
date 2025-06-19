package main

import (
	"monorepo/arbol"
	"monorepo/dhistorias"
	"monorepo/ust"

	"github.com/pargomx/gecko"
	"github.com/pargomx/gecko/gko"
)

func (s *writehdl) postTarea(c *gecko.Context, tx *handlerTx) error {
	args := arbol.ArgsAgregarTarea{
		Tipo:     "TAR",
		NodoID:   ust.NewRandomID(),
		PadreID:  c.PathInt("historia_id"),
		Titulo:   c.FormVal("descripcion"),
		Estimado: c.FormVal("segundos_estimado"),
	}
	err := tx.app.AgregarTarea(args)
	if err != nil {
		return err
	}
	defer s.reloader.brodcastReload(c)
	return c.AskedForFallback("/historias/%v", args.PadreID)
}

func (s *servidor) modificarTarea(c *gecko.Context) error {
	estimado, err := ust.NuevaDuraciónSegundos(c.FormVal("segundos_estimado"))
	if err != nil {
		return err
	}
	tarea := ust.Tarea{
		TareaID:          c.FormInt("tarea_id"),
		HistoriaID:       c.FormInt("historia_id"),
		Tipo:             ust.SetTipoTareaDB(c.FormVal("tipo")),
		Descripcion:      c.FormVal("descripcion"),
		Impedimentos:     c.FormVal("impedimentos"),
		SegundosEstimado: estimado,
		Importancia:      ust.SetImportanciaTareaDB(c.FormVal("importancia")),
	}
	err = dhistorias.ActualizarTarea(c.PathInt("tarea_id"), tarea, s.repoOld)
	if err != nil {
		return err
	}
	defer s.reloader.brodcastReload(c)
	return c.AskedForFallback("/historias/%v#%v", tarea.HistoriaID, tarea.TareaID)
}

func (s *servidor) ciclarImportanciaTarea(c *gecko.Context) error {
	tarea, err := s.repoOld.GetTarea(c.PathInt("tarea_id"))
	if err != nil {
		return err
	}
	if tarea.Importancia.EsIdea() {
		tarea.Importancia = ust.ImportanciaTareaNecesaria
	} else if tarea.Importancia.EsNecesaria() {
		tarea.Importancia = ust.ImportanciaTareaMejora
	} else if tarea.Importancia.EsMejora() {
		tarea.Importancia = ust.ImportanciaTareaIdea
	} else {
		tarea.Importancia = ust.ImportanciaTareaIdea
	}
	err = dhistorias.ActualizarTarea(tarea.TareaID, *tarea, s.repoOld)
	if err != nil {
		return err
	}
	defer s.reloader.brodcastReload(c)
	return c.AskedForFallback("/historias/%v", tarea.HistoriaID)
}

func (s *servidor) cambiarEstimadoTarea(c *gecko.Context) error {
	estimado, err := ust.NuevaDuraciónSegundos(c.PromptVal())
	if err != nil {
		return err
	}
	if estimado <= 0 {
		return gko.ErrDatoInvalido.Msg("El estimado debe ser mayor a 0")
	}
	tarea, err := s.repoOld.GetTarea(c.PathInt("tarea_id"))
	if err != nil {
		return err
	}
	tarea.SegundosEstimado = estimado
	err = dhistorias.ActualizarTarea(c.PathInt("tarea_id"), *tarea, s.repoOld)
	if err != nil {
		return err
	}
	defer s.reloader.brodcastReload(c)
	return c.AskedForFallback("/historias/%v#%v", tarea.HistoriaID, tarea.TareaID)
}

func (s *servidor) eliminarTarea(c *gecko.Context) error {
	historiaID, err := dhistorias.EliminarTarea(c.PathInt("tarea_id"), s.repoOld)
	if err != nil {
		return err
	}
	gko.LogInfof("Tarea %d eliminada", c.PathInt("tarea_id"))
	defer s.reloader.brodcastReload(c)
	return c.AskedForFallback("/historias/%v", historiaID)
}

func (s *servidor) iniciarTarea(c *gecko.Context) error {
	historiaID, err := dhistorias.IniciarTarea(c.PathInt("tarea_id"), s.repoOld)
	if err != nil {
		return err
	}
	defer s.reloader.brodcastReload(c)
	return c.AskedForFallback("/historias/%v", historiaID)
}
func (s *servidor) pausarTarea(c *gecko.Context) error {
	historiaID, err := dhistorias.PausarTarea(c.PathInt("tarea_id"), s.repoOld)
	if err != nil {
		return err
	}
	defer s.reloader.brodcastReload(c)
	return c.AskedForFallback("/historias/%v", historiaID)
}
func (s *servidor) terminarTarea(c *gecko.Context) error {
	historiaID, err := dhistorias.FinalizarTarea(c.PathInt("tarea_id"), s.repoOld)
	if err != nil {
		return err
	}
	defer s.reloader.brodcastReload(c)
	return c.AskedForFallback("/historias/%v", historiaID)
}

/*
func (s *servidor) materializarTiemposTareas(c *gecko.Context) error {
	err := dhistorias.MaterializarTiempoRealTareas(s.repo)
	if err != nil {
		return err
	}
	return c.StringOk("Tiempos actualizados según los intervalos de trabajo")
}
*/

func (s *readhdl) getTarea(c *gecko.Context) error {
	tarea, err := s.repoOld.GetTarea(c.PathInt("tarea_id"))
	if err != nil {
		return err
	}
	intervalos, err := s.repoOld.ListIntervalosByTareaID(tarea.TareaID)
	if err != nil {
		return err
	}
	data := map[string]any{
		"Tarea":      tarea,
		"Intervalos": intervalos,
	}
	return c.RenderOk("tarea", data)
}

func (s *readhdl) getIntervalos(c *gecko.Context) error {
	recientes, err := s.repoOld.ListIntervalosRecientes()
	if err != nil {
		return err
	}
	abiertos, err := s.repoOld.ListIntervalosRecientesAbiertos()
	if err != nil {
		return err
	}
	data := map[string]any{
		"Titulo":    "Intervalos de trabajo",
		"Recientes": recientes,
		"Abiertos":  abiertos,
	}
	return c.RenderOk("intervalos", data)
}

func (s *servidor) patchIntervalo(c *gecko.Context) error {
	historiaID, err := dhistorias.ParcharIntervalo(c.PathInt("tarea_id"), c.PathVal("inicio"), c.FormVal("inicio"), c.FormVal("fin"), s.repoOld)
	if err != nil {
		return err
	}
	// defer s.reloader.brodcastReload(c)
	return c.AskedForFallback("/historias/%v", historiaID)
}

// ================================================================ //
// ========== QUICK TASKS ========================================= //

func (s *writehdl) postQuickTask(c *gecko.Context, tx *handlerTx) error {
	args := arbol.ArgsAgregarTarea{
		Tipo:     "TAR",
		NodoID:   ust.NewRandomID(),
		PadreID:  0, // TODO: poner default historia en migración?
		Titulo:   c.PromptVal(),
		Estimado: "1h",
	}
	err := tx.app.AgregarTarea(args)
	if err != nil {
		return err
	}
	// _, err = dhistorias.IniciarTarea(tarea.TareaID, s.repoOld)
	// if err != nil {
	// 	return err
	// }
	defer s.reloader.brodcastReload(c)
	return c.AskedForFallback("/tareas")
}

func (s *readhdl) getQuickTasks(c *gecko.Context) error {
	tareas, err := s.repoOld.ListTareasByHistoriaID(dhistorias.QUICK_TASK_HISTORIA_ID)
	if err != nil {
		return err
	}
	data := map[string]any{
		"Titulo": "Quick tasks",
		"Tareas": tareas,
	}
	return c.RenderOk("tareas", data)
}
