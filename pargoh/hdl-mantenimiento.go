package main

import (
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
