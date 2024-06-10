package sqliteust

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"

	"monorepo/gecko"
	"monorepo/historias_de_usuario/ust"
)

//  ================================================================  //
//  ========== MYSQL/CONSTANTES ====================================  //

// Lista de columnas separadas por coma para usar en consulta SELECT
// en conjunto con scanRow o scanRows, ya que las columnas coinciden
// con los campos escaneados.
const columnasPersona string = "persona_id, nombre, descripcion"

// Origen de los datos de ust.Persona
//
// FROM personas
const fromPersona string = "FROM personas "

//  ================================================================  //
//  ========== MYSQL/TBL-INSERT ====================================  //

// InsertPersona valida el registro y lo inserta en la base de datos.
func (s *Repositorio) InsertPersona(per ust.Persona) error {
	const op string = "mysqlust.InsertPersona"
	if per.PersonaID == 0 {
		return gecko.NewErr(http.StatusBadRequest).Msg("PersonaID sin especificar").Ctx(op, "pk_indefinida")
	}
	if per.Nombre == "" {
		return gecko.NewErr(http.StatusBadRequest).Msg("Nombre sin especificar").Ctx(op, "required_sin_valor")
	}
	err := per.Validar()
	if err != nil {
		return gecko.NewErr(http.StatusBadRequest).Err(err).Op(op).Msg(err.Error())
	}
	_, err = s.db.Exec("INSERT INTO personas "+
		"(persona_id, nombre, descripcion) "+
		"VALUES (?, ?, ?) ",
		per.PersonaID, per.Nombre, per.Descripcion,
	)
	if err != nil {
		if strings.HasPrefix(err.Error(), "Error 1062 (23000)") {
			return gecko.NewErr(http.StatusConflict).Err(err).Op(op)
		} else if strings.HasPrefix(err.Error(), "Error 1452 (23000)") {
			return gecko.NewErr(http.StatusBadRequest).Err(err).Op(op).Msg("No se puede insertar la información porque el registro asociado no existe")
		} else {
			return gecko.NewErr(http.StatusInternalServerError).Err(err).Op(op)
		}
	}
	return nil
}

//  ================================================================  //
//  ========== MYSQL/TBL-UPDATE ====================================  //

// UpdatePersona valida y sobreescribe el registro en la base de datos.
func (s *Repositorio) UpdatePersona(per ust.Persona) error {
	const op string = "mysqlust.UpdatePersona"
	if per.PersonaID == 0 {
		return gecko.NewErr(http.StatusBadRequest).Msg("PersonaID sin especificar").Ctx(op, "pk_indefinida")
	}
	if per.Nombre == "" {
		return gecko.NewErr(http.StatusBadRequest).Msg("Nombre sin especificar").Ctx(op, "required_sin_valor")
	}
	err := per.Validar()
	if err != nil {
		return gecko.NewErr(http.StatusBadRequest).Err(err).Op(op).Msg(err.Error())
	}
	_, err = s.db.Exec(
		"UPDATE personas SET "+
			"persona_id=?, nombre=?, descripcion=? "+
			"WHERE persona_id = ?",
		per.PersonaID, per.Nombre, per.Descripcion,
		per.PersonaID,
	)
	if err != nil {
		return gecko.NewErr(http.StatusInternalServerError).Err(err).Op(op)
	}
	return nil
}

//  ================================================================  //
//  ========== MYSQL/TBL-DELETE ====================================  //

