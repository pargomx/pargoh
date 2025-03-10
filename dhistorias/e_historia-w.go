package dhistorias

import (
	"monorepo/ust"

	"github.com/pargomx/gecko/gko"
	"github.com/pargomx/gecko/gkt"
)

const prioridadInvalidaMsg = "La prioridad debe estar entre 0 y 3"

func prioridadValida(prioridad int) bool {
	return prioridad >= 0 && prioridad <= 3
}

// ================================================================ //

func AgregarHistoria(padreID int, his ust.Historia, repo Repo) error {
	op := gko.Op("AgregarHistoria").Ctx("padreID", padreID)

	// Validar historia
	if his.HistoriaID == 0 {
		return op.Msg("el ID de la historia debe estar definido")
	}
	if his.Titulo == "" {
		return op.Msg("el título no puede estar vacío")
	}
	if !prioridadValida(his.Prioridad) {
		return op.Msg(prioridadInvalidaMsg)
	}

	// Validar padre
	padre, err := repo.GetNodo(padreID)
	if err != nil {
		return op.Err(err)
	}
	if padre.EsTarea() {
		return op.Msg("el nodo padre es una tarea y no puede tener historias")
	}

	// Popular indice de proyecto_id y persona_id
	if padre.EsPersona() {
		pers, err := repo.GetPersona(padre.NodoID)
		if err != nil {
			return op.Err(err).Op("set_historia_ancestros_index")
		}
		his.PersonaID = pers.PersonaID
		his.ProyectoID = pers.ProyectoID
	} else {
		hisPadre, err := repo.GetHistoria(padre.NodoID)
		if err != nil {
			return op.Err(err).Op("set_historia_ancestros_index")
		}
		his.PersonaID = hisPadre.PersonaID
		his.ProyectoID = hisPadre.ProyectoID
	}

	// Insertar en la base de datos
	err = repo.InsertHistoria(his)
	if err != nil {
		return op.Err(err)
	}
	err = agregarNodo(padreID, his.HistoriaID, "his", repo)
	if err != nil {
		return op.Err(err)
	}

	return nil
}

func ActualizarHistoria(historiaID int, his ust.Historia, repo Repo) error {
	op := gko.Op("ActualizarHistoria").Ctx("historiaID", historiaID)

	if his.HistoriaID == 0 {
		return op.Msg("el ID de la historia debe estar definido")
	}
	if his.Titulo == "" {
		return op.Msg("el título no puede estar vacío")
	}
	if !prioridadValida(his.Prioridad) {
		return op.Msg(prioridadInvalidaMsg)
	}

	oldHis, err := repo.GetHistoria(historiaID)
	if err != nil {
		return op.Err(err)
	}

	if oldHis.HistoriaID != his.HistoriaID {
		return op.Msg("no se puede cambiar el ID de la historia aún")
	}
	oldHis.Titulo = his.Titulo
	oldHis.Objetivo = his.Objetivo
	oldHis.Prioridad = his.Prioridad
	oldHis.Completada = his.Completada

	err = repo.UpdateHistoria(*oldHis)
	if err != nil {
		return op.Err(err)
	}
	return nil
}

func ParcharHistoria(historiaID int, param string, newVal string, repo Repo) error {
	op := gko.Op("ParcharHistoria").Ctx("historiaID", historiaID)
	if historiaID == 0 {
		return op.Msg("el ID de la historia debe estar definido")
	}
	Hist, err := repo.GetHistoria(historiaID)
	if err != nil {
		return op.Err(err)
	}
	switch param {
	case "titulo":
		Hist.Titulo = gkt.SinEspaciosExtra(newVal)

	case "objetivo":
		Hist.Objetivo = gkt.SinEspaciosExtraConSaltos(newVal)

	case "descripcion":
		Hist.Descripcion = gkt.SinEspaciosExtraConSaltos(newVal)

	case "notas":
		Hist.Notas = gkt.SinEspaciosExtraConSaltos(newVal)

	case "prioridad":
		Hist.Prioridad, _ = gkt.ToInt(newVal)
		if !prioridadValida(Hist.Prioridad) {
			return op.Msg(prioridadInvalidaMsg)
		}

	case "completada":
		num, _ := gkt.ToInt(newVal)
		Hist.Completada = num > 0

	case "presupuesto":
		if newVal == "" {
			Hist.SegundosPresupuesto = 0
		} else {
			horas, err := gkt.ToInt(newVal)
			if err != nil {
				return op.Err(err)
			}
			if horas < 0 {
				return op.ErrDatoInvalido().Msg("El presupuesto debe ser positivo")
			}
			if horas > 30 {
				return op.ErrDatoInvalido().Msgf("Establezca un presupuesto menor, %v son demasiadas horas para una sola historia.", horas)
			}
			Hist.SegundosPresupuesto = horas * 3600
		}

	default:
		gko.LogWarnf("Nada cambió para historia %v", Hist.HistoriaID)
		return nil
	}
	err = repo.UpdateHistoria(*Hist)
	if err != nil {
		return op.Err(err)
	}
	return nil
}

func PriorizarHistoria(historiaID int, prioridad int, repo Repo) error {
	op := gko.Op("PriorizarHistoria").Ctx("historiaID", historiaID)

	if !prioridadValida(prioridad) {
		return op.Msg(prioridadInvalidaMsg)
	}

	his, err := repo.GetHistoria(historiaID)
	if err != nil {
		return op.Err(err)
	}
	if his.Prioridad == prioridad {
		return nil
	}
	his.Prioridad = prioridad

	err = repo.UpdateHistoria(*his)
	if err != nil {
		return op.Err(err)
	}
	return nil
}

