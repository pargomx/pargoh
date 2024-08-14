package main

import (
	"monorepo/dhistorias"
	"monorepo/sqliteust"

	"github.com/pargomx/gecko"
	"github.com/pargomx/gecko/gko"
)

func (s *servidor) exportarJSON(c *gecko.Context) error {
	out, err := dhistorias.ExportarProyecto(c.PathVal("proyecto_id"), s.repo)
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

func (s *servidor) exportarFile(c *gecko.Context) error {
	err := dhistorias.ExportarDocx(c.PathVal("proyecto_id"), s.repo, "export.docx")
	if err != nil {
		return err
	}
	return c.StatusOk("Exportación realizada")
}
