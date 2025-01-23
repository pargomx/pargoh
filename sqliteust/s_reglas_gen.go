package sqliteust

import (
	"database/sql"
	"errors"

	"github.com/pargomx/gecko/gko"

	"monorepo/ust"
)

//  ================================================================  //
//  ========== INSERT ==============================================  //

func (s *Repositorio) InsertRegla(reg ust.Regla) error {
	const op string = "InsertRegla"
	if reg.HistoriaID == 0 {
		return gko.ErrDatoIndef().Op(op).Msg("HistoriaID sin especificar").Str("pk_indefinida")
	}
	if reg.Posicion == 0 {
		return gko.ErrDatoIndef().Op(op).Msg("Posicion sin especificar").Str("pk_indefinida")
	}
	if reg.Texto == "" {
		return gko.ErrDatoIndef().Op(op).Msg("Texto sin especificar").Str("required_sin_valor")
	}
	_, err := s.db.Exec("INSERT INTO reglas "+
		"(historia_id, posicion, texto, implementada, probada) "+
		"VALUES (?, ?, ?, ?, ?) ",
		reg.HistoriaID, reg.Posicion, reg.Texto, reg.Implementada, reg.Probada,
	)
	if err != nil {
		return gko.ErrAlEscribir().Err(err).Op(op)
	}
	return nil
}

//  ================================================================  //
//  ========== UPDATE ==============================================  //

// UpdateRegla valida y sobreescribe el registro en la base de datos.
func (s *Repositorio) UpdateRegla(reg ust.Regla) error {
	const op string = "UpdateRegla"
	if reg.HistoriaID == 0 {
		return gko.ErrDatoIndef().Op(op).Msg("HistoriaID sin especificar").Str("pk_indefinida")
	}
	if reg.Posicion == 0 {
		return gko.ErrDatoIndef().Op(op).Msg("Posicion sin especificar").Str("pk_indefinida")
	}
	if reg.Texto == "" {
		return gko.ErrDatoIndef().Op(op).Msg("Texto sin especificar").Str("required_sin_valor")
	}
	_, err := s.db.Exec(
		"UPDATE reglas SET "+
			"historia_id=?, posicion=?, texto=?, implementada=?, probada=? "+
			"WHERE historia_id = ? AND posicion = ?",
		reg.HistoriaID, reg.Posicion, reg.Texto, reg.Implementada, reg.Probada,
		reg.HistoriaID, reg.Posicion,
	)
	if err != nil {
		return gko.ErrInesperado().Err(err).Op(op)
	}
	return nil
}

//  ================================================================  //
//  ========== EXISTE ==============================================  //

// Retorna error nil si existe solo un registro con esta clave primaria.
func (s *Repositorio) ExisteRegla(HistoriaID int, Posicion int) error {
	const op string = "ExisteRegla"
	var num int
	err := s.db.QueryRow("SELECT COUNT(historia_id) FROM reglas WHERE historia_id = ? AND posicion = ?",
		HistoriaID, Posicion,
	).Scan(&num)
	if err != nil {
		if err == sql.ErrNoRows {
			return gko.ErrNoEncontrado().Err(ust.ErrReglaNotFound).Op(op)
		}
		return gko.ErrInesperado().Err(err).Op(op)
	}
	if num > 1 {
		return gko.ErrInesperado().Err(nil).Op(op).Str("existen más de un registro para la pk").Ctx("registros", num)
	} else if num == 0 {
		return gko.ErrNoEncontrado().Err(ust.ErrReglaNotFound).Op(op)
	}
	return nil
}

//  ================================================================  //
//  ========== DELETE ==============================================  //

