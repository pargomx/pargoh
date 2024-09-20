package dhistorias

import "github.com/pargomx/gecko/gko"

func (p *ProyectoExport) EliminarPorCompleto(repo Repo) error {
	for _, persona := range p.Personas {
		for _, historia := range persona.Historias {
			err := deleteHistoriaRecursiva(historia, repo)
			if err != nil {
				return err
			}
		}
		err := EliminarPersona(persona.Persona.PersonaID, repo)
		if err != nil {
			return err
		}
	}
	err := EliminarProyecto(p.Proyecto.ProyectoID, repo)
	if err != nil {
		return err
	}
	gko.LogEventof("Proyecto %v eliminado por completo", p.Proyecto.ProyectoID)
	return nil
}

func deleteHistoriaRecursiva(his HistoriaExport, repo Repo) error {
	for _, tramo := range his.Tramos {
		if tramo.Imagen != "" {
			err := EliminarFotoTramo(tramo.HistoriaID, tramo.Posicion, "imagenes", repo) // Todo: parametrizar directorio
			if err != nil {
				return err
			}
		}
	}
	err := repo.DeleteAllTramos(his.Historia.HistoriaID)
	if err != nil {
		return err
	}
	for _, tarea := range his.Tareas {
		for _, intervalo := range tarea.Intervalos {
			err := repo.DeleteIntervalo(intervalo.TareaID, intervalo.Inicio)
			if err != nil {
				return err
			}
		}
	}
	err = repo.DeleteAllTareas(his.Historia.HistoriaID)
	if err != nil {
		return err
	}
	err = repo.DeleteAllReglas(his.Historia.HistoriaID)
	if err != nil {
		return err
	}
	for _, h := range his.Historias {
		err := deleteHistoriaRecursiva(h, repo)
		if err != nil {
			return err
		}
	}
	_, err = EliminarHistoria(his.Historia.HistoriaID, repo)
	if err != nil {
		return err
	}
	return nil
}
