package dhistorias

import (
	"monorepo/ust"
	"strings"

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
		Texto:      strings.TrimSpace(texto),
		Posicion:   len(tramos) + 1,
	}
	err = repo.InsertTramo(tramo)
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
	texto = strings.TrimSpace(texto)
	if texto == "" {
		return op.Msg("El texto no puede estar vacío")
		// return EliminarTramoDeViaje(repo, historiaID, posicion)
	}
	tramo.Texto = texto
	err = repo.UpdateTramo(*tramo)
	if err != nil {
		return op.Err(err)
	}
	return nil
}
