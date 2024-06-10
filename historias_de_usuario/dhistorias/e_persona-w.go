package dhistorias

import (
	"monorepo/gecko"
	"monorepo/historias_de_usuario/ust"
)

func InsertarPersona(per ust.Persona, repo Repo) error {
	op := gecko.NewOp("InsertarPersona")
	if per.Nombre == "" {
		return op.Msg("Persona sin nombre")
	}
	err := repo.InsertPersona(per)
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
	op := gecko.NewOp("ActualizarPersona")
	if per.Nombre == "" {
		return op.Msg("Persona sin nombre")
	}
	err := repo.UpdatePersona(per)
	if err != nil {
		return op.Err(err)
	}
	return nil
}

func EliminarPersona(personaID int, repo Repo) error {
	op := gecko.NewOp("EliminarPersona")
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
