package main

import (
	"monorepo/dhistorias"

	"github.com/pargomx/gecko"
)

func (s *servidor) reordenarNodo(c *gecko.Context) error {
	err := dhistorias.ReordenarNodo(c.PathInt("nodo_id"), c.FormInt("newPosicion"), s.repo)
	if err != nil {
		return err
	}
	return c.StatusOkf("Nodo %v reordenado", c.PathInt("nodo_id"))
}
