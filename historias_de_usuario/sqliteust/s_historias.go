package sqliteust

import (
	"monorepo/gecko"
	"monorepo/historias_de_usuario/ust"
	"net/http"
)

// ListHistorias devuelve todos los registros de las historias de usuario
func (s *Repositorio) ListHistoriasByPadreID(nodoID int) ([]ust.Historia, error) {
	const op string = "mysqlust.ListHistoriasByPadreID"
	rows, err := s.db.Query(
		"SELECT "+columnasHistoria+" "+fromHistoria+
			"JOIN nodos ON nodo_id = historia_id WHERE padre_id = ? ORDER BY posicion", nodoID,
	)
	if err != nil {
		return nil, gecko.NewErr(http.StatusInternalServerError).Err(err).Op(op)
	}
	return s.scanRowsHistoria(rows, op)
}

func (s *Repositorio) ListHistoriasPrioritarias() ([]ust.NodoHistoria, error) {
	const op string = "mysqlust.ListHistoriasPrioritarias"
	rows, err := s.db.Query(
		"SELECT " + columnasNodoHistoria + " " + fromNodoHistoria +
			"WHERE his.prioridad > 0 AND completada == 0 " +
			"ORDER BY (his.prioridad * nod.nivel) + 20 - nod.posicion DESC LIMIT 50",
	)
	if err != nil {
		return nil, gecko.NewErr(http.StatusInternalServerError).Err(err).Op(op)
	}
	return s.scanRowsNodoHistoria(rows, op)
}

/*

Update:
- set new siblings order++ where order >= me.order
- set me.parent = new parent
- if new parent != old parent {  }

Get:
- get where ID
- list where parent = ID order by order

Insert:
- set parent
- get last sibling order
- set order = last sibling order + 1

Delete:
- delete by ID
- update siblings set order-- where order > deleted.order


*/
