package main

import (
	"monorepo/arbol"
	"monorepo/dhistorias"
	"monorepo/ust"

	"github.com/pargomx/gecko"
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

func (s *writehdl) ciclarImportanciaTarea(c *gecko.Context, tx *handlerTx) error {
	args := arbol.ArgsParcharNodo{
		NodoID: c.PathInt("tarea_id"),
		Campo:  "importancia",
		NewVal: "",
	}
	err := tx.app.ParcharNodo(args)
	if err != nil {
		return err
	}
	nod, err := tx.repo.GetNodo(args.NodoID)
	if err != nil {
		return err
	}
	defer s.reloader.brodcastReload(c)
	return c.AskedForFallback("/historias/%v", nod.PadreID)
}

func (s *writehdl) cambiarEstimadoTarea(c *gecko.Context, tx *handlerTx) error {
	args := arbol.ArgsParcharNodo{
		NodoID: c.PathInt("tarea_id"),
		Campo:  "estimado",
		NewVal: c.PromptVal(),
	}
	err := tx.app.ParcharNodo(args)
	if err != nil {
		return err
	}
	nod, err := tx.repo.GetNodo(args.NodoID)
	if err != nil {
		return err
	}
	defer s.reloader.brodcastReload(c)
	return c.AskedForFallback("/historias/%v#%v", nod.PadreID, nod.NodoID)
}

func (s *writehdl) iniciarTarea(c *gecko.Context, tx *handlerTx) error {
	tareaID := c.PathInt("tarea_id")
	err := tx.app.IniciarTarea(tareaID)
	if err != nil {
		return err
	}
	tar, err := tx.repo.GetNodo(tareaID)
	if err != nil {
		return err
	}
	defer s.reloader.brodcastReload(c)
	return c.AskedForFallback("/historias/%v", tar.PadreID)
}

func (s *writehdl) pausarTarea(c *gecko.Context, tx *handlerTx) error {
	tareaID := c.PathInt("tarea_id")
	err := tx.app.PausarTarea(tareaID)
	if err != nil {
		return err
	}
	tar, err := tx.repo.GetNodo(tareaID)
	if err != nil {
		return err
	}
	defer s.reloader.brodcastReload(c)
	return c.AskedForFallback("/historias/%v", tar.PadreID)
}

func (s *writehdl) terminarTarea(c *gecko.Context, tx *handlerTx) error {
	tareaID := c.PathInt("tarea_id")
	err := tx.app.FinalizarTarea(tareaID)
	if err != nil {
		return err
	}
	tar, err := tx.repo.GetNodo(tareaID)
	if err != nil {
		return err
	}
	defer s.reloader.brodcastReload(c)
	return c.AskedForFallback("/historias/%v", tar.PadreID)
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
