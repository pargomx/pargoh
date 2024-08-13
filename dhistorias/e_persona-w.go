package dhistorias

import (
	"monorepo/ust"

	"github.com/pargomx/gecko/gko"
)

func InsertarPersona(per ust.Persona, repo Repo) error {
	op := gko.Op("InsertarPersona")
	if per.Nombre == "" {
		return op.Msg("Persona sin nombre")
	}
	err := repo.ExisteProyecto(per.ProyectoID)
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

func ActualizarPersona(per ust.Persona, repo Repo) error {
	op := gko.Op("ActualizarPersona")
	if per.Nombre == "" {
		return op.Msg("Persona sin nombre")
	}
	err := repo.ExisteProyecto(per.ProyectoID)
	if err != nil {
		return op.Err(err)
	}
	err = repo.UpdatePersona(per)
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
		return op.Msg("Para eliminar una persona, primero elimine sus historias y tareas")
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
