package main

import (
	"fmt"
	"monorepo/dhistorias"

	"github.com/pargomx/gecko"
	"github.com/pargomx/gecko/gko"
)

func (s *writehdl) deleteProyecto(c *gecko.Context, tx *handlerTx) error {
	_, err := tx.app.EliminarNodo(c.PathInt("proyecto_id"))
	if err != nil {
		return err
	}
	return c.RedirOtro("/")
}

func (s *writehdl) deleteProyectoPorCompleto(c *gecko.Context, tx *handlerTx) error {
	pry, err := tx.repo.GetProyecto(c.PathInt("proyecto_id"))
	if err != nil {
		return err
	}
	if c.PromptVal() != fmt.Sprintf("eliminar_%v", pry.ProyectoID) {
		return gko.ErrDatoInvalido.Msg("No se confirmó la eliminación")
	}
	_, err = tx.app.EliminarRama(pry.ProyectoID)
	if err != nil {
		return err
	}
	return c.RedirOtro("/")
}

func (s *writehdl) deletePersona(c *gecko.Context, tx *handlerTx) error {
	_, err := tx.app.EliminarNodo(c.PathInt("persona_id"))
	if err != nil {
		return err
	}
	return c.RefreshHTMX()
}

func (s *writehdl) deleteHistoria(c *gecko.Context, tx *handlerTx) error {
	padre, err := tx.app.EliminarNodo(c.PathInt("historia_id"))
	if err != nil {
		return err
	}

	// TODO: AskedFor
	if padre.EsHistoriaDeUsuario() {
		return c.RedirOtrof("/h/%v", padre.NodoID)
	} else if padre.EsPersona() {
		return c.RedirOtrof("/h/%v", padre.NodoID)
	} else {
		gko.LogWarnf("deleteHistoria: padre %v no es persona ni historia", padre.NodoID)
		return c.RedirOtro("/h")
	}
}

func (s *writehdl) deleteRegla(c *gecko.Context, tx *handlerTx) error {
	_, err := tx.app.EliminarNodo(c.PathInt("regla_id"))
	if err != nil {
		return err
	}
	defer s.reloader.brodcastReload(c)
	// TODO: Solo enviar el fragmento.
	return c.RedirOtrof("/h/%v", c.PathInt("historia_id"))
}

func (s *writehdl) eliminarTarea(c *gecko.Context, tx *handlerTx) error {
	padre, err := tx.app.EliminarNodo(c.PathInt("tarea_id"))
	if err != nil {
		return err
	}
	defer s.reloader.brodcastReload(c)
	return c.AskedForFallback("/h/%v", padre.NodoID)
}

func (s *writehdl) deleteTramoDeViaje(c *gecko.Context, tx *handlerTx) error {
	_, err := tx.app.EliminarNodo(c.PathInt("tramo_id"))
	if err != nil {
		return err
	}
	defer s.reloader.brodcastReload(c)
	return c.RedirOtrof("/h/%v", c.PathInt("historia_id"))
}

// ================================================================ //

func (s *servidor) deleteImagenTramo(c *gecko.Context, tx *handlerTx) error {
	err := dhistorias.EliminarFotoTramo(c.PathInt("historia_id"), c.PathInt("posicion"), s.cfg.ImagesDir, s.repoOld)
	if err != nil {
		return err
	}
	defer s.reloader.brodcastReload(c)
	return c.RedirOtrof("/h/%v", c.PathInt("historia_id"))
}
