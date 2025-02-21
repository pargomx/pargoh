package dhistorias

import (
	"monorepo/ust"

	"github.com/pargomx/gecko/gko"
)

type flagGet int

// Flags para popular tramos, reglas y tareas de cada historia recursiva usando OR.
const (
	GetDescendientes flagGet = 1 << iota
	GetTramos
	GetReglas
	GetTareas
	GetRelacionadas
)

func GetHistoria(historiaID int, flags flagGet, repo Repo) (*HistoriaAgregado, error) {
	op := gko.Op("dhistorias.GetHistoria").Ctx("historiaID", historiaID)

	historia, err := repo.GetNodoHistoria(historiaID)
	if err != nil {
		return nil, op.Err(err)
	}
	item := HistoriaAgregado{
		Historia: *historia,
	}

	item.Tareas, err = repo.ListTareasByHistoriaID(historiaID)
	if err != nil {
		return nil, op.Err(err)
	}

	// Agregar tiempo transcurrido hasta ahora para tareas activas.
	for i, tarea := range item.Tareas {
		if tarea.EnCurso() {
			itrvs, err := repo.ListIntervalosByTareaID(tarea.TareaID)
			if err != nil {
				return nil, op.Err(err)
			}
			for _, itrv := range itrvs {
				if itrv.Fin == "" {
					item.Tareas[i].SegundosUtilizado += itrv.Segundos()
				}
			}
		}
	}

	item.Tramos, err = repo.ListTramosByHistoriaID(historiaID)
	if err != nil {
		return nil, op.Err(err)
	}

	item.Reglas, err = repo.ListReglasByHistoriaID(historiaID)
	if err != nil {
		return nil, op.Err(err)
	}

	item.Relacionadas, err = repo.ListNodoHistoriasRelacionadas(historiaID)
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

	// Obtener historias descendientes.
	if flags&GetDescendientes != 0 {
		item.Descendientes, err = GetHistoriasDescendientes(historiaID, 0, flags, repo)
		if err != nil {
			return nil, op.Err(err)
		}
	}
	return &item, nil
}

// Obtener árbol de historias descendientes de forma recursiva.
// Si el nivel es 1, solo se obtienen las historias inmediatas.
// Si el nivel es 2, se obtienen las historias inmediatas y sus historias inmediatas.
// Y así sucesivamente.
// Si el nivel es 0 o negativo no se limita la recursión y se traen todos los descendientes.
func GetHistoriasDescendientes(padreID int, niveles int, flags flagGet, repo Repo) ([]HistoriaRecursiva, error) {
	historias, err := repo.ListNodoHistoriasByPadreID(padreID)
	if err != nil {
		return nil, gko.Err(err).Strf("padreID:%v niveles:%v", padreID, niveles)
	}
	res := make([]HistoriaRecursiva, len(historias))
	for i, his := range historias {
		res[i].NodoHistoria = his

		if flags&GetTramos != 0 {
			res[i].Tramos, err = repo.ListTramosByHistoriaID(his.HistoriaID)
			if err != nil {
				return nil, err
			}
		}
		if flags&GetReglas != 0 {
			res[i].Reglas, err = repo.ListReglasByHistoriaID(his.HistoriaID)
			if err != nil {
				return nil, err
			}
		}
		if flags&GetTareas != 0 {
			res[i].Tareas, err = repo.ListTareasByHistoriaID(his.HistoriaID)
			if err != nil {
				return nil, err
			}
		}
		if flags&GetRelacionadas != 0 {
			res[i].Relacionadas, err = repo.ListNodoHistoriasRelacionadas(his.HistoriaID)
			if err != nil {
				return nil, err
			}
		}
		if niveles == 1 {
			continue // limitar la recursión cuando se da un nivel positivo.
		}
		res[i].Descendientes, err = GetHistoriasDescendientes(his.HistoriaID, niveles-1, flags, repo)
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}
