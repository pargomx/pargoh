package dhistorias

import (
	"monorepo/ust"

	"github.com/pargomx/gecko/gko"
)

func AgregarTramoDeViaje(repo Repo, historiaID int, texto string) error {
	op := gko.Op("NuevoTramoDeViaje")
	if historiaID == 0 {
		return op.Msg("Debe asignarse un ID de historia al tramo de viaje")
	}
	if texto == "" {
		return op.Msg("El texto del tramo de viaje no puede estar vacío")
	}
	err := repo.ExisteHistoria(historiaID)
	if err != nil {
		return op.Err(err)
	}
	tramos, err := repo.ListTramosByHistoriaID(historiaID)
	if err != nil {
		return op.Err(err)
	}
	tramo := ust.Tramo{
		HistoriaID: historiaID,
		Texto:      texto,
		Posicion:   len(tramos) + 1,
	}
	err = repo.InsertTramo(tramo)
	if err != nil {
		return op.Err(err)
	}
	return nil
}

func EliminarTramoDeViaje(repo Repo, historiaID int, posicion int) error {
	op := gko.Op("EliminarTramoDeViaje")
	if historiaID == 0 {
		return op.Msg("falta historiaID")
	}
	if posicion == 0 {
		return op.Msg("falta posición del tramo")
	}
	Tramo, err := repo.GetTramo(historiaID, posicion)
	if err != nil {
		return op.Err(err)
	}
	if Tramo.Imagen != "" {
		return op.Msg("Antes de eliminar el tramo quite la imagen")
	}
	err = repo.DeleteTramo(historiaID, posicion)
	if err != nil {
		return op.Err(err)
	}
	return nil
}

// Si el texto está vacío, elimina el tramo.
func EditarTramoDeViaje(repo Repo, historiaID int, posicion int, texto string) error {
	op := gko.Op("EditarTramoDeViaje")
	tramo, err := repo.GetTramo(historiaID, posicion)
	if err != nil {
		return op.Err(err)
	}
	if texto == "" {
		// return op.Msg("El texto no puede estar vacío")
		return EliminarTramoDeViaje(repo, historiaID, posicion)
	}
	tramo.Texto = texto
	err = repo.UpdateTramo(*tramo)
	if err != nil {
		return op.Err(err)
	}
	return nil
}

func ReordenarTramo(repo Repo, historiaID, oldPos, newPos int) error {
	if historiaID == 0 {
		return gko.Op("ReordenarTramo").Msg("falta historiaID")
	}
	if oldPos == newPos {
		return nil
	}
	err := repo.ReordenarTramo(historiaID, oldPos, newPos)
	if err != nil {
		return gko.Op("ReordenarTramo").Err(err)
	}
	return nil
}

func MoverTramo(historiaID int, posicion int, newHistoriaID int, repo Repo) (int, error) {
	op := gko.Op("MoverTramo")
	err := repo.ExisteTramo(historiaID, posicion)
	if err != nil {
		return 0, op.Err(err)
	}
	err = repo.ExisteHistoria(newHistoriaID)
	if err != nil {
		return 0, op.Err(err).Msg("Debe mover el tramo a una historia, no persona ni proyecto")
	}
	err = repo.MoverTramo(historiaID, posicion, newHistoriaID)
	if err != nil {
		return 0, err
	}
	return newHistoriaID, nil
}
