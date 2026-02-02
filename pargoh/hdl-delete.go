package main

import (
	"fmt"
	"monorepo/dhistorias"

	"github.com/pargomx/gecko"
	"github.com/pargomx/gecko/gko"
)

func (s *writehdl) eliminarNodo(c *gecko.Context, tx *handlerTx) error {
	padre, err := tx.app.EliminarNodo(c.PathInt("nodo_id"))
	if err != nil {
		return err
	}
	return c.RedirOtrof("/h/%v", padre.NodoID)
}

func (s *writehdl) eliminarRama(c *gecko.Context, tx *handlerTx) error {
	nodoID, err := tx.repo.GetNodo(c.PathInt("nodo_id"))
	if err != nil {
		return err
	}
	if c.PromptVal() != fmt.Sprintf("eliminar_%v", nodoID.NodoID) {
		return gko.ErrDatoInvalido.Msg("No se confirmó la eliminación")
	}
	padre, err := tx.app.EliminarRama(nodoID.NodoID)
	if err != nil {
		return err
	}
	return c.RedirOtrof("/h/%v", padre.NodoID)
}

func (s *servidor) eliminarImagen(c *gecko.Context, tx *handlerTx) error {
	err := dhistorias.EliminarImagen(c.PathInt("historia_id"), c.PathInt("posicion"), s.cfg.ImagesDir, s.repoOld)
	if err != nil {
		return err
	}
	defer s.reloader.brodcastReload(c)
	return c.RedirOtrof("/h/%v", c.PathInt("historia_id"))
}
