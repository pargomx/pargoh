package main

import (
	"bytes"
	"monorepo/historias_de_usuario/dhistorias"
	"os"

	"github.com/pargomx/gecko"
	"github.com/pargomx/gecko/gko"
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

func (s *servidor) exportarFile() {
	buf := new(bytes.Buffer)
	err := dhistorias.ExportarMarkdown(buf, s.repo)
	if err != nil {
		gko.FatalError(err)
	}
	os.WriteFile("/home/andrew/proyectos/PARGO/pargoh/export.md", buf.Bytes(), 0644)

	err = dhistorias.ExportarDocx(s.repo, "/home/andrew/proyectos/PARGO/pargoh/export.docx")
	if err != nil {
		gko.FatalError(err)
	}
}
