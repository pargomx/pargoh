package dhistorias

import (
	"monorepo/ust"

	"github.com/pargomx/gecko/gko"
)

func validarNombreDescrDePersona(per *ust.Persona) error {
	per.Nombre = txtQuitarEspaciosYSaltos(per.Nombre)
	per.Descripcion = txtQuitarEspaciosYSaltos(per.Descripcion)
	if per.Nombre == "" {
		return gko.ErrDatoIndef.Msg("Nombre de la persona indefinido")
	}
	if len(per.Nombre) > 40 {
		return gko.ErrDatoInvalido.Msg("El nombre del personaje no debe superar los 40 caracteres").Strf("Nombre muy largo: %d", len(per.Nombre))
	}
	if len(per.Nombre) > 5000 {
		return gko.ErrDatoInvalido.Msg("La descripción del personaje no debe superar los 5000 caracteres").Strf("Descripcion muy largo: %d", len(per.Descripcion))
	}
	return nil
}

func InsertarPersona(per ust.Persona, repo Repo) error {
	op := gko.Op("InsertarPersona")
	if per.Nombre == "" {
		return op.Msg("Persona sin nombre")
	}
	err := repo.ExisteProyecto(per.ProyectoID)
	if err != nil {
		return op.Err(err)
	}
	err = validarNombreDescrDePersona(&per)
	if err != nil {
		return op.Err(err)
	}
	err = repo.InsertPersona(per)
	if err != nil {
		return op.Err(err)
	}
	err = agregarNodo(ust.RootNodoID, per.PersonaID, ust.TipoNodoPersona, repo)
	if err != nil {
		return op.Err(err)
	}
	return nil
}

// Acepta cambiar de proyecto a la persona.
func ActualizarPersona(per ust.Persona, repo Repo) error {
	op := gko.Op("ActualizarPersona")
	if per.Nombre == "" {
		return op.Msg("Persona sin nombre")
	}
	oldPer, err := repo.GetPersona(per.PersonaID)
	if err != nil {
		return op.Err(err)
	}
	if oldPer.ProyectoID != per.ProyectoID {
		err = repo.ExisteProyecto(per.ProyectoID)
		if err != nil {
			return op.Err(err)
		}
		err := repo.CambiarProyectoDeHistoriasByPersonaID(per.PersonaID, per.ProyectoID)
		if err != nil {
			return op.Err(err)
		}
	}
	err = validarNombreDescrDePersona(&per)
	if err != nil {
		return op.Err(err)
	}
	err = repo.UpdatePersona(per)
	if err != nil {
		return op.Err(err)
	}
	return nil
}

func ParcharPersona(personaID int, param string, newVal string, repo Repo) error {
	op := gko.Op("ParcharPersona").Ctx("personaID", personaID)
	if personaID == 0 {
		return op.Msg("personaID debe estar definido")
	}
	Persona, err := repo.GetPersona(personaID)
	if err != nil {
		return op.Err(err)
	}
	switch param {
	case "nombre":
		Persona.Nombre = newVal
	case "descripcion":
		Persona.Descripcion = newVal
	default:
		return op.Msgf("Parámetro no soportado: %v", param)
	}
	err = repo.UpdatePersona(*Persona)
	if err != nil {
		return op.Err(err)
	}
	return nil
}

func EliminarPersona(personaID int, repo Repo) error {
	op := gko.Op("EliminarPersona")
	per, err := repo.GetPersona(personaID)
	if err != nil {
		return op.Err(err)
	}
	// Verificar que no tenga hijos
	hijos, err := GetHijosDeNodo(per.PersonaID, repo)
	if err != nil {
		return op.Err(err)
	}
	if len(hijos) > 0 {
		return op.Msg("Para eliminar una persona, primero elimine sus historias")
	}
	err = repo.EliminarNodo(per.PersonaID)
	if err != nil {
		return op.Err(err)
	}
	err = repo.DeletePersona(per.PersonaID)
	if err != nil {
		return op.Err(err)
	}
	return nil
}
