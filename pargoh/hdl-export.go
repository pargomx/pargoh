package main

import (
	"monorepo/historias_de_usuario/dhistorias"

	"github.com/pargomx/gecko"
)

func (s *servidor) exportarMarkdown(c *gecko.Context) error {
	c.Response().WriteHeader(200)
	c.Response().Header().Set("Content-Type", "text/markdown")
	err := dhistorias.ExportarMarkdown(c.Response().Writer, s.repo)
	if err != nil {
		return err
	}
	return nil
}
