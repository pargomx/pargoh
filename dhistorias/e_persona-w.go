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
