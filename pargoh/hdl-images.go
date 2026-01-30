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
	if s.cfg.ImagesDir == "" {
		return gko.ErrDatoIndef.Msg("Directorio para guardar imágenes indefinido")
	}
	inf, err := os.Stat(s.cfg.ImagesDir)
	if err != nil {
		if os.IsNotExist(err) {
			err := os.MkdirAll(s.cfg.ImagesDir, 0750)
			if err != nil {
				return gko.Err(err).Msg("No se puede crear directorio para imágenes")
			}
			gko.LogInfof("Directorio de imágenes creado: %v", s.cfg.ImagesDir)
		} else {
			return gko.Err(err).Msg("No se puede acceder a directorio de imágenes")
		}
	} else {
		if !inf.IsDir() {
			return gko.ErrDatoInvalido.Msgf("No es un directorio válido para imágenes: %v", s.cfg.ImagesDir)
		}
	}
	outFile, err := os.CreateTemp(s.cfg.ImagesDir, "*.jpeg")
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

func (s *servidor) setImagenTramo(c *gecko.Context, tx *handlerTx) error {
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
	err = dhistorias.SetFotoTramo(
		c.FormInt("historia_id"),
		c.FormInt("posicion"),
		foto,
		s.cfg.ImagesDir,
		file.Header.Get("Content-Type"),
		s.repoOld,
	)
	if err != nil {
		return err
	}
	defer s.reloader.brodcastReload(c)
	return c.RedirOtrof("/h/%v", c.FormInt("historia_id"))
}
