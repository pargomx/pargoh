package sqliteust

import (
	"database/sql"
	"net/http"

	"monorepo/historias_de_usuario/ust"

	"github.com/pargomx/gecko"
)

//  ================================================================  //
//  ========== MYSQL/CONSTANTES ====================================  //

// Lista de columnas separadas por coma para usar en consulta SELECT
// en conjunto con scanRow o scanRows, ya que las columnas coinciden
// con los campos escaneados.
const columnasIntervaloReciente string = "coalesce(tar.historia_id, 0), interv.tarea_id, interv.inicio, interv.fin, coalesce(tar.tipo, ''), coalesce(tar.descripcion, ''), coalesce(tar.impedimentos, ''), coalesce(tar.tiempo_estimado, 0), coalesce(tar.tiempo_real, 0), coalesce(tar.estatus, 0), coalesce(his.titulo, ''), coalesce(his.objetivo, ''), coalesce(his.completada, 0), coalesce(his.prioridad, 0)"

// Origen de los datos de ust.IntervaloReciente
//
// FROM intervalos interv
// INNER JOIN tareas tar ON tar.tarea_id = interv.tarea_id
// INNER JOIN historias his ON his.historia_id = tar.historia_id
const fromIntervaloReciente string = "FROM intervalos interv INNER JOIN tareas tar ON tar.tarea_id = interv.tarea_id INNER JOIN historias his ON his.historia_id = tar.historia_id "

//  ================================================================  //
//  ========== MYSQL/SCAN-ROWS =====================================  //

// scanRowsIntervaloReciente escanea cada row en la struct IntervaloReciente
// y devuelve un slice con todos los items.
// Siempre se encarga de llamar rows.Close()
func (s *Repositorio) scanRowsIntervaloReciente(rows *sql.Rows, op string) ([]ust.IntervaloReciente, error) {
	defer rows.Close()
	items := []ust.IntervaloReciente{}
	for rows.Next() {
		itvr := ust.IntervaloReciente{}
		var tipo string
		err := rows.Scan(
			&itvr.HistoriaID, &itvr.TareaID, &itvr.Inicio, &itvr.Fin, &tipo, &itvr.Descripcion, &itvr.Impedimentos, &itvr.TiempoEstimado, &itvr.TiempoReal, &itvr.Estatus, &itvr.Titulo, &itvr.Objetivo, &itvr.Completada, &itvr.Prioridad,
		)
		if err != nil {
			return nil, gecko.NewErr(http.StatusInternalServerError).Err(err).Op(op)
		}
		itvr.Tipo = ust.SetTipoTareaDB(tipo)
		items = append(items, itvr)
	}
	return items, nil
}

//  ================================================================  //
//  ========== MYSQL/LIST ==========================================  //

func (s *Repositorio) ListIntervalosRecientes() ([]ust.IntervaloReciente, error) {
	const op string = "mysqlust.ListIntervalosRecientes"
	rows, err := s.db.Query(
		"SELECT " + columnasIntervaloReciente + " " + fromIntervaloReciente +
			"WHERE interv.fin <> '' ORDER BY interv.inicio DESC LIMIT 20",
	)
	if err != nil {
		return nil, gecko.NewErr(http.StatusInternalServerError).Err(err).Op(op)
	}
	return s.scanRowsIntervaloReciente(rows, op)
}

func (s *Repositorio) ListIntervalosAbiertos() ([]ust.IntervaloReciente, error) {
	const op string = "mysqlust.ListIntervalosAbiertos"
	rows, err := s.db.Query(
		"SELECT " + columnasIntervaloReciente + " " + fromIntervaloReciente +
			"WHERE interv.fin == '' ORDER BY interv.inicio DESC LIMIT 20",
	)
	if err != nil {
		return nil, gecko.NewErr(http.StatusInternalServerError).Err(err).Op(op)
	}
	return s.scanRowsIntervaloReciente(rows, op)
}
