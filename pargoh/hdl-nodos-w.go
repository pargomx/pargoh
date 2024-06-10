package main

import (
	"monorepo/gecko"
	"monorepo/historias_de_usuario/dhistorias"
)

func (s *servidor) reordenarNodo(c *gecko.Context) error {
	err := dhistorias.ReordenarNodo(c.PathInt("nodo_id"), c.FormInt("newPosicion"), s.repo)
	if err != nil {
		return err
	}
	return c.StatusOkf("Nodo %v reordenado", c.PathInt("nodo_id"))
}
