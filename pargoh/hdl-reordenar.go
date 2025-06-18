package main

import (
	"monorepo/arbol"

	"github.com/pargomx/gecko"
	"github.com/pargomx/gecko/gko"
)

func (s *servidor) reordenarPersona(c *gecko.Context, tx *handlerTx) error {
	err := tx.app.ReordenarEntidad(arbol.ArgsReordenar{
		NodoID: c.FormInt("persona_id"),
		NewPos: c.FormInt("new_pos"),
	})
	if err != nil {
		return err
	}
	pers, err := tx.repo.GetPersona(c.FormInt("persona_id"))
	if err != nil {
		return err
	}
	return c.RedirOtrof("/proyectos/%v", pers.ProyectoID)
}

func (s *servidor) reordenarHistoria(c *gecko.Context, tx *handlerTx) error {
	err := tx.app.ReordenarEntidad(arbol.ArgsReordenar{
		NodoID: c.FormInt("historia_id"),
		NewPos: c.FormInt("new_pos"),
	})
	if err != nil {
		return err
	}

	his, err := tx.repo.GetHistoria(c.FormInt("historia_id"))
	if err != nil {
		return err
	}
	defer s.reloader.brodcastReload(c)

	padre, err := tx.repo.GetNodo(his.PadreID)
	if err != nil {
		return err
	}
	if padre.EsPersona() {
		return c.RedirOtrof("/personas/%v", his.PadreID)
	} else if padre.EsHistoriaDeUsuario() {
		return c.RedirOtrof("/historias/%v", his.PadreID)
	}
	return gko.ErrInesperado.Msgf("reordenarHistoria: padre %v no es persona ni historia, sino %v",
		padre.NodoID, padre.Tipo)
}

func (s *servidor) reordenarTramo(c *gecko.Context, tx *handlerTx) error {
	err := tx.app.ReordenarEntidad(arbol.ArgsReordenar{
		NodoID: c.FormInt("historia_id"),
		NewPos: c.FormInt("new_pos"),
	})
	if err != nil {
		return err
	}
	defer s.reloader.brodcastReload(c)
	return c.RedirOtrof("/historias/%v", c.FormInt("historia_id"))
}

func (s *servidor) reordenarRegla(c *gecko.Context, tx *handlerTx) error {
	err := tx.app.ReordenarEntidad(arbol.ArgsReordenar{
		NodoID: c.FormInt("historia_id"),
		NewPos: c.FormInt("new_pos"),
	})
	if err != nil {
		return err
	}
	defer s.reloader.brodcastReload(c)
	// TODO: Solo enviar el fragmento.
	return c.RedirOtrof("/historias/%v", c.FormInt("historia_id"))
}
