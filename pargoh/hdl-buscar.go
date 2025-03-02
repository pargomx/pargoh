package main

import (
	"html/template"
	"regexp"
	"strings"

	"github.com/pargomx/gecko"
)

var (
	reMatch = regexp.MustCompile(`ſ(.*?)ſ`)
)

func highlight(text string) template.HTML {
	escapedText := template.HTMLEscapeString(text)
	escapedText = reMatch.ReplaceAllStringFunc(escapedText, func(match string) string {
		return "<span class=\"text-cyan-400 font-bold\">" + strings.Trim(match, "ſ") + "</span>"
	})
	return template.HTML(escapedText)
}

func (s *servidor) buscar(c *gecko.Context) error {
	query := c.QueryVal("q")
	resultados, err := s.repo.FullTextSearch(query)
	if err != nil {
		return err
	}
	for i, r := range resultados {
		if strings.HasPrefix(r.Texto, "... ") {
			resultados[i].Texto = "..." + strings.TrimPrefix(r.Texto, "... ")
		}
		if strings.HasSuffix(r.Texto, " ...") {
			resultados[i].Texto = strings.TrimSuffix(r.Texto, " ...") + "..."
		}
		resultados[i].Subrallado = highlight(resultados[i].Texto)
	}
	data := map[string]any{
		"Titulo":     "Búsqueda",
		"Busqueda":   query,
		"Resultados": resultados,
	}
	return c.RenderOk("busqueda", data)
}
