package dhistorias

import (
	"monorepo/ust"
	"strings"

	"github.com/pargomx/gecko/gko"
)

func AgregarRegla(repo Repo, historiaID int, texto string) error {
	op := gko.Op("NuevoRegla")
	if historiaID == 0 {
		return op.Msg("Debe asignarse un ID de historia a la regla")
	}
	if texto == "" {
		return op.Msg("El texto del regla de viaje no puede estar vacío")
	}
	err := repo.ExisteHistoria(historiaID)
	if err != nil {
		return op.Err(err)
	}
	reglas, err := repo.ListReglasByHistoriaID(historiaID)
	if err != nil {
		return op.Err(err)
	}
	regla := ust.Regla{
		HistoriaID: historiaID,
		Texto:      strings.TrimSpace(texto),
		Posicion:   len(reglas) + 1,
	}
	err = repo.InsertRegla(regla)
	if err != nil {
		return op.Err(err)
	}
	return nil
}

func EliminarRegla(repo Repo, historiaID int, posicion int) error {
	op := gko.Op("EliminarRegla")
	if historiaID == 0 {
		return op.Msg("falta historiaID")
	}
	if posicion == 0 {
		return op.Msg("falta posición del regla")
	}
	err := repo.ExisteHistoria(historiaID)
	if err != nil {
		return op.Err(err)
	}
	reglas, err := repo.ListReglasByHistoriaID(historiaID)
	if err != nil {
		return op.Err(err)
	}
	if posicion < 1 || posicion > len(reglas) {
		op.Msg("posición de regla inválida").Ctx("historia", historiaID).Ctx("pos", posicion).Ctx("hermanos", len(reglas)).Alert() // Solo alertar
	}
	err = repo.DeleteRegla(historiaID, posicion)
	if err != nil {
		return op.Err(err)
	}
	return nil
}

// Si el texto está vacío, elimina el tramo.
func EditarRegla(repo Repo, historiaID int, posicion int, texto string) error {
	op := gko.Op("EditarRegla")
	regla, err := repo.GetRegla(historiaID, posicion)
	if err != nil {
		return op.Err(err)
	}
	if texto == "" {
		// return op.Msg("El texto no puede estar vacío")
		return EliminarRegla(repo, historiaID, posicion)
	}
	regla.Texto = strings.TrimSpace(texto)
	err = repo.UpdateRegla(*regla)
	if err != nil {
		return op.Err(err)
	}
	return nil
}

func MarcarRegla(repo Repo, historiaID int, posicion int) error {
	op := gko.Op("MarcarRegla")
	regla, err := repo.GetRegla(historiaID, posicion)
	if err != nil {
		return op.Err(err)
	}
	switch {
	case !regla.Implementada && !regla.Probada:
		regla.Implementada = true
	case regla.Implementada && !regla.Probada:
		regla.Probada = true
	default:
		regla.Implementada = false
		regla.Probada = false
	}
	err = repo.UpdateRegla(*regla)
	if err != nil {
		return op.Err(err)
	}
	return nil
}

func ReordenarRegla(repo Repo, historiaID, oldPos, newPos int) error {
	if historiaID == 0 {
		return gko.Op("ReordenarRegla").Msg("falta historiaID")
	}
	if oldPos == newPos {
		return nil
	}
	err := repo.ReordenarRegla(historiaID, oldPos, newPos)
	if err != nil {
		return gko.Op("ReordenarRegla").Err(err)
	}
	return nil
}
