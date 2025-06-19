package main

import (
	"fmt"
	"monorepo/arbol"
	"monorepo/dhistorias"

	"github.com/pargomx/gecko"
	"github.com/pargomx/gecko/gko"
)

func (s *writehdl) deleteProyecto(c *gecko.Context, tx *handlerTx) error {
	_, err := tx.app.EliminarNodo(arbol.ArgsEliminarNodo{
		NodoID: c.PathInt("proyecto_id"),
	})
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
	_, err = tx.app.EliminarRama(arbol.ArgsEliminarNodo{
		NodoID: pry.ProyectoID,
	})
	if err != nil {
		return err
	}
	return c.RedirOtro("/")
}

func (s *writehdl) deletePersona(c *gecko.Context, tx *handlerTx) error {
	_, err := tx.app.EliminarNodo(arbol.ArgsEliminarNodo{
		NodoID: c.PathInt("persona_id"),
	})
	if err != nil {
		return err
	}
	return c.RefreshHTMX()
}

func (s *writehdl) deleteHistoria(c *gecko.Context, tx *handlerTx) error {
	padre, err := tx.app.EliminarNodo(arbol.ArgsEliminarNodo{
		NodoID: c.PathInt("historia_id"),
	})
	if err != nil {
		return err
	}

	// TODO: AskedFor
	if padre.EsHistoriaDeUsuario() {
		return c.RedirOtrof("/historias/%v", padre.NodoID)
	} else if padre.EsPersona() {
		return c.RedirOtrof("/personas/%v", padre.NodoID)
	} else {
		gko.LogWarnf("deleteHistoria: padre %v no es persona ni historia", padre.NodoID)
		return c.RedirOtro("/proyectos")
	}
}

func (s *writehdl) deleteRegla(c *gecko.Context, tx *handlerTx) error {
	_, err := tx.app.EliminarNodo(arbol.ArgsEliminarNodo{
		NodoID: c.PathInt("regla_id"),
	})
	if err != nil {
		return err
	}
	defer s.reloader.brodcastReload(c)
	// TODO: Solo enviar el fragmento.
	return c.RedirOtrof("/historias/%v", c.PathInt("historia_id"))
}

func (s *writehdl) eliminarTarea(c *gecko.Context, tx *handlerTx) error {
	padre, err := tx.app.EliminarNodo(arbol.ArgsEliminarNodo{
		NodoID: c.PathInt("tarea_id"),
	})
	if err != nil {
		return err
	}
	defer s.reloader.brodcastReload(c)
	return c.AskedForFallback("/historias/%v", padre.NodoID)
}

func (s *writehdl) deleteTramoDeViaje(c *gecko.Context, tx *handlerTx) error {
	_, err := tx.app.EliminarNodo(arbol.ArgsEliminarNodo{
		NodoID: c.PathInt("tramo_id"),
	})
	if err != nil {
		return err
	}
	defer s.reloader.brodcastReload(c)
	return c.RedirOtrof("/historias/%v", c.PathInt("historia_id"))
}

// ================================================================ //

func (s *servidor) deleteImagenTramo(c *gecko.Context) error {
	err := dhistorias.EliminarFotoTramo(c.PathInt("historia_id"), c.PathInt("posicion"), s.cfg.imagesDir, s.repoOld)
	if err != nil {
		return err
	}
	defer s.reloader.brodcastReload(c)
	return c.RedirOtrof("/historias/%v", c.PathInt("historia_id"))
}
