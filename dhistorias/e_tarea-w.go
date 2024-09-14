package dhistorias

import (
	"monorepo/ust"
	"strings"
	"time"

	"github.com/pargomx/gecko/gko"
)

func AgregarTarea(tarea ust.Tarea, repo Repo) error {
	op := gko.Op("AgregarTarea")
	err := validarTarea(&tarea, op, repo)
	if err != nil {
		return err
	}
	tarea.Estatus = 0
	err = repo.InsertTarea(tarea)
	if err != nil {
		return err
	}
	return nil
}

func EliminarTarea(tareaID int, repo Repo) error {
	op := gko.Op("EliminarTarea").Ctx("tareaID", tareaID)
	err := repo.DeleteTarea(tareaID)
	if err != nil {
		return op.Err(err)
	}
	return nil
}

func ActualizarTarea(tareaID int, nueva ust.Tarea, repo Repo) error {
	op := gko.Op("ActualizarTarea")
	tar, err := repo.GetTarea(tareaID)
	if err != nil {
		return op.Err(err)
	}

	if tareaID != nueva.TareaID {
		return op.Msg("El ID de la tarea no se puede cambiar")
	}
	tar.HistoriaID = nueva.HistoriaID
	tar.Tipo = nueva.Tipo
	tar.Descripcion = nueva.Descripcion
	tar.Impedimentos = nueva.Impedimentos
	tar.TiempoEstimado = nueva.TiempoEstimado

	err = actualizarTiempoReal(tar, op, repo)
	if err != nil {
		return err
	}

	err = validarTarea(tar, op, repo)
	if err != nil {
		return err
	}
	err = repo.UpdateTarea(*tar)
	if err != nil {
		return err
	}
	return nil
}

func validarTarea(tarea *ust.Tarea, op *gko.Error, repo Repo) error {
	if tarea.TareaID == 0 {
		return op.Msg("Debe asignarse un ID nuevo a la tarea")
	}
	if tarea.HistoriaID == 0 {
		return op.Msg("La tarea debe pertenecer a una historia")
	}
	if tarea.Descripcion == "" {
		return op.Msg("La tarea debe tener una descripción")
	}
	_, err := repo.GetNodoHistoria(tarea.HistoriaID)
	if err != nil {
		return op.Err(err).Ctx("historiaID", tarea.HistoriaID)
	}
	if strings.HasPrefix(strings.ToLower(tarea.Descripcion), "bug:") {
		gko.LogInfo("Tarea de tipo bug")
		tarea.Tipo = ust.TipoTareaBug
		tarea.Descripcion = strings.TrimSpace(tarea.Descripcion[4:])
	} else {
		gko.LogInfof("Hmm: '%v'", strings.ToLower(tarea.Descripcion))
	}
	return nil
}

// ================================================================ //
// ========== Intervalos ========================================== //

func actualizarTiempoReal(tar *ust.Tarea, op *gko.Error, repo Repo) error {
	intervalos, err := repo.ListIntervalosByTareaID(tar.TareaID)
	if err != nil {
		return op.Err(err)
	}
	tar.TiempoReal = 0
	for _, itv := range intervalos {
		if itv.Fin == "" {
			continue
		}
		inicio, err := time.Parse("2006-01-02 15:04:05", itv.Inicio) // UTC
		if err != nil {
			return op.Err(err).Op("ParseInicio").Ctx("string", itv.Inicio)
		}
		fin, err := time.Parse("2006-01-02 15:04:05", itv.Fin) // UTC
		if err != nil {
			return op.Err(err).Op("ParseFin").Ctx("string", itv.Fin)
		}
		tar.TiempoReal += int(fin.Sub(inicio).Seconds())
	}
	return nil
}

func IniciarTarea(tareaID int, repo Repo) error {
	op := gko.Op("IniciarTarea").Ctx("tareaID", tareaID)
	tar, err := repo.GetTarea(tareaID)
	if err != nil {
		return op.Err(err)
	}

	// No debe haber otro intervalo en curso.
	intervalos, err := repo.ListIntervalosByTareaID(tar.TareaID)
	if err != nil {
		return op.Err(err)
	}
	for _, itv := range intervalos {
		if itv.Fin == "" {
			return op.Msgf("Esta tarea ya se había iniciado en %v", itv.Inicio)
		}
	}

	// Iniciar intervalo.
	interv := ust.Intervalo{
		TareaID: tareaID,
		Inicio:  time.Now().UTC().Format("2006-01-02 15:04:05"),
	}
	err = repo.InsertIntervalo(interv)
	if err != nil {
		return op.Err(err)
	}

	// Declarar tarea en curso.
	tar.Estatus = ust.EstatusTareaEnCurso
	err = repo.UpdateTarea(*tar)
	if err != nil {
		return op.Err(err)
	}
	return nil
}

func PausarTarea(tareaID int, repo Repo) error {
	op := gko.Op("PausarTarea").Ctx("tareaID", tareaID)
	tar, err := repo.GetTarea(tareaID)
	if err != nil {
		return op.Err(err)
	}

	// Debe haber un intervalo en curso.
	intervalos, err := repo.ListIntervalosByTareaID(tar.TareaID)
	if err != nil {
		return op.Err(err)
	}
	var interv *ust.Intervalo
	for _, itv := range intervalos {
		if itv.Fin == "" {
			interv = &itv
			break
		}
	}
	if interv == nil {
		return op.Msg("No hay ningún intervalo en curso para esta tarea")
	}

	// Finalizar intervalo.
	interv.Fin = time.Now().UTC().Format("2006-01-02 15:04:05")
	err = repo.UpdateIntervalo(*interv)
	if err != nil {
		return op.Err(err)
	}

	// Actualizar tiempo real de tarea y declarar como finalizada.
	err = actualizarTiempoReal(tar, op, repo)
	if err != nil {
		return err
	}
	tar.Estatus = ust.EstatusTareaEnPausa
	err = repo.UpdateTarea(*tar)
	if err != nil {
		return op.Err(err)
	}
	return nil
}

func FinalizarTarea(tareaID int, repo Repo) error {
	op := gko.Op("FinalizarTarea").Ctx("tareaID", tareaID)
	tar, err := repo.GetTarea(tareaID)
	if err != nil {
		return op.Err(err)
	}

	// Debe haber un intervalo en curso.
	intervalos, err := repo.ListIntervalosByTareaID(tar.TareaID)
	if err != nil {
		return op.Err(err)
	}
	var interv *ust.Intervalo
	for _, itv := range intervalos {
		if itv.Fin == "" {
			interv = &itv
			break
		}
	}
	if interv == nil {
		return op.Msg("No hay ningún intervalo en curso para esta tarea")
	}

	// Finalizar intervalo.
	interv.Fin = time.Now().UTC().Format("2006-01-02 15:04:05")
	err = repo.UpdateIntervalo(*interv)
	if err != nil {
		return op.Err(err)
	}

	// Actualizar tiempo real de tarea y declarar como finalizada.
	err = actualizarTiempoReal(tar, op, repo)
	if err != nil {
		return err
	}
	tar.Estatus = ust.EstatusTareaFinalizada
	err = repo.UpdateTarea(*tar)
	if err != nil {
		return op.Err(err)
	}
	return nil
}
