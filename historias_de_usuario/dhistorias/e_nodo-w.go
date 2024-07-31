package dhistorias

import (
	"monorepo/historias_de_usuario/ust"

	"github.com/pargomx/gecko/gko"
)

// Inserta el nodo en la última posición dentro de los hijos del nodo padre dado.
func agregarNodo(padreID int, nodoID int, tipo string, repo Repo) error {
	op := gko.Op("dhistorias.agregarNodo").Ctx("padreID", padreID).Ctx("tipo", tipo)
	if padreID == 0 && tipo != ust.TipoNodoPersona {
		return op.Msg("padreID sin especificar para nodo no persona")
	}
	if nodoID == 0 {
		return op.Msg("nodoID sin especificar")
	}
	if tipo == "" {
		return op.Msg("tipo de nodo sin especificar")
	}
	padre, err := repo.GetNodo(padreID)
	if err != nil {
		return op.Op("get_nodo_padre").Err(err)
	}
	nod := ust.Nodo{
		NodoID:   nodoID,
		NodoTbl:  tipo,
		PadreID:  padre.NodoID,
		PadreTbl: padre.NodoTbl,
	}
	if !nod.EsPersona() && !nod.EsHistoria() && !nod.EsTarea() {
		return op.Msg("tipo de nodo inválido")
	}
	err = repo.InsertNodo(nod)
	if err != nil {
		return op.Err(err)
	}
	return nil
}

func ReordenarNodo(nodoID int, newPosicion int, repo Repo) error {
	nodo, err := repo.GetNodo(nodoID)
	if err != nil {
		return err
	}
	err = repo.ReordenarNodo(nodo.NodoID, nodo.Posicion, newPosicion)
	if err != nil {
		return err
	}
	return nil
}
