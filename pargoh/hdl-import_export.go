package main

import (
	"monorepo/dhistorias"
	"monorepo/sqliteust"
	"strings"

	"github.com/pargomx/gecko"
	"github.com/pargomx/gecko/gko"
)

func (s *servidor) exportarJSON(c *gecko.Context) error {
	out, err := dhistorias.GetProyectoExport(c.PathVal("proyecto_id"), s.repo)
	if err != nil {
		return err
	}
	return c.JSON(200, out)
}

func (s *servidor) importarJSON(c *gecko.Context) error {
	proyecto := dhistorias.ProyectoExport{}
	err := c.JSONUnmarshalFile("proyecto", &proyecto)
	if err != nil {
		return gko.Err(err).Op("Unmarshall")
	}
	tx, err := s.db.Begin()
	if err != nil {
		return gko.Err(err).Op("Begin")
	}
	err = dhistorias.Importar(proyecto, sqliteust.NuevoRepo(tx))
	if err != nil {
		tx.Rollback()
		return gko.Err(err).Op("Importar")
	}
	err = tx.Commit()
	if err != nil {
		return gko.Err(err).Op("Commit")
	}
	return c.RefreshHTMX()
}

func (s *servidor) exportarMarkdown(c *gecko.Context) error {
	c.Response().WriteHeader(200)
	c.Response().Header().Set("Content-Type", "text/markdown")
	err := dhistorias.ExportarMarkdown(c.PathVal("proyecto_id"), c.Response().Writer, s.repo)
	if err != nil {
		return err
	}
	return nil
}

func (s *servidor) exportarDocx(c *gecko.Context) error {
	err := dhistorias.ExportarDocx(c.PathVal("proyecto_id"), s.repo, "export.docx")
	if err != nil {
		return err
	}
	return c.StatusOk("Exportación realizada")
}

func (s *servidor) exportarPDF(c *gecko.Context) error {
	tex, err := dhistorias.GetProyectoLaTeX(s.repo, c.PathVal("proyecto_id"))
	if err != nil {
		return err
	}
	pdf, err := tex.ToPDF()
	if err != nil {
		gko.LogError(err)
		return c.String(500, tex.String())
	}
	c.Response().Header().Set("Content-Type", "application/pdf")
	c.Response().Header().Set("Content-Disposition", "inline; filename=\"document.pdf\"")
	return c.Blob(200, "application/pdf", pdf)
}

func (s *servidor) exportarProyectoTeX(c *gecko.Context) error {
	tex, err := dhistorias.GetProyectoLaTeX(s.repo, c.PathVal("proyecto_id"))
	if err != nil {
		return err
	}
	return c.StringOk(tex.String())
}

func (s *servidor) exportarPersonaPDF(c *gecko.Context) error {
	tex, err := dhistorias.GetPersonaLaTeX(s.repo, c.PathInt("persona_id"))
	if err != nil {
		return err
	}
	pdf, err := tex.ToPDF()
	if err != nil {
		gko.LogError(err)
		return c.String(500, tex.String())
	}
	c.Response().Header().Set("Content-Type", "application/pdf")
	c.Response().Header().Set("Content-Disposition", "inline; filename=\"document.pdf\"")
	return c.Blob(200, "application/pdf", pdf)
}

// ================================================================ //

func (s *servidor) exportarArbolTXT(c *gecko.Context) error {
	proyectos, err := dhistorias.GetArbolCompleto(s.repo)
	if err != nil {
		return err
	}
	res := ""
	for _, pry := range proyectos {
		res += "\n" + pry.Proyecto.Titulo + "\n"
		for _, per := range pry.Personas {
			res += "\n" + per.Persona.Nombre + "\n"
			for _, his := range per.Historias {
				res += printHistRec(his, 1)
			}
		}
	}
	return c.StatusOk(res)
}
func printHistRec(his dhistorias.HistoriaExport, nivel int) string {
	res := strings.Repeat(" ", nivel) + "-" + his.Historia.Titulo + "\n"
	for _, hijo := range his.Historias {
		res += printHistRec(hijo, nivel+1)
	}
	return res
}
