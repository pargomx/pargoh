package dhistorias

import (
	"monorepo/ust"

	"github.com/pargomx/gecko/gko"
)

func AgregarReferencia(repo Repo, HistoriaID int, refHistoriaID int) error {
	op := gko.Op("AgregarReferencia")
	if HistoriaID == 0 {
		return op.Msg("falta historiaID")
	}
	if refHistoriaID == 0 {
		return op.Msg("falta refHistoriaID")
	}
	if HistoriaID == refHistoriaID {
		return op.Msg("no se puede referenciar a sí misma")
	}
	err := repo.ExisteHistoria(HistoriaID)
	if err != nil {
		return op.Err(err)
	}
	err = repo.ExisteHistoria(refHistoriaID)
	if err != nil {
		return op.Err(err)
	}
	ref := ust.Referencia{
		HistoriaID:    HistoriaID,
		RefHistoriaID: refHistoriaID,
	}
	err = repo.InsertReferencia(ref)
	if err != nil {
		return op.Err(err)
	}
	return nil
}

func EliminarReferencia(repo Repo, HistoriaID int, refHistoriaID int) error {
	op := gko.Op("EliminarReferencia")
	if HistoriaID == 0 {
		return op.Msg("falta historiaID")
	}
	if refHistoriaID == 0 {
		return op.Msg("falta refHistoriaID")
	}
	err := repo.DeleteReferencia(HistoriaID, refHistoriaID)
	if err != nil {
		return op.Err(err)
	}
	return nil
}
