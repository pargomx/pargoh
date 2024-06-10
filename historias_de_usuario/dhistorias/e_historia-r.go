package dhistorias

import (
	"fmt"
	"monorepo/gecko"
	"monorepo/historias_de_usuario/ust"
)

func GetHistoriasDePadre(padreID int, repo Repo) (*HistoriaConNietos, error) {
	op := gecko.NewOp("GetHistoriaConHistorias").Ctx("padreID", padreID)

	padre, err := repo.GetNodo(padreID)
	if err != nil {
		return nil, op.Err(err)
	}

	tareas, err := repo.ListTareasByPadreID(padre.NodoID)
	if err != nil {
		return nil, op.Err(err)
	}

	item := HistoriaConNietos{
		Tareas: tareas,
	}

	if padre.EsTarea() {
		return nil, op.Msgf("el nodo %v es una tarea y no puede ser padre de historias", padreID)
	}

	if padre.EsPersona() {
		persona, err := repo.GetPersona(padre.NodoID)
		if err != nil {
			return nil, op.Err(err)
		}
		item.Persona = *persona
		// fmt.Println("Padre es persona")
	}

	if padre.EsHistoria() {
		historia, err := repo.GetNodoHistoria(padre.NodoID)
		if err != nil {
			return nil, op.Err(err)
		}
		item.Abuelo = historia
		// fmt.Println("Padre es historia")
	}

	// Obtener todos los ancestros hasta llegar a la persona.
	esteNivel := padre.Nivel
	sigAncestroID := padre.PadreID
	for esteNivel > 1 {
		esteAncestro, err := repo.GetNodo(sigAncestroID)
		if err != nil {
			return nil, op.Err(err)
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
				fmt.Printf("ancestro es historia %v nivel %v actually %v \n", esteAncestro.NodoID, esteNivel, esteAncestro.Nivel)
			}

		case esteAncestro.EsPersona():
			persona, err := repo.GetPersona(esteAncestro.NodoID)
			if err != nil {
				return nil, op.Err(err)
			}
			item.Persona = *persona
			if esteNivel != esteAncestro.Nivel {
				fmt.Printf("ancestro es persona %v nivel %v actually %v \n", esteAncestro.NodoID, esteNivel, esteAncestro.Nivel)
			}

		default:
			return nil, op.Msgf("el nodo %v es un %v y no puede ser ancestro de historias", sigAncestroID, esteAncestro.NodoTbl)
		}
	}

	if item.Persona.PersonaID == 0 {
		return nil, op.Msgf("no se encontró la persona del nodo %v en los ancestros", padreID)
	}

	// Obtener todos los hijos y nietos de la historia.
	historias, err := repo.ListNodoHistoriasByPadreID(padre.NodoID)
	if err != nil {
		return nil, op.Err(err)
	}
	for _, hijo := range historias {
		historiasDeHijo, err := repo.ListNodoHistoriasByPadreID(hijo.HistoriaID)
		if err != nil {
			return nil, op.Err(err)
		}
		tareasDeNieto, err := repo.ListTareasByPadreID(hijo.HistoriaID)
		if err != nil {
			return nil, op.Err(err)
		}
		item.Padres = append(item.Padres, HistoriaConHijos{
			Padre:  hijo,
			Hijos:  historiasDeHijo,
			Tareas: tareasDeNieto,
		})
	}
	return &item, nil
}
