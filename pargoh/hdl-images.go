package main

import (
	"monorepo/dhistorias"
	"os"

	"github.com/pargomx/gecko"
	"github.com/pargomx/gecko/gko"
)

// Verificar que se pueda escribir en el directorio.
// Si no existe lo intenta crear.
func (s *servidor) verificarDirectorioImagenes() error {
	if s.cfg.imagesDir == "" {
		return gko.ErrDatoIndef().Msg("Directorio para guardar imágenes indefinido")
	}
	inf, err := os.Stat(s.cfg.imagesDir)
	if err != nil {
		if os.IsNotExist(err) {
			err := os.MkdirAll(s.cfg.imagesDir, 0750)
			if err != nil {
				return gko.Err(err).Msg("No se puede crear directorio para imágenes")
			}
			gko.LogInfof("Directorio de imágenes creado: %v", s.cfg.imagesDir)
		} else {
			return gko.Err(err).Msg("No se puede acceder a directorio de imágenes")
		}
	} else {
		if !inf.IsDir() {
			return gko.ErrDatoInvalido().Msgf("No es un directorio válido para imágenes: %v", s.cfg.imagesDir)
		}
	}
	outFile, err := os.CreateTemp(s.cfg.imagesDir, "*.jpeg")
	if err != nil {
		return gko.Err(err).Msg("No se puede guardar imágenes")
	}
	defer outFile.Close()
	_, err = outFile.WriteString("prueba")
	if err != nil {
		return gko.Err(err).Msg("No se puede escribir imágenes")
	}
	err = os.Remove(outFile.Name())
	if err != nil {
		return gko.Err(err).Msg("No se puede eliminar imágenes")
	}
	return nil
}

func (s *servidor) setImagenTramo(c *gecko.Context) error {
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
	err = dhistorias.SetFotoTramo(c.FormInt("historia_id"), c.FormInt("posicion"), foto, s.cfg.imagesDir, s.repo)
	if err != nil {
		return err
	}
	return c.RefreshHTMX()
}

func (s *servidor) setImagenProyecto(c *gecko.Context) error {
	hdr, err := c.FormFile("imagen")
	if err != nil {
		return err
	}
	file, err := hdr.Open()
	if err != nil {
		return err
	}
	defer file.Close()
	gko.LogDebugf("Imagen recibida: %v\t Tamaño: %v\t MIME:%v", hdr.Filename, hdr.Size, hdr.Header.Get("Content-Type"))
	err = dhistorias.SetImagenProyecto(c.PathVal("proyecto_id"), file, s.cfg.imagesDir, s.repo)
	if err != nil {
		return err
	}
	return c.RefreshHTMX()
}

func (s *servidor) deleteImagenTramo(c *gecko.Context) error {
	err := dhistorias.EliminarFotoTramo(c.PathInt("historia_id"), c.PathInt("posicion"), s.cfg.imagesDir, s.repo)
	if err != nil {
		return err
	}
	return c.RefreshHTMX()
}
