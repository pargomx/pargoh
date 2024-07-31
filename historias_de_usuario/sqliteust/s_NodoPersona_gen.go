package sqliteust

import (
	"database/sql"

	"monorepo/historias_de_usuario/ust"

	"github.com/pargomx/gecko/gko"
)

//  ================================================================  //
//  ========== MYSQL/CONSTANTES ====================================  //

// Lista de columnas separadas por coma para usar en consulta SELECT
// en conjunto con scanRow o scanRows, ya que las columnas coinciden
// con los campos escaneados.
const columnasNodoPersona string = "per.persona_id, per.nombre, per.descripcion, coalesce(nod.padre_id, 0), coalesce(nod.padre_tbl, ''), coalesce(nod.nivel, 0), coalesce(nod.posicion, 0)"

// Origen de los datos de ust.NodoPersona
//
// FROM personas per
// INNER JOIN nodos nod ON nodo_id = persona_id
const fromNodoPersona string = "FROM personas per INNER JOIN nodos nod ON nodo_id = persona_id "

//  ================================================================  //
//  ========== MYSQL/SCAN-ROWS =====================================  //

// scanRowsNodoPersona escanea cada row en la struct NodoPersona
// y devuelve un slice con todos los items.
// Siempre se encarga de llamar rows.Close()
func (s *Repositorio) scanRowsNodoPersona(rows *sql.Rows, op string) ([]ust.NodoPersona, error) {
	defer rows.Close()
	items := []ust.NodoPersona{}
	for rows.Next() {
		nper := ust.NodoPersona{}

		err := rows.Scan(
			&nper.PersonaID, &nper.Nombre, &nper.Descripcion, &nper.PadreID, &nper.PadreTbl, &nper.Nivel, &nper.Posicion,
		)
		if err != nil {
			return nil, gko.ErrInesperado().Err(err).Op(op)
		}

		items = append(items, nper)
	}
	return items, nil
}

//  ================================================================  //
//  ========== MYSQL/LIST ==========================================  //

func (s *Repositorio) ListNodosPersonas() ([]ust.NodoPersona, error) {
	const op string = "mysqlust.ListNodosPersonas"
	rows, err := s.db.Query(
		"SELECT " + columnasNodoPersona + " " + fromNodoPersona + " ORDER BY nod.posicion",
	)
	if err != nil {
		return nil, gko.ErrInesperado().Err(err).Op(op)
	}
	return s.scanRowsNodoPersona(rows, op)
}
