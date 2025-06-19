package dhistorias

/*
// Actualiza los campos PersonaID y ProyectoID de todas las historias.
func MaterializarAncestrosDeHistorias(repo Repo) error {
	op := gko.Op("MaterializarAncestrosDeHistorias")
	historias, err := repo.ListHistorias()
	if err != nil {
		return op.Err(err)
	}
	for _, his := range historias {
		hist, err := GetHistoria(his.HistoriaID, 0, repo)
		if err != nil {
			return op.Err(err).Ctx("historiaID", his.HistoriaID)
		}
		his.PersonaID = hist.Persona.PersonaID
		his.ProyectoID = hist.Proyecto.ProyectoID
		err = repo.UpdateHistoria(his)
		if err != nil {
			return op.Err(err).Ctx("historiaID", his.HistoriaID)
		}
	}
	return nil
}
*/
