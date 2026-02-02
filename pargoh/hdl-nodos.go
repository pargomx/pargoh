package main

import (
	"monorepo/arbol"

	"github.com/pargomx/gecko"
	"github.com/pargomx/gecko/gko"
)

func (s *readhdl) getRawNodoEditor(c *gecko.Context) error {
	nodo, err := s.repo.GetNodo(c.PathInt("nodo_id"))
	if err != nil {
		return err
	}
	return c.RenderOk("nodo-raw", map[string]any{
		"Titulo": nodo.Titulo,
		"Nodo":   nodo,
	})
}

func (s *writehdl) parcharNodo(c *gecko.Context, tx *handlerTx) error {
	propiedad := c.PathVal("param")
	if propiedad == "" {
		return gko.ErrDatoIndef.Msg("Propiedad a parchar indefinida")
	}
	err := tx.app.ParcharNodo(arbol.ArgsParcharNodo{
		NodoID: c.PathInt("nodo_id"),
		Campo:  propiedad,
		NewVal: c.FormValue(propiedad),
	})
	if err != nil {
		return err
	}
	defer s.reloader.brodcastReload(c)
	return c.AskedForFallback("/h/%v", c.PathInt("nodo_id"))
	// return c.AskedFor("Cambio guardado")
	// return c.RedirOtrof("/h/%v", c.PathInt("nodo_id"))
}
