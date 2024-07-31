package sqliteust

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/pargomx/gecko/gko"

	"monorepo/historias_de_usuario/ust"
)

//  ================================================================  //
//  ========== MYSQL/CONSTANTES ====================================  //

// Lista de columnas separadas por coma para usar en consulta SELECT
// en conjunto con scanRow o scanRows, ya que las columnas coinciden
// con los campos escaneados.
const columnasHistoria string = "historia_id, titulo, objetivo, prioridad, completada"

// Origen de los datos de ust.Historia
//
// FROM historias
const fromHistoria string = "FROM historias "

//  ================================================================  //
//  ========== MYSQL/TBL-INSERT ====================================  //

// InsertHistoria valida el registro y lo inserta en la base de datos.
func (s *Repositorio) InsertHistoria(his ust.Historia) error {
	const op string = "mysqlust.InsertHistoria"
	if his.HistoriaID == 0 {
		return gko.ErrDatoInvalido().Msg("HistoriaID sin especificar").Ctx(op, "pk_indefinida")
	}
	if his.Titulo == "" {
		return gko.ErrDatoInvalido().Msg("Titulo sin especificar").Ctx(op, "required_sin_valor")
	}
	err := his.Validar()
	if err != nil {
		return gko.ErrDatoInvalido().Err(err).Op(op).Msg(err.Error())
	}
	_, err = s.db.Exec("INSERT INTO historias "+
		"(historia_id, titulo, objetivo, prioridad, completada) "+
		"VALUES (?, ?, ?, ?, ?) ",
		his.HistoriaID, his.Titulo, his.Objetivo, his.Prioridad, his.Completada,
	)
	if err != nil {
		if strings.HasPrefix(err.Error(), "Error 1062 (23000)") {
			return gko.ErrYaExiste().Err(err).Op(op)
		} else if strings.HasPrefix(err.Error(), "Error 1452 (23000)") {
			return gko.ErrDatoInvalido().Err(err).Op(op).Msg("No se puede insertar la información porque el registro asociado no existe")
		} else {
			return gko.ErrInesperado().Err(err).Op(op)
		}
	}
	return nil
}

//  ================================================================  //
//  ========== MYSQL/TBL-UPDATE ====================================  //

// UpdateHistoria valida y sobreescribe el registro en la base de datos.
func (s *Repositorio) UpdateHistoria(his ust.Historia) error {
	const op string = "mysqlust.UpdateHistoria"
	if his.HistoriaID == 0 {
		return gko.ErrDatoInvalido().Msg("HistoriaID sin especificar").Ctx(op, "pk_indefinida")
	}
	if his.Titulo == "" {
		return gko.ErrDatoInvalido().Msg("Titulo sin especificar").Ctx(op, "required_sin_valor")
	}
	err := his.Validar()
	if err != nil {
		return gko.ErrDatoInvalido().Err(err).Op(op).Msg(err.Error())
	}
	_, err = s.db.Exec(
		"UPDATE historias SET "+
			"historia_id=?, titulo=?, objetivo=?, prioridad=?, completada=? "+
			"WHERE historia_id = ?",
		his.HistoriaID, his.Titulo, his.Objetivo, his.Prioridad, his.Completada,
		his.HistoriaID,
	)
	if err != nil {
		return gko.ErrInesperado().Err(err).Op(op)
	}
	return nil
}

//  ================================================================  //
//  ========== MYSQL/TBL-DELETE ====================================  //

func (s *Repositorio) DeleteHistoria(HistoriaID int) error {
	const op string = "mysqlust.DeleteHistoria"
	if HistoriaID == 0 {
		return gko.ErrDatoInvalido().Msg("HistoriaID sin especificar").Ctx(op, "pk_indefinida")
	}
	// Verificar que solo se borre un registro.
	var num int
	err := s.db.QueryRow("SELECT COUNT(historia_id) FROM historias WHERE historia_id = ?",
		HistoriaID,
	).Scan(&num)
	if err != nil {
		if err == sql.ErrNoRows {
			return gko.ErrNoEncontrado().Err(ust.ErrHistoriaNotFound).Op(op)
		}
		return gko.ErrInesperado().Err(err).Op(op)
	}
	if num > 1 {
		return gko.ErrInesperado().Err(nil).Op(op).Msgf("abortado porque serían borrados %v registros", num)
	} else if num == 0 {
		return gko.ErrNoEncontrado().Err(ust.ErrHistoriaNotFound).Op(op).Msg("cero resultados")
	}
	// Eliminar registro
	_, err = s.db.Exec(
		"DELETE FROM historias WHERE historia_id = ?",
		HistoriaID,
	)
	if err != nil {
		if strings.HasPrefix(err.Error(), "Error 1451 (23000)") {
			return gko.ErrYaExiste().Err(err).Op(op).Msg("Este registro es referenciado por otros y no se puede eliminar")
		} else {
			return gko.ErrInesperado().Err(err).Op(op)
		}
	}
	return nil
}

//  ================================================================  //
//  ========== MYSQL/SCAN-ROW ======================================  //

// Utilizar luego de un sql.QueryRow(). No es necesario hacer row.Close()
func (s *Repositorio) scanRowHistoria(row *sql.Row, his *ust.Historia, op string) error {

	err := row.Scan(
		&his.HistoriaID, &his.Titulo, &his.Objetivo, &his.Prioridad, &his.Completada,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return gko.ErrNoEncontrado().Msg("Historia de usuario no se encuentra").Op(op)
		}
		return gko.ErrInesperado().Err(err).Op(op)
	}

	return nil
}

//  ================================================================  //
//  ========== MYSQL/GET ===========================================  //

// GetHistoria devuelve un Historia de la DB por su clave primaria.
// Error si no encuentra ninguno, o si encuentra más de uno.
func (s *Repositorio) GetHistoria(HistoriaID int) (*ust.Historia, error) {
	const op string = "mysqlust.GetHistoria"
	if HistoriaID == 0 {
		return nil, gko.ErrDatoInvalido().Msg("HistoriaID sin especificar").Ctx(op, "pk_indefinida")
	}
	row := s.db.QueryRow(
		"SELECT "+columnasHistoria+" "+fromHistoria+
			"WHERE historia_id = ?",
		HistoriaID,
	)
	his := &ust.Historia{}
	return his, s.scanRowHistoria(row, his, op)
}

//  ================================================================  //
//  ========== MYSQL/SCAN-ROWS =====================================  //

// scanRowsHistoria escanea cada row en la struct Historia
// y devuelve un slice con todos los items.
// Siempre se encarga de llamar rows.Close()
func (s *Repositorio) scanRowsHistoria(rows *sql.Rows, op string) ([]ust.Historia, error) {
	defer rows.Close()
	items := []ust.Historia{}
	for rows.Next() {
		his := ust.Historia{}

		err := rows.Scan(
			&his.HistoriaID, &his.Titulo, &his.Objetivo, &his.Prioridad, &his.Completada,
		)
		if err != nil {
			return nil, gko.ErrInesperado().Err(err).Op(op)
		}

		items = append(items, his)
	}
	return items, nil
}

//  ================================================================  //
//  ========== MYSQL/LIST ==========================================  //

func (s *Repositorio) ListHistorias() ([]ust.Historia, error) {
	const op string = "mysqlust.ListHistorias"
	rows, err := s.db.Query(
		"SELECT " + columnasHistoria + " " + fromHistoria,
	)
	if err != nil {
		return nil, gko.ErrInesperado().Err(err).Op(op)
	}
	return s.scanRowsHistoria(rows, op)
}
