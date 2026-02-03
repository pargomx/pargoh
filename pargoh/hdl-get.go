package main

import (
	"monorepo/dhistorias"

	"github.com/pargomx/gecko"
)

func (s *readhdl) getNodoTablero(c *gecko.Context) error {
	Historia, err := dhistorias.GetHistoria(c.PathInt("nodo_id"), dhistorias.GetDescendientes, s.repoOld)
	if err != nil {
		return err
	}
	data := map[string]any{
		"Titulo":   Historia.Historia.Titulo,
		"Agregado": Historia,
	}
	return c.RenderOk("hist_tablero", data)
}
