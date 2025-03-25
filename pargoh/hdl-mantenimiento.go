package main

import (
	"monorepo/dhistorias"

	"github.com/pargomx/gecko"
)

func (s *servidor) continuar(c *gecko.Context) error {
	if c.QueryBool("set") {
		s.noContinuar = !s.noContinuar
		return c.RedirFull("/continuar")
	}
	return c.Render(200, "app/continuar", !s.noContinuar)
}

func (s *servidor) offline(c *gecko.Context) error {
	return c.RenderOk("app/offline", nil)
}

func (s *servidor) materializarHistorias(c *gecko.Context) error {
	err := dhistorias.MaterializarAncestrosDeHistorias(s.repo)
	if err != nil {
		return err
	}
	return c.StringOk("Proyecto y persona materializados para historias")
}

func (s *servidor) moverHistoriaForm(c *gecko.Context) error {
	historia, err := s.repo.GetNodoHistoria(c.PathInt("historia_id"))
	if err != nil {
		return err
	}
	arboles, err := dhistorias.GetArbolCompleto(s.repo)
	if err != nil {
		return err
	}
	data := map[string]any{
		"Titulo":   "Mover historia",
		"Arboles":  arboles,
		"Historia": historia,
	}
	return c.RenderOk("hist_mover", data)
}
