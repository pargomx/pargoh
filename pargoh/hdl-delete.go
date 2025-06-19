package main

import (
	"monorepo/dhistorias"

	"github.com/pargomx/gecko"
	"github.com/pargomx/gecko/gko"
)

func (s *servidor) deleteProyecto(c *gecko.Context) error {
	err := dhistorias.EliminarProyecto(c.PathVal("proyecto_id"), s.repoOld)
	if err != nil {
		return err
	}
	return c.RedirOtro("/")
}

func (s *servidor) deleteProyectoPorCompleto(c *gecko.Context) error {
	pry, err := dhistorias.GetProyectoExport(c.PathVal("proyecto_id"), s.repoOld)
	if err != nil {
		return err
	}
	if c.PromptVal() != "eliminar_"+pry.Proyecto.ProyectoID {
		return gko.ErrDatoInvalido.Msg("No se confirmó la eliminación")
	}
	err = pry.EliminarPorCompleto(s.repoOld)
	if err != nil {
		return err
	}
	return c.RedirOtro("/")
}

func (s *servidor) deletePersona(c *gecko.Context) error {
	err := dhistorias.EliminarPersona(c.PathInt("persona_id"), s.repoOld)
	if err != nil {
		return err
	}
	return c.RefreshHTMX()
}

func (s *servidor) deleteHistoria(c *gecko.Context) error {
	padreID, err := dhistorias.EliminarHistoria(c.PathInt("historia_id"), s.repoOld)
	if err != nil {
		return err
	}
	padre, err := s.repoOld.GetNodo(padreID)
	if err != nil {
		return err
	}
	// TODO: AskedFor
	if padre.EsHistoria() {
		return c.RedirOtrof("/historias/%v", padreID)
	} else if padre.EsPersona() {
		return c.RedirOtrof("/personas/%v", padreID)
	} else {
		gko.LogWarnf("deleteHistoria: padre %v no es persona ni historia", padreID)
		return c.RedirOtro("/proyectos")
	}
}

func (s *servidor) deleteRegla(c *gecko.Context) error {
	tx, err := s.newRepoTx()
	if err != nil {
		return err
	}
	err = dhistorias.EliminarRegla(tx.repoOld, c.PathInt("historia_id"), c.PathInt("posicion"))
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	defer s.reloader.brodcastReload(c)
	// TODO: Solo enviar el fragmento.
	return c.RedirOtrof("/historias/%v", c.PathInt("historia_id"))
}

func (s *servidor) deleteImagenTramo(c *gecko.Context) error {
	err := dhistorias.EliminarFotoTramo(c.PathInt("historia_id"), c.PathInt("posicion"), s.cfg.imagesDir, s.repoOld)
	if err != nil {
		return err
	}
	defer s.reloader.brodcastReload(c)
	return c.RedirOtrof("/historias/%v", c.PathInt("historia_id"))
}

func (s *servidor) eliminarTarea(c *gecko.Context) error {
	historiaID, err := dhistorias.EliminarTarea(c.PathInt("tarea_id"), s.repoOld)
	if err != nil {
		return err
	}
	gko.LogInfof("Tarea %d eliminada", c.PathInt("tarea_id"))
	defer s.reloader.brodcastReload(c)
	return c.AskedForFallback("/historias/%v", historiaID)
}

func (s *servidor) deleteTramoDeViaje(c *gecko.Context) error {
	tx, err := s.newRepoTx()
	if err != nil {
		return err
	}
	err = dhistorias.EliminarTramoDeViaje(tx.repoOld, c.PathInt("historia_id"), c.PathInt("posicion"))
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	defer s.reloader.brodcastReload(c)
	return c.RedirOtrof("/historias/%v", c.PathInt("historia_id"))
}
