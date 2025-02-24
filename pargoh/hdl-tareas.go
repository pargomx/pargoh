package main

import (
	"monorepo/dhistorias"
	"monorepo/ust"

	"github.com/pargomx/gecko"
	"github.com/pargomx/gecko/gko"
)

func (s *servidor) postTarea(c *gecko.Context) error {
	estimado, err := ust.NuevaDuraciónSegundos(c.FormVal("segundos_estimado"))
	if err != nil {
		return err
	}
	tarea := ust.Tarea{
		TareaID:          ust.NewRandomID(),
		HistoriaID:       c.PathInt("historia_id"),
		Tipo:             ust.SetTipoTareaDB(c.FormVal("tipo")),
		Descripcion:      c.FormVal("descripcion"),
		Impedimentos:     c.FormVal("impedimentos"),
		SegundosEstimado: estimado,
		Importancia:      ust.SetImportanciaTareaDB(c.FormVal("importancia")),
	}
	err = dhistorias.AgregarTarea(tarea, s.repo)
	if err != nil {
		return err
	}
	defer s.reloader.brodcastReload(c)
	return c.AskedForFallback("/historias/%v", tarea.HistoriaID)
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
	err = dhistorias.ActualizarTarea(c.PathInt("tarea_id"), tarea, s.repo)
	if err != nil {
		return err
	}
	defer s.reloader.brodcastReload(c)
	return c.AskedForFallback("/historias/%v#%v", tarea.HistoriaID, tarea.TareaID)
}

func (s *servidor) ciclarImportanciaTarea(c *gecko.Context) error {
	tarea, err := s.repo.GetTarea(c.PathInt("tarea_id"))
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
	err = dhistorias.ActualizarTarea(tarea.TareaID, *tarea, s.repo)
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
		return gko.ErrDatoInvalido().Msg("El estimado debe ser mayor a 0")
	}
	tarea, err := s.repo.GetTarea(c.PathInt("tarea_id"))
	if err != nil {
		return err
	}
	tarea.SegundosEstimado = estimado
	err = dhistorias.ActualizarTarea(c.PathInt("tarea_id"), *tarea, s.repo)
	if err != nil {
		return err
	}
	defer s.reloader.brodcastReload(c)
	return c.AskedForFallback("/historias/%v#%v", tarea.HistoriaID, tarea.TareaID)
}

func (s *servidor) moverTarea(c *gecko.Context) error {
	historiaID, err := dhistorias.MoverTarea(c.FormInt("tarea_id"), c.FormInt("target_historia_id"), s.repo)
	if err != nil {
		return err
	}
	defer s.reloader.brodcastReload(c)
	return c.AskedForFallback("/historias/%v", historiaID)
}

func (s *servidor) eliminarTarea(c *gecko.Context) error {
	historiaID, err := dhistorias.EliminarTarea(c.PathInt("tarea_id"), s.repo)
	if err != nil {
		return err
	}
	gko.LogInfof("Tarea %d eliminada", c.PathInt("tarea_id"))
	defer s.reloader.brodcastReload(c)
	return c.AskedForFallback("/historias/%v", historiaID)
}

func (s *servidor) iniciarTarea(c *gecko.Context) error {
	historiaID, err := dhistorias.IniciarTarea(c.PathInt("tarea_id"), s.repo)
	if err != nil {
		return err
	}
	defer s.reloader.brodcastReload(c)

	return c.AskedForFallback("/historias/%v", historiaID)
}
func (s *servidor) pausarTarea(c *gecko.Context) error {
	historiaID, err := dhistorias.PausarTarea(c.PathInt("tarea_id"), s.repo)
	if err != nil {
		return err
	}
	defer s.reloader.brodcastReload(c)
	return c.AskedForFallback("/historias/%v", historiaID)
}
func (s *servidor) terminarTarea(c *gecko.Context) error {
	historiaID, err := dhistorias.FinalizarTarea(c.PathInt("tarea_id"), s.repo)
	if err != nil {
		return err
	}
	defer s.reloader.brodcastReload(c)
	return c.AskedForFallback("/historias/%v", historiaID)
}

func (s *servidor) materializarTiemposTareas(c *gecko.Context) error {
	err := dhistorias.MaterializarTiempoRealTareas(s.repo)
	if err != nil {
		return err
	}
	return c.StringOk("Tiempos actualizados según los intervalos de trabajo")
}

func (s *servidor) getTarea(c *gecko.Context) error {
	tarea, err := s.repo.GetTarea(c.PathInt("tarea_id"))
	if err != nil {
		return err
	}
	intervalos, err := s.repo.ListIntervalosByTareaID(tarea.TareaID)
	if err != nil {
		return err
	}
	data := map[string]any{
		"Tarea":      tarea,
		"Intervalos": intervalos,
	}
	return c.RenderOk("tarea", data)
}

func (s *servidor) getIntervalos(c *gecko.Context) error {
	recientes, err := s.repo.ListIntervalosRecientes()
	if err != nil {
		return err
	}
	abiertos, err := s.repo.ListIntervalosRecientesAbiertos()
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
	historiaID, err := dhistorias.ParcharIntervalo(c.PathInt("tarea_id"), c.PathVal("inicio"), c.FormVal("inicio"), c.FormVal("fin"), s.repo)
	if err != nil {
		return err
	}
	// defer s.reloader.brodcastReload(c)
	return c.AskedForFallback("/historias/%v", historiaID)
}

// ================================================================ //
// ========== QUICK TASKS ========================================= //

// Tareas sin proyecto.
//
//	INSERT INTO historias(historia_id, titulo, objetivo, prioridad, completada) VALUES (1,'QuickTasksParent','Esta historia sirve de padre para las tareas sin proyecto',0,0);
const QUICK_TASK_HISTORIA_ID = 0001

func (s *servidor) postQuickTask(c *gecko.Context) error {
	tarea := ust.Tarea{
		TareaID:     ust.NewRandomID(),
		HistoriaID:  QUICK_TASK_HISTORIA_ID,
		Descripcion: c.PromptVal(),
	}
	err := dhistorias.AgregarTarea(tarea, s.repo)
	if err != nil {
		return err
	}
	_, err = dhistorias.IniciarTarea(tarea.TareaID, s.repo)
	if err != nil {
		return err
	}
	return c.AskedForFallback("/tareas")
}

func (s *servidor) getQuickTasks(c *gecko.Context) error {
	tareas, err := s.repo.ListTareasByHistoriaID(QUICK_TASK_HISTORIA_ID)
	if err != nil {
		return err
	}
	data := map[string]any{
		"Titulo": "Quick tasks",
		"Tareas": tareas,
	}
	return c.RenderOk("tareas", data)
}
