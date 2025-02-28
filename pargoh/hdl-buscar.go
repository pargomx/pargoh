package main

import (
	"github.com/pargomx/gecko"
)

func (s *servidor) buscar(c *gecko.Context) error {
	resultados, err := s.repo.FullTextSearch(c.QueryVal("q"))
	if err != nil {
		return err
	}
	data := map[string]any{
		"Titulo":     "Búsqueda",
		"Resultados": resultados,
	}
	return c.RenderOk("busqueda", data)
}
