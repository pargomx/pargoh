package sqliteust

import (
	"github.com/pargomx/gecko/gko"
)

// Actualiza el campo materializado proyecto_id para todas las historias utilizando el otro campo materializado persona_id.
func (s *Repositorio) CambiarProyectoDeHistoriasByPersonaID(personaID int, proyectoID string) error {
	const op string = "CambiarProyectoDeHistoriasByPersonaID"
	if personaID == 0 {
		return gko.ErrDatoIndef.Str("personaID sin especificar").Op(op)
	}
	if proyectoID == "" {
		return gko.ErrDatoIndef.Str("proyectoID sin especificar").Op(op)
	}
	_, err := s.db.Exec(
		"UPDATE historias SET proyecto_id = ? WHERE persona_id = ?",
		proyectoID, personaID,
	)
	if err != nil {
		return gko.ErrInesperado.Err(err).Op(op)
	}
	return nil
}
