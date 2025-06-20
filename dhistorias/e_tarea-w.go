package dhistorias

import (
	"monorepo/ust"
	"strings"
	"time"

	"github.com/pargomx/gecko/gko"
	"github.com/pargomx/gecko/gkt"
)

func AgregarTarea(tarea ust.Tarea, repo Repo) error {
	op := gko.Op("AgregarTarea")
	err := validarTarea(&tarea, op, repo)
	if err != nil {
		return err
	}
	tarea.Estatus = 0
	if tarea.Importancia.EsIndefinido() {
		tarea.Importancia = ust.ImportanciaTareaNecesaria
	}
	err = repo.InsertTarea(tarea)
	if err != nil {
		return err
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
	tar.SegundosEstimado = nueva.SegundosEstimado
	tar.Importancia = nueva.Importancia
	if tar.Importancia.EsIndefinido() {
		tar.Importancia = ust.ImportanciaTareaIdea
	}
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
	err := repo.ExisteHistoria(tarea.HistoriaID)
	if err != nil {
		return op.Err(err).Ctx("historiaID", tarea.HistoriaID)
	}
	// Helpers para descripción.
	comparableDesc := strings.ToLower(tarea.Descripcion)
	if strings.HasPrefix(comparableDesc, "bug:") {
		tarea.Tipo = ust.TipoTareaBug
		tarea.Descripcion = strings.TrimSpace(tarea.Descripcion[4:])
	}
	if strings.HasPrefix(comparableDesc, "ui:") {
		tarea.Tipo = ust.TipoTareaWebUi
		tarea.Descripcion = strings.TrimSpace(tarea.Descripcion[3:])
	}
	if strings.HasPrefix(comparableDesc, "db:") {
		tarea.Tipo = ust.TipoTareaDb
		tarea.Descripcion = strings.TrimSpace(tarea.Descripcion[3:])
	}
	if strings.HasPrefix(comparableDesc, "idea:") {
		tarea.Importancia = ust.ImportanciaTareaIdea
		tarea.Descripcion = strings.TrimSpace(tarea.Descripcion[5:])
	}
	if strings.HasPrefix(comparableDesc, "mejora:") {
		tarea.Importancia = ust.ImportanciaTareaMejora
		tarea.Descripcion = strings.TrimSpace(tarea.Descripcion[7:])
	}
	if len(tarea.Descripcion) < 3 {
		return op.Msg("La descripción de la tarea debe tener al menos 3 caracteres")
	}
	// Primera letra en mayúscula.
	tarea.Descripcion = gkt.PrimeraMayusc(tarea.Descripcion)
	return nil
}

// ================================================================ //
// ========== Intervalos ========================================== //

// TODO: revisar error y quizá hacer configurable. También en pargoh/hdl-metricas.go.
var locationMexicoCity, _ = time.LoadLocation("America/Mexico_City")

func actualizarTiempoReal(tar *ust.Tarea, op *gko.Error, repo Repo) error {
	intervalos, err := repo.ListIntervalosByTareaID(tar.TareaID)
	if err != nil {
		return op.Err(err)
	}
	tar.SegundosUtilizado = 0
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
		tar.SegundosUtilizado += int(fin.Sub(inicio).Seconds())
	}
	return nil
}

/*
func MaterializarTiempoRealTareas(repo Repo) error {
	tareas, err := repo.ListTareas()
	if err != nil {
		return err
	}
	for _, tar := range tareas {
		err = actualizarTiempoReal(&tar, gko.Op("MaterializarTiempoRealTareas"), repo)
		if err != nil {
			return err
		}
		err = repo.UpdateTarea(tar)
		if err != nil {
			return err
		}
	}
	return nil
}
*/