func MarcarHistoria(historiaID int, completada bool, repo Repo) error {
	op := gko.Op("MarcarHistoria").Ctx("historiaID", historiaID)
	his, err := repo.GetHistoria(historiaID)
	if err != nil {
		return op.Err(err)
	}
	if his.Completada == completada {
		return nil
	}
	his.Completada = completada
	err = repo.UpdateHistoria(*his)
	if err != nil {
		return op.Err(err)
	}
	return nil
}

func EliminarHistoria(historiaID int, repo Repo) (padreID int, err error) {
	op := gko.Op("EliminarHistoria").Ctx("historiaID", historiaID)
	his, err := repo.GetNodoHistoria(historiaID)
	if err != nil {
		return 0, op.Err(err)
	}
	hijos, err := repo.ListNodosByPadreID(his.HistoriaID)
	if err != nil {
		return 0, op.Err(err)
	}
	if len(hijos) > 0 {
		return 0, op.Msg("la historia tiene descendientes y no puede ser eliminada")
	}
	tramos, err := repo.ListTramosByHistoriaID(his.HistoriaID)
	if err != nil {
		return 0, op.Err(err)
	}
	if len(tramos) > 0 {
		return 0, op.Msg("la historia tiene tramos y no puede ser eliminada")
	}
	tareas, err := repo.ListTareasByHistoriaID(his.HistoriaID)
	if err != nil {
		return 0, op.Err(err)
	}
	if len(tareas) > 0 {
		return 0, op.Msg("la historia tiene tareas y no puede ser eliminada")
	}
	reglas, err := repo.ListReglasByHistoriaID(his.HistoriaID)
	if err != nil {
		return 0, op.Err(err)
	}
	if len(reglas) > 0 {
		return 0, op.Msg("la historia tiene reglas y no puede ser eliminada")
	}
	err = repo.EliminarNodo(his.HistoriaID)
	if err != nil {
		return 0, op.Err(err)
	}
	err = repo.DeleteHistoria(his.HistoriaID)
	if err != nil {
		return 0, op.Err(err)
	}
	return his.PadreID, nil
}

func MoverHistoria(historiaID int, nuevoPadreID int, repo Repo) error {
	op := gko.Op("MoverHistoria")
	if historiaID == 0 {
		return op.Msg("No se especificó qué historia mover")
	}
	if nuevoPadreID == 0 {
		return op.Msg("No se especificó a qué padre mover la historia")
	}
	if historiaID == nuevoPadreID {
		return op.Msg("No se puede mover la historia hacia sí misma")
	}

	nodoHistoria, err := repo.GetNodoHistoria(historiaID)
	if err != nil {
		return op.Err(err)
	}
	if nodoHistoria.PadreID == nuevoPadreID {
		return op.Msg("No se moverá porque sigue siendo el mismo padre")
	}

	nuevoPadre, err := repo.GetNodo(nuevoPadreID)
	if err != nil {
		return op.Err(err)
	}
	if !(nuevoPadre.EsPersona() || nuevoPadre.EsHistoria()) {
		return op.Msgf("El nuevo padre debe ser historia o persona pero es %v", nuevoPadre.NodoTbl)
	}
	if nuevoPadre.EsHistoria() {
		nuevaHistoriaPadre, err := GetHistoria(nuevoPadreID, 0, repo)
		if err != nil {
			return op.Err(err)
		}
		for _, ancestro := range nuevaHistoriaPadre.Ancestros {
			if ancestro.HistoriaID == historiaID {
				return op.Msg("La historia no puede ser hija de su propio descendiente")
			}
		}
	}
	err = repo.MoverNodo(historiaID, nuevoPadreID)
	if err != nil {
		return op.Err(err)
	}

	// Actualizar índices proyecto_id y persona_id.
	hist, err := GetHistoria(historiaID, 0, repo)
	if err != nil {
		return op.Err(err).Ctx("historiaID", historiaID)
	}
	his, err := repo.GetHistoria(historiaID)
	if err != nil {
		return op.Err(err).Ctx("historiaID", historiaID)
	}
	his.PersonaID = hist.Persona.PersonaID
	his.ProyectoID = hist.Proyecto.ProyectoID
	err = repo.UpdateHistoria(*his)
	if err != nil {
		return op.Err(err).Ctx("historiaID", historiaID)
	}

	return nil
}

// Actualiza los campos PersonaID y ProyectoID de todas las historias.
func MaterializarAncestrosDeHistorias(repo Repo) error {
	op := gko.Op("MaterializarAncestrosDeHistorias")
	historias, err := repo.ListHistorias()
	if err != nil {
		return op.Err(err)
	}
	for _, his := range historias {
		hist, err := GetHistoria(his.HistoriaID, 0, repo)
		if err != nil {
			return op.Err(err).Ctx("historiaID", his.HistoriaID)
		}
		his.PersonaID = hist.Persona.PersonaID
		his.ProyectoID = hist.Proyecto.ProyectoID
		err = repo.UpdateHistoria(his)
		if err != nil {
			return op.Err(err).Ctx("historiaID", his.HistoriaID)
		}
	}
	return nil
}
