package dhistorias

import (
	"monorepo/ust"

	"github.com/pargomx/gecko/gko"
)

func GetHistoria(historiaID int, repo Repo) (*Historia, error) {
	op := gko.Op("dhistorias.GetHistoria").Ctx("historiaID", historiaID)

	historia, err := repo.GetNodoHistoria(historiaID)
	if err != nil {
		return nil, op.Err(err)
	}
	item := Historia{
		Historia: *historia,
	}

	item.Tareas, err = repo.ListTareasByHistoriaID(historiaID)
	if err != nil {
		return nil, op.Err(err)
	}

	item.Tramos, err = repo.ListTramosByHistoriaID(historiaID)
	if err != nil {
		return nil, op.Err(err)
	}

	item.Reglas, err = repo.ListReglasByHistoriaID(historiaID)
	if err != nil {
		return nil, op.Err(err)
	}

	// Obtener historias ascendientes hasta llegar a la persona.
	esteNivel := item.Historia.Nivel
	sigAncestroID := item.Historia.PadreID
	for esteNivel > 1 {
		esteAncestro, err := repo.GetNodo(sigAncestroID)
		if err != nil {
			return nil, op.Err(err).Ctx("id", sigAncestroID)
		}
		esteNivel--
		sigAncestroID = esteAncestro.PadreID
		switch {
		case esteAncestro.EsHistoria():
			historia, err := repo.GetNodoHistoria(esteAncestro.NodoID)
			if err != nil {
				return nil, op.Err(err)
			}
			item.Ancestros = append([]ust.NodoHistoria{*historia}, item.Ancestros...) // prepend
			if esteNivel != esteAncestro.Nivel {
				gko.LogWarnf("ancestro es historia %v nivel %v actually %v \n", esteAncestro.NodoID, esteNivel, esteAncestro.Nivel)
			}
		case esteAncestro.EsPersona():
			persona, err := repo.GetPersona(esteAncestro.NodoID)
			if err != nil {
				return nil, op.Err(err)
			}
			item.Persona = *persona
			if esteNivel != esteAncestro.Nivel {
				gko.LogWarnf("ancestro es persona %v nivel %v actually %v \n", esteAncestro.NodoID, esteNivel, esteAncestro.Nivel)
			}
		default:
			return nil, op.Msgf("el nodo %v es un %v y no puede ser ancestro de historias", sigAncestroID, esteAncestro.NodoTbl)
		}
	}
	if item.Persona.PersonaID == 0 {
		return nil, op.Msgf("no se encontró la persona de la historia %v en sus ancestros", historiaID)
	}

	// Obtener proyecto
	proy, err := repo.GetProyecto(item.Persona.ProyectoID)
	if err != nil {
		return nil, op.Err(err)
	}
	item.Proyecto = *proy

	// Obtener historias descendientes 2 niveles.
	historias1, err := repo.ListNodoHistorias(historiaID)
	if err != nil {
		return nil, op.Err(err)
	}
	item.Descendientes = make([]HistoriaRecursiva, len(historias1))
	for i, his1 := range historias1 { // TODO: usar función recursiva.
		item.Descendientes[i].Historia = his1
		historias2, err := repo.ListNodoHistorias(his1.HistoriaID)
		if err != nil {
			return nil, op.Err(err)
		}
		item.Descendientes[i].Descendientes = make([]HistoriaRecursiva, len(historias2))
		for j, his2 := range historias2 {
			item.Descendientes[i].Descendientes[j].Historia = his2
		}
	}
	return &item, nil
}
