package main

import (
	"monorepo/dhistorias"
	"os"

	"github.com/pargomx/gecko"
	"github.com/pargomx/gecko/gko"
)

func (s *servidor) putImagen(directorio string) gecko.HandlerFunc {
	// Verificar que podamos escribir en el directorio
	if directorio == "" {
		gko.LogWarn("Directorio para guardar fotos indefinido")
		return func(c *gecko.Context) error { return gko.ErrNoDisponible().Msg("Fotos no disponibles") }
	}
	// TODO: crear directorio si no existe...?
	outFile, err := os.CreateTemp(directorio, "*.jpeg")
	if err != nil {
		gko.LogWarn("No se podrán guardar fotos: " + err.Error())
		return func(c *gecko.Context) error { return gko.ErrNoDisponible().Msg("Fotos no disponibles") }
	}
	outFile.WriteString("prueba")
	defer outFile.Close()
	err = os.Remove(outFile.Name())
	if err != nil {
		gko.LogWarn("No se podrán eliminar fotos: " + err.Error())
		return func(c *gecko.Context) error { return gko.ErrNoDisponible().Msg("Fotos no disponibles") }
	}

	// Handler para recibir y guardar imagen.
	return func(c *gecko.Context) error {
		file, err := c.FormFile("imagen")
		if err != nil {
			return err
		}
		foto, err := file.Open()
		if err != nil {
			return err
		}
		defer foto.Close()
		// gko.LogDebugf("Imagen recibida: %v\t Tamaño: %v\t MIME:%v", file.Filename, file.Size, file.Header.Get("Content-Type"))
		err = dhistorias.SetFotoTramo(c.FormInt("historia_id"), c.FormInt("posicion"), foto, directorio, s.repo)
		if err != nil {
			return err
		}
		return c.RefreshHTMX()
	}
}

func (s *servidor) deleteImagen(directorio string) gecko.HandlerFunc {
	if directorio == "" {
		gko.LogWarn("Directorio para guardar fotos indefinido")
		return func(c *gecko.Context) error { return gko.ErrNoDisponible().Msg("Fotos no disponibles") }
	}
	return func(c *gecko.Context) error {
		err := dhistorias.EliminarFotoTramo(c.PathInt("historia_id"), c.PathInt("posicion"), directorio, s.repo)
		if err != nil {
			return err
		}
		return c.RefreshHTMX()
	}
}
