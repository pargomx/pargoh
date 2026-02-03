package main

import (
	"monorepo/arbol"
	"monorepo/dhistorias"
	"monorepo/ust"

	"github.com/pargomx/gecko"
)

func (s *servidor) modificarTarea(c *gecko.Context, tx *handlerTx) error {
	estimado, err := ust.NuevaDuraciónSegundos(c.FormVal("segundos_estimado"))
	if err != nil {
		return err
	}
	tarea := ust.Tarea{
		TareaID:          c.FormInt("nodo_id"),
		HistoriaID:       c.FormInt("historia_id"),
		Tipo:             ust.SetTipoTareaDB(c.FormVal("tipo")),
		Descripcion:      c.FormVal("descripcion"),
		Impedimentos:     c.FormVal("impedimentos"),
		SegundosEstimado: estimado,
		Importancia:      ust.SetImportanciaTareaDB(c.FormVal("importancia")),
	}
	err = dhistorias.ActualizarTarea(c.PathInt("nodo_id"), tarea, s.repoOld)
	if err != nil {
		return err
	}
	defer s.reloader.brodcastReload(c)
	return c.AskedForFallback("/h/%v#%v", tarea.HistoriaID, tarea.TareaID)
}

func (s *writehdl) cambiarEstimadoPrompt(c *gecko.Context, tx *handlerTx) error {
	args := arbol.ArgsParcharNodo{
		NodoID: c.PathInt("nodo_id"),
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
	return c.AskedForFallback("/h/%v#%v", nod.PadreID, nod.NodoID)
}

// ================================================================ //
// ========== Intervalos ========================================== //

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

func (s *writehdl) patchIntervalo(c *gecko.Context, tx *handlerTx) error {
	args := arbol.ArgsParcharIntervalo{
		NodoID: c.PathInt("nodo_id"),
		TsID:   c.PathVal("ts_id"),
	}
	if c.PathVal("cambiar") == "ini" {
		args.NewTS = c.FormVal("ts_ini")
		args.Cambiar = "INI"
	} else {
		args.NewTS = c.FormVal("ts_fin")
		args.Cambiar = "FIN"
	}
	err := tx.app.ParcharIntervalo(args)
	if err != nil {
		return err
	}
	padre, err := tx.repo.GetNodo(args.NodoID)
	if err != nil {
		return err
	}
	// defer s.reloader.brodcastReload(c)
	return c.AskedForFallback("/h/%v", padre.NodoID)
}

// ================================================================ //
// ========== QUICK TASKS ========================================= //

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
