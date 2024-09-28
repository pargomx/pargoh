package main

import (
	"monorepo/dhistorias"

	"github.com/pargomx/gecko"
)

func (s *servidor) materializarHistorias(c *gecko.Context) error {
	err := dhistorias.MaterializarAncestrosDeHistorias(s.repo)
	if err != nil {
		return err
	}
	return c.StringOk("Proyecto y persona materializados para historias")
}
