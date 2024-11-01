package sqliteust

import (
	"database/sql"
	"errors"

	"github.com/pargomx/gecko/gko"

	"monorepo/ust"
)

//  ================================================================  //
//  ========== INSERT ==============================================  //

func (s *Repositorio) InsertPersona(per ust.Persona) error {
	const op string = "InsertPersona"
	if per.PersonaID == 0 {
		return gko.ErrDatoIndef().Op(op).Msg("PersonaID sin especificar").Str("pk_indefinida")
	}
	if per.Nombre == "" {
		return gko.ErrDatoIndef().Op(op).Msg("Nombre sin especificar").Str("required_sin_valor")
	}
	_, err := s.db.Exec("INSERT INTO personas "+
		"(persona_id, proyecto_id, nombre, descripcion, segundos_gestion) "+
		"VALUES (?, ?, ?, ?, ?) ",
		per.PersonaID, per.ProyectoID, per.Nombre, per.Descripcion, per.SegundosGestion,
	)
	if err != nil {
		return gko.ErrAlEscribir().Err(err).Op(op)
	}
	return nil
}

//  ================================================================  //
//  ========== UPDATE ==============================================  //

// UpdatePersona valida y sobreescribe el registro en la base de datos.
func (s *Repositorio) UpdatePersona(per ust.Persona) error {
	const op string = "UpdatePersona"
	if per.PersonaID == 0 {
		return gko.ErrDatoIndef().Op(op).Msg("PersonaID sin especificar").Str("pk_indefinida")
	}
	if per.Nombre == "" {
		return gko.ErrDatoIndef().Op(op).Msg("Nombre sin especificar").Str("required_sin_valor")
	}
	_, err := s.db.Exec(
		"UPDATE personas SET "+
			"persona_id=?, proyecto_id=?, nombre=?, descripcion=?, segundos_gestion=? "+
			"WHERE persona_id = ?",
		per.PersonaID, per.ProyectoID, per.Nombre, per.Descripcion, per.SegundosGestion,
		per.PersonaID,
	)
	if err != nil {
		return gko.ErrInesperado().Err(err).Op(op)
	}
	return nil
}

//  ================================================================  //
//  ========== EXISTE ==============================================  //

// Retorna error nil si existe solo un registro con esta clave primaria.
func (s *Repositorio) ExistePersona(PersonaID int) error {
	const op string = "ExistePersona"
	var num int
	err := s.db.QueryRow("SELECT COUNT(persona_id) FROM personas WHERE persona_id = ?",
		PersonaID,
	).Scan(&num)
	if err != nil {
		if err == sql.ErrNoRows {
			return gko.ErrNoEncontrado().Err(ust.ErrPersonaNotFound).Op(op)
		}
		return gko.ErrInesperado().Err(err).Op(op)
	}
	if num > 1 {
		return gko.ErrInesperado().Err(nil).Op(op).Str("existen más de un registro para la pk").Ctx("registros", num)
	} else if num == 0 {
		return gko.ErrNoEncontrado().Err(ust.ErrPersonaNotFound).Op(op)
	}
	return nil
}

//  ================================================================  //
//  ========== DELETE ==============================================  //

func (s *Repositorio) DeletePersona(PersonaID int) error {
	const op string = "DeletePersona"
	if PersonaID == 0 {
		return gko.ErrDatoIndef().Op(op).Msg("PersonaID sin especificar").Str("pk_indefinida")
	}
	err := s.ExistePersona(PersonaID)
	if err != nil {
		return gko.Err(err).Op(op)
	}
	_, err = s.db.Exec(
		"DELETE FROM latidos WHERE persona_id = ?",
		PersonaID,
	)
	if err != nil {
		return gko.ErrAlEscribir().Err(err).Op(op)
	}
	_, err = s.db.Exec(
		"DELETE FROM personas WHERE persona_id = ?",
		PersonaID,
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
//	persona_id,
//	proyecto_id,
//	nombre,
//	descripcion,
//	segundos_gestion
const columnasPersona string = "persona_id, proyecto_id, nombre, descripcion, segundos_gestion"

// Origen de los datos de ust.Persona
//
//	FROM personas
const fromPersona string = "FROM personas "

//  ================================================================  //
//  ========== SCAN ================================================  //

// Utilizar luego de un sql.QueryRow(). No es necesario hacer row.Close()
func (s *Repositorio) scanRowPersona(row *sql.Row, per *ust.Persona) error {
	err := row.Scan(
		&per.PersonaID, &per.ProyectoID, &per.Nombre, &per.Descripcion, &per.SegundosGestion,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return gko.ErrNoEncontrado().Msg("Persona del dominio no se encuentra")
		}
		return gko.ErrInesperado().Err(err)
	}
	return nil
}

//  ================================================================  //
//  ========== GET =================================================  //

// GetPersona devuelve un Persona de la DB por su clave primaria.
// Error si no encuentra ninguno, o si encuentra más de uno.
func (s *Repositorio) GetPersona(PersonaID int) (*ust.Persona, error) {
	const op string = "GetPersona"
	if PersonaID == 0 {
		return nil, gko.ErrDatoIndef().Op(op).Msg("PersonaID sin especificar").Str("pk_indefinida")
	}
	row := s.db.QueryRow(
		"SELECT "+columnasPersona+" "+fromPersona+
			"WHERE persona_id = ?",
		PersonaID,
	)
	per := &ust.Persona{}
	err := s.scanRowPersona(row, per)
	if err != nil {
		return nil, err
	}
	return per, nil
}
