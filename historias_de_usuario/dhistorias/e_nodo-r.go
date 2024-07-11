package dhistorias

import (
	"fmt"
	"monorepo/historias_de_usuario/ust"

	"github.com/pargomx/gecko"
)

func GetNodo(id int, repo Repo) (*Nodo, error) {
	op := gecko.NewOp("GetNodo").Ctx("id", id)
	nod, err := repo.GetNodo(id)
	if err != nil {
		return nil, op.Err(err)
	}
	nodo := Nodo{
		Nodo: *nod,
	}
	switch {
	case nod.EsPersona():
		nodo.Persona, err = repo.GetPersona(nod.NodoID)
	case nod.EsHistoria():
		nodo.Historia, err = repo.GetNodoHistoria(nod.NodoID)
	case nod.EsTarea():
		nodo.Tarea, err = repo.GetTarea(nod.NodoID)
	default:
		return nil, op.Msgf("tipo de nodo '%v' inválido", nod.NodoTbl)
	}
	if err != nil {
		return nil, op.Err(err)
	}
	return &nodo, nil
}

func GetHijosDeNodo(id int, repo Repo) ([]Nodo, error) {
	op := gecko.NewOp("GetHijosDeNodo").Ctx("id", id)
	nod, err := repo.GetNodo(id)
	if err != nil {
		return nil, op.Err(err)
	}
	hijos := []Nodo{}
	nodos, err := repo.ListNodosByPadreID(nod.NodoID)
	if err != nil {
		return nil, op.Err(err)
	}
	for _, n := range nodos {
		nodo := Nodo{
			Nodo: n,
		}
		switch {
		case n.EsPersona():
			nodo.Persona, err = repo.GetPersona(n.NodoID)
		case n.EsHistoria():
			nodo.Historia, err = repo.GetNodoHistoria(n.NodoID)
		case n.EsTarea():
			nodo.Tarea, err = repo.GetTarea(n.NodoID)
		default:
			return nil, op.Msgf("tipo de nodo '%v' inválido", n.NodoTbl).Ctx("nodo_id", n.NodoID)
		}
		if err != nil {
			return nil, op.Err(err)
		}
		hijos = append(hijos, nodo)
	}
	return hijos, nil
}

func GetNodoConHijos(id int, repo Repo) (*NodoConHijos, error) {
	op := gecko.NewOp("GetNodoConHijos").Ctx("id", id)
	nod, err := GetNodo(id, repo)
	if err != nil {
		return nil, op.Err(err)
	}
	hijos, err := GetHijosDeNodo(nod.Nodo.NodoID, repo)
	if err != nil {
		return nil, op.Err(err)
	}
	padre, _ := GetNodo(nod.Nodo.PadreID, repo)
	// Puede que no tenga padre y se deja nil
	nodo := NodoConHijos{
		Nodo:     nod.Nodo,
		Persona:  nod.Persona,
		Historia: nod.Historia,
		Tarea:    nod.Tarea,
		Padre:    padre,
		Hijos:    hijos,
	}
	return &nodo, nil
}

func GetArbolCompleto(repo Repo) ([]Arbol, error) {
	op := gecko.NewOp("GetArbolCompleto")
	personas, err := repo.ListNodosPersonas()
	if err != nil {
		return nil, op.Err(err)
	}
	arboles := []Arbol{}
	for _, p := range personas {
		arbol := Arbol{
			Persona: p,
		}
		historias, err := repo.ListNodoHistoriasByPadreID(p.PersonaID)
		if err != nil {
			return nil, op.Err(err)
		}
		for _, h := range historias {
			arbol.Historias = append(arbol.Historias, getHistoriaRecursiva(h, repo))
		}
		arboles = append(arboles, arbol)
	}
	return arboles, nil
}

func getHistoriaRecursiva(his ust.NodoHistoria, repo Repo) HistoriaRecursiva {
	historia := HistoriaRecursiva{
		Padre:  his,
		Hijos:  nil,
		Tareas: nil,
	}
	hijos, err := repo.ListNodoHistoriasByPadreID(his.HistoriaID)
	if err != nil {
		fmt.Println("getHistoriaConHijos: %w", err)
	}
	historia.Tareas, err = repo.ListTareasByPadreID(his.HistoriaID)
	if err != nil {
		fmt.Println("getHistoriaConHijos: %w", err)
	}
	for _, hijo := range hijos {
		historia.Hijos = append(historia.Hijos, getHistoriaRecursiva(hijo, repo))
	}
	return historia
}
