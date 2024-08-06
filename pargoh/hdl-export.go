package main

import (
	"bytes"
	"monorepo/dhistorias"
	"os"

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

func (s *servidor) exportarFile(c *gecko.Context) error {
	buf := new(bytes.Buffer)
	err := dhistorias.ExportarMarkdown(buf, s.repo)
	if err != nil {
		return err
	}
	os.WriteFile("/home/andrew/proyectos/PARGO/pargoh/export.md", buf.Bytes(), 0644)

	err = dhistorias.ExportarDocx(s.repo, "/home/andrew/proyectos/PARGO/pargoh/export.docx")
	if err != nil {
		return err
	}
	return c.StatusOk("Exportación realizada")
}
