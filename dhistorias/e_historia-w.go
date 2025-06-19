package dhistorias

import (
	"github.com/pargomx/gecko/gko"
)

const prioridadInvalidaMsg = "La prioridad debe estar entre 0 y 3"

func prioridadValida(prioridad int) bool {
	return prioridad >= 0 && prioridad <= 3
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

/*
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
*/