func (s *Repositorio) DeleteRegla(HistoriaID int, Posicion int) error {
	const op string = "DeleteRegla"
	if HistoriaID == 0 {
		return gko.ErrDatoIndef().Op(op).Msg("HistoriaID sin especificar").Str("pk_indefinida")
	}
	if Posicion == 0 {
		return gko.ErrDatoIndef().Op(op).Msg("Posicion sin especificar").Str("pk_indefinida")
	}
	err := s.ExisteRegla(HistoriaID, Posicion)
	if err != nil {
		return gko.Err(err).Op(op)
	}
	_, err = s.db.Exec(
		"DELETE FROM reglas WHERE historia_id = ? AND posicion = ?",
		HistoriaID, Posicion,
	)
	if err != nil {
		return gko.ErrAlEscribir().Err(err).Op(op)
	}
	// porque sqlite no mueve los IDs de manera que no haya conflictos.
	_, err = s.db.Exec(
		"UPDATE reglas SET posicion = -(posicion - 1) WHERE historia_id = ? AND posicion > ?",
		HistoriaID, Posicion,
	)
	if err != nil {
		return gko.ErrAlEscribir().Err(err).Op(op)
	}
	_, err = s.db.Exec(
		"UPDATE reglas SET posicion = -posicion WHERE historia_id = ? AND posicion < 0",
		HistoriaID,
	)
	if err != nil {
		return gko.ErrAlEscribir().Err(err).Op(op)
	}
	return nil
}

//  ================================================================  //
//  ========== CONSTANTES ==========================================  //

// Lista de columnas separadas por coma para usar en consulta SELECT
// en conjunto con scanRow o scanRows, ya que las columnas coinciden
// con los campos escaneados.
//
//	historia_id,
//	posicion,
//	texto,
//	implementada,
//	probada
const columnasRegla string = "historia_id, posicion, texto, implementada, probada"

// Origen de los datos de ust.Regla
//
//	FROM reglas
const fromRegla string = "FROM reglas "

//  ================================================================  //
//  ========== SCAN ================================================  //

// Utilizar luego de un sql.QueryRow(). No es necesario hacer row.Close()
func (s *Repositorio) scanRowRegla(row *sql.Row, reg *ust.Regla) error {
	err := row.Scan(
		&reg.HistoriaID, &reg.Posicion, &reg.Texto, &reg.Implementada, &reg.Probada,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return gko.ErrNoEncontrado().Msg("Regla no se encuentra")
		}
		return gko.ErrInesperado().Err(err)
	}
	return nil
}

//  ================================================================  //
//  ========== GET =================================================  //

// GetRegla devuelve un Regla de la DB por su clave primaria.
// Error si no encuentra ninguno, o si encuentra más de uno.
func (s *Repositorio) GetRegla(HistoriaID int, Posicion int) (*ust.Regla, error) {
	const op string = "GetRegla"
	if HistoriaID == 0 {
		return nil, gko.ErrDatoIndef().Op(op).Msg("HistoriaID sin especificar").Str("pk_indefinida")
	}
	if Posicion == 0 {
		return nil, gko.ErrDatoIndef().Op(op).Msg("Posicion sin especificar").Str("pk_indefinida")
	}
	row := s.db.QueryRow(
		"SELECT "+columnasRegla+" "+fromRegla+
			"WHERE historia_id = ? AND posicion = ?",
		HistoriaID, Posicion,
	)
	reg := &ust.Regla{}
	err := s.scanRowRegla(row, reg)
	if err != nil {
		return nil, err
	}
	return reg, nil
}

//  ================================================================  //
//  ========== SCAN ================================================  //

// scanRowsRegla escanea cada row en la struct Regla
// y devuelve un slice con todos los items.
// Siempre se encarga de llamar rows.Close()
func (s *Repositorio) scanRowsRegla(rows *sql.Rows, op string) ([]ust.Regla, error) {
	defer rows.Close()
	items := []ust.Regla{}
	for rows.Next() {
		reg := ust.Regla{}
		err := rows.Scan(
			&reg.HistoriaID, &reg.Posicion, &reg.Texto, &reg.Implementada, &reg.Probada,
		)
		if err != nil {
			return nil, gko.ErrInesperado().Err(err).Op(op)
		}
		items = append(items, reg)
	}
	return items, nil
}

//  ================================================================  //
//  ========== LIST_BY =============================================  //

func (s *Repositorio) ListReglasByHistoriaID(HistoriaID int) ([]ust.Regla, error) {
	const op string = "ListReglasByHistoriaID"
	if HistoriaID == 0 {
		return nil, gko.ErrDatoIndef().Op(op).Msg("HistoriaID sin especificar").Str("param_indefinido")
	}
	rows, err := s.db.Query(
		"SELECT "+columnasRegla+" "+fromRegla+
			"WHERE historia_id = ?",
		HistoriaID,
	)
	if err != nil {
		return nil, gko.ErrInesperado().Err(err).Op(op)
	}
	return s.scanRowsRegla(rows, op)
}
