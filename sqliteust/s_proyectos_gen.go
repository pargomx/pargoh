package sqliteust

import (
	"database/sql"
	"errors"

	"github.com/pargomx/gecko/gko"

	"monorepo/ust"
)

//  ================================================================  //
//  ========== INSERT ==============================================  //

func (s *Repositorio) InsertProyecto(pro ust.Proyecto) error {
	const op string = "InsertProyecto"
	if pro.ProyectoID == "" {
		return gko.ErrDatoIndef().Op(op).Msg("ProyectoID sin especificar").Str("pk_indefinida")
	}
	if pro.Titulo == "" {
		return gko.ErrDatoIndef().Op(op).Msg("Titulo sin especificar").Str("required_sin_valor")
	}
	_, err := s.db.Exec("INSERT INTO proyectos "+
		"(proyecto_id, titulo, imagen, descripcion) "+
		"VALUES (?, ?, ?, ?) ",
		pro.ProyectoID, pro.Titulo, pro.Imagen, pro.Descripcion,
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
//	proyecto_id,
//	titulo,
//	imagen,
//	descripcion
const columnasProyecto string = "proyecto_id, titulo, imagen, descripcion"

// Origen de los datos de ust.Proyecto
//
//	FROM proyectos
const fromProyecto string = "FROM proyectos "

//  ================================================================  //
//  ========== SCAN ================================================  //

// Utilizar luego de un sql.QueryRow(). No es necesario hacer row.Close()
func (s *Repositorio) scanRowProyecto(row *sql.Row, pro *ust.Proyecto) error {
	err := row.Scan(
		&pro.ProyectoID, &pro.Titulo, &pro.Imagen, &pro.Descripcion,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return gko.ErrNoEncontrado().Msg("Proyecto no se encuentra")
		}
		return gko.ErrInesperado().Err(err)
	}
	return nil
}

//  ================================================================  //
//  ========== GET =================================================  //

// GetProyecto devuelve un Proyecto de la DB por su clave primaria.
// Error si no encuentra ninguno, o si encuentra más de uno.
func (s *Repositorio) GetProyecto(ProyectoID string) (*ust.Proyecto, error) {
	const op string = "GetProyecto"
	if ProyectoID == "" {
		return nil, gko.ErrDatoIndef().Op(op).Msg("ProyectoID sin especificar").Str("pk_indefinida")
	}
	row := s.db.QueryRow(
		"SELECT "+columnasProyecto+" "+fromProyecto+
			"WHERE proyecto_id = ?",
		ProyectoID,
	)
	pro := &ust.Proyecto{}
	err := s.scanRowProyecto(row, pro)
	if err != nil {
		return nil, err
	}
	return pro, nil
}

//  ================================================================  //
//  ========== UPDATE ==============================================  //

// UpdateProyecto valida y sobreescribe el registro en la base de datos.
func (s *Repositorio) UpdateProyecto(pro ust.Proyecto) error {
	const op string = "UpdateProyecto"
	if pro.ProyectoID == "" {
		return gko.ErrDatoIndef().Op(op).Msg("ProyectoID sin especificar").Str("pk_indefinida")
	}
	if pro.Titulo == "" {
		return gko.ErrDatoIndef().Op(op).Msg("Titulo sin especificar").Str("required_sin_valor")
	}
	_, err := s.db.Exec(
		"UPDATE proyectos SET "+
			"proyecto_id=?, titulo=?, imagen=?, descripcion=? "+
			"WHERE proyecto_id = ?",
		pro.ProyectoID, pro.Titulo, pro.Imagen, pro.Descripcion,
		pro.ProyectoID,
	)
	if err != nil {
		return gko.ErrInesperado().Err(err).Op(op)
	}
	return nil
}

//  ================================================================  //
//  ========== EXISTE ==============================================  //

// Retorna error nil si existe solo un registro con esta clave primaria.
func (s *Repositorio) ExisteProyecto(ProyectoID string) error {
	const op string = "ExisteProyecto"
	var num int
	err := s.db.QueryRow("SELECT COUNT(proyecto_id) FROM proyectos WHERE proyecto_id = ?",
		ProyectoID,
	).Scan(&num)
	if err != nil {
		if err == sql.ErrNoRows {
			return gko.ErrNoEncontrado().Err(ust.ErrProyectoNotFound).Op(op)
		}
		return gko.ErrInesperado().Err(err).Op(op)
	}
	if num > 1 {
		return gko.ErrInesperado().Err(nil).Op(op).Str("existen más de un registro para la pk").Ctx("registros", num)
	} else if num == 0 {
		return gko.ErrNoEncontrado().Err(ust.ErrProyectoNotFound).Op(op)
	}
	return nil
}

//  ================================================================  //
//  ========== DELETE ==============================================  //

func (s *Repositorio) DeleteProyecto(ProyectoID string) error {
	const op string = "DeleteProyecto"
	if ProyectoID == "" {
		return gko.ErrDatoIndef().Op(op).Msg("ProyectoID sin especificar").Str("pk_indefinida")
	}
	err := s.ExisteProyecto(ProyectoID)
	if err != nil {
		return gko.Err(err).Op(op)
	}
	_, err = s.db.Exec(
		"DELETE FROM proyectos WHERE proyecto_id = ?",
		ProyectoID,
	)
	if err != nil {
		return gko.ErrAlEscribir().Err(err).Op(op)
	}
	return nil
}

//  ================================================================  //
//  ========== SCAN ================================================  //

// scanRowsProyecto escanea cada row en la struct Proyecto
// y devuelve un slice con todos los items.
// Siempre se encarga de llamar rows.Close()
func (s *Repositorio) scanRowsProyecto(rows *sql.Rows, op string) ([]ust.Proyecto, error) {
	defer rows.Close()
	items := []ust.Proyecto{}
	for rows.Next() {
		pro := ust.Proyecto{}
		err := rows.Scan(
			&pro.ProyectoID, &pro.Titulo, &pro.Imagen, &pro.Descripcion,
		)
		if err != nil {
			return nil, gko.ErrInesperado().Err(err).Op(op)
		}
		items = append(items, pro)
	}
	return items, nil
}

//  ================================================================  //
//  ========== LIST ================================================  //

func (s *Repositorio) ListProyectos() ([]ust.Proyecto, error) {
	const op string = "ListProyectos"
	rows, err := s.db.Query(
		"SELECT " + columnasProyecto + " " + fromProyecto,
	)
	if err != nil {
		return nil, gko.ErrInesperado().Err(err).Op(op)
	}
	return s.scanRowsProyecto(rows, op)
}
