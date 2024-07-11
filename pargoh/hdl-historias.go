package main

import (
	"monorepo/historias_de_usuario/dhistorias"
	"monorepo/historias_de_usuario/ust"
	"strings"

	"github.com/pargomx/gecko"
)

func (s *servidor) getHistoriasLista(c *gecko.Context) error {
	agg, err := dhistorias.GetHistoriasDePadre(c.PathInt("nodo_id"), s.repo)
	if err != nil {
		return err
	}
	titulo := "Nodo"
	if agg.Abuelo != nil {
		titulo = "Historia" // agg.Abuelo.Titulo
	} else {
		titulo = agg.Persona.Nombre
	}
	data := map[string]any{
		"Titulo":   titulo,
		"Agregado": agg,
	}
	return c.RenderOk("hist_lista", data)
}

func (s *servidor) getHistoriasTablero(c *gecko.Context) error {
	agg, err := dhistorias.GetHistoriasDePadre(c.PathInt("nodo_id"), s.repo)
	if err != nil {
		return err
	}
	titulo := "Nodo"
	if agg.Abuelo != nil {
		titulo = "Historia" // agg.Abuelo.Titulo
	} else {
		titulo = agg.Persona.Nombre
	}
	data := map[string]any{
		"Titulo":   titulo,
		"Agregado": agg,
	}
	return c.RenderOk("hist_tablero", data)
}

func (s *servidor) getHistoriasPrioritarias(c *gecko.Context) error {
	Historias, err := s.repo.ListHistoriasPrioritarias()
	if err != nil {
		return err
	}
	data := map[string]any{
		"Titulo":    "Historias prioritarias",
		"Historias": Historias,
	}
	return c.RenderOk("hist_prioritarias", data)
}

func (s *servidor) formHistoria(c *gecko.Context) error {
	historia, err := s.repo.GetHistoria(c.PathInt("historia_id"))
	if err != nil {
		return err
	}
	data := map[string]any{
		"Titulo":   historia.Titulo,
		"Historia": historia,
	}
	return c.RenderOk("hist_form", data)
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

func (s *servidor) getTareasDeHistoria(c *gecko.Context) error {
	historia, err := s.repo.GetNodoHistoria(c.PathInt("historia_id"))
	if err != nil {
		return err
	}
	tareas, err := s.repo.ListTareasByHistoriaID(historia.HistoriaID)
	if err != nil {
		return err
	}
	agg, err := dhistorias.GetHistoriasDePadre(historia.HistoriaID, s.repo)
	if err != nil {
		return err
	}
	data := map[string]any{
		"Titulo":   "Tareas",
		"Historia": historia,
		"Tareas":   tareas,
		"Agregado": agg,

		"ListaTipoTarea": ust.ListaTipoTarea,
	}
	return c.RenderOk("hist_tareas", data)
}

func (s *servidor) getArbolCompleto(c *gecko.Context) error {
	arboles, err := dhistorias.GetArbolCompleto(s.repo)
	if err != nil {
		return err
	}
	res := "HISTORIAS DE USUARIO\n"
	for _, arbol := range arboles {
		res += "\n" + arbol.Persona.Nombre + "\n"
		for _, his := range arbol.Historias {
			res += printHistRec(his, 1)
		}
	}
	return c.StatusOk(res)
}
func printHistRec(his dhistorias.HistoriaRecursiva, nivel int) string {
	res := strings.Repeat(" ", nivel) + "-" + his.Padre.Titulo + "\n"
	for _, hijo := range his.Hijos {
		res += printHistRec(hijo, nivel+1)
	}
	return res
}
