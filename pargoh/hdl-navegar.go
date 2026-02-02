package main

import (
	"monorepo/arbol"

	"github.com/pargomx/gecko"
)

func (s *readhdl) navDesdeRoot(c *gecko.Context) error {
	nodo, err := s.repo.GetNodoConArbol(arbol.NODO_ROOT)
	if err != nil {
		return err
	}
	data := map[string]any{
		"Nodo": nodo,
	}
	return c.RenderOk("nav_tree", data)
}

func (s *readhdl) navDesdeNodo(c *gecko.Context) error {
	nodo, err := s.repo.GetNodoConArbol(c.PathInt("nodo_id"))
	if err != nil {
		return err
	}
	data := map[string]any{
		"Nodo": nodo,
	}
	return c.RenderOk("nav_tree", data)
}
