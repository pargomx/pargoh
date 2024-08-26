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
	err := repo.ExisteHistoria(historiaID)
	if err != nil {
		return op.Err(err)
	}
	tramos, err := repo.ListTramosByHistoriaID(historiaID)
	if err != nil {
		return op.Err(err)
	}
	if posicion < 1 || posicion > len(tramos) {
		op.Msg("posición de tramo inválida").Ctx("historia", historiaID).Ctx("pos", posicion).Ctx("hermanos", len(tramos)).Alert() // Solo alertar
	}
	err = repo.DeleteTramo(historiaID, posicion)
	if err != nil {
		return op.Err(err)
	}
	return nil
}