// DeletePersona elimina permanentemente un registro de la persona del dominio de la base de datos.
// Error si el registro no existe o si no se da la clave primaria.
func (s *Repositorio) DeletePersona(PersonaID int) error {
	const op string = "mysqlust.DeletePersona"
	if PersonaID == 0 {
		return gecko.NewErr(http.StatusBadRequest).Msg("PersonaID sin especificar").Ctx(op, "pk_indefinida")
	}
	// Verificar que solo se borre un registro.
	var num int
	err := s.db.QueryRow("SELECT COUNT(persona_id) FROM personas WHERE persona_id = ?",
		PersonaID,
	).Scan(&num)
	if err != nil {
		if err == sql.ErrNoRows {
			return gecko.NewErr(http.StatusNotFound).Err(ust.ErrPersonaNotFound).Op(op)
		}
		return gecko.NewErr(http.StatusInternalServerError).Err(err).Op(op)
	}
	if num > 1 {
		return gecko.NewErr(http.StatusInternalServerError).Err(nil).Op(op).Msgf("abortado porque serían borrados %v registros", num)
	} else if num == 0 {
		return gecko.NewErr(http.StatusNotFound).Err(ust.ErrPersonaNotFound).Op(op).Msg("cero resultados")
	}
	// Eliminar registro
	_, err = s.db.Exec(
		"DELETE FROM personas WHERE persona_id = ?",
		PersonaID,
	)
	if err != nil {
		if strings.HasPrefix(err.Error(), "Error 1451 (23000)") {
			return gecko.NewErr(http.StatusConflict).Err(err).Op(op).Msg("Este registro es referenciado por otros y no se puede eliminar")
		} else {
			return gecko.NewErr(http.StatusInternalServerError).Err(err).Op(op)
		}
	}
	return nil
}

//  ================================================================  //
//  ========== MYSQL/SCAN-ROW ======================================  //

// Utilizar luego de un sql.QueryRow(). No es necesario hacer row.Close()
func (s *Repositorio) scanRowPersona(row *sql.Row, per *ust.Persona, op string) error {

	err := row.Scan(
		&per.PersonaID, &per.Nombre, &per.Descripcion,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return gecko.NewErr(http.StatusNotFound).Msg("la persona del dominio no se encuentra").Op(op)
		}
		return gecko.NewErr(http.StatusInternalServerError).Err(err).Op(op)
	}

	return nil
}

//  ================================================================  //
//  ========== MYSQL/GET ===========================================  //

// GetPersona devuelve un Persona de la DB por su clave primaria.
// Error si no encuentra ninguno, o si encuentra más de uno.
func (s *Repositorio) GetPersona(PersonaID int) (*ust.Persona, error) {
	const op string = "mysqlust.GetPersona"
	if PersonaID == 0 {
		return nil, gecko.NewErr(http.StatusBadRequest).Msg("PersonaID sin especificar").Ctx(op, "pk_indefinida")
	}
	row := s.db.QueryRow(
		"SELECT "+columnasPersona+" "+fromPersona+
			"WHERE persona_id = ?",
		PersonaID,
	)
	per := &ust.Persona{}
	return per, s.scanRowPersona(row, per, op)
}

//  ================================================================  //
//  ========== MYSQL/SCAN-ROWS =====================================  //

// scanRowsPersona escanea cada row en la struct Persona
// y devuelve un slice con todos los items.
// Siempre se encarga de llamar rows.Close()
func (s *Repositorio) scanRowsPersona(rows *sql.Rows, op string) ([]ust.Persona, error) {
	defer rows.Close()
	items := []ust.Persona{}
	for rows.Next() {
		per := ust.Persona{}

		err := rows.Scan(
			&per.PersonaID, &per.Nombre, &per.Descripcion,
		)
		if err != nil {
			return nil, gecko.NewErr(http.StatusInternalServerError).Err(err).Op(op)
		}

		items = append(items, per)
	}
	return items, nil
}

//  ================================================================  //
//  ========== MYSQL/LIST ==========================================  //

// ListPersonas devuelve todos los registros de las personajes
func (s *Repositorio) ListPersonas() ([]ust.Persona, error) {
	const op string = "mysqlust.ListPersonas"
	rows, err := s.db.Query(
		"SELECT " + columnasPersona + " " + fromPersona,
	)
	if err != nil {
		return nil, gecko.NewErr(http.StatusInternalServerError).Err(err).Op(op)
	}
	return s.scanRowsPersona(rows, op)
}
