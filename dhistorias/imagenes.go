package dhistorias

import (
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"

	"github.com/pargomx/gecko/gko"
)

// Escribe una imagen en un archivo con el formato especificado.
func GuardarImagen(input io.Reader, format string, outputPath string, maxPix int) error {
	img, _, err := image.Decode(input)
	if err != nil {
		return gko.ErrNoSoportado().Msg("Imposible decodificar imagen").Err(err)
	}
	if img.Bounds().Max.X > maxPix {
		return gko.ErrTooBig().Msgf("Suba máximo una imagen de %dpx, no %dpx", maxPix, img.Bounds().Max.X)
	}
	if img.Bounds().Max.Y > maxPix {
		return gko.ErrTooBig().Msgf("Suba máximo una imagen de %dpx, no %dpx", maxPix, img.Bounds().Max.Y)
	}

	// Crear nuevo archivo vacío
	outFile, err := os.OpenFile(outputPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0640)
	if err != nil {
		return gko.Err(err).Msg("Crear archivo para imagen").Ctx("path", outputPath)
	}
	defer outFile.Close()

	// Codificar y escribir.
	switch format {
	case "jpeg", "jpg":
		err = jpeg.Encode(outFile, img, nil)
	case "png":
		err = png.Encode(outFile, img)
	case "gif":
		err = gif.Encode(outFile, img, nil)
	default:
		return gko.ErrNoSoportado().Msgf("Formato de imagen no soportado: %v", format)
	}
	if err != nil {
		return gko.Err(err).Msgf("codificar archivo %v", format)
	}
	return nil
}

// ================================================================ //

func SetFotoTramo(HistoriaID int, Posicion int, foto io.Reader, directorio string, repo Repo) error {
	Tramo, err := repo.GetTramo(HistoriaID, Posicion)
	if err != nil {
		return err
	}
	Filename := fmt.Sprintf("h_%d_%d.jpeg", HistoriaID, Posicion)
	Filepath := filepath.Join(directorio, Filename)

	if Tramo.Imagen != "" { // Mover archivo anterior a _trash
		os.Rename(Filepath, filepath.Join(directorio, "trash_"+Filename))
	}
	err = GuardarImagen(foto, "jpeg", Filepath, 3000)
	if err != nil {
		return err
	}
	Tramo.Imagen = Filename
	err = repo.UpdateTramo(*Tramo)
	if err != nil {
		return err
	}
	gko.LogInfof("Imagen nueva %v", Tramo.Imagen)
	return nil
}

func EliminarFotoTramo(HistoriaID int, Posicion int, directorio string, repo Repo) error {
	Tramo, err := repo.GetTramo(HistoriaID, Posicion)
	if err != nil {
		return err
	}
	if Tramo.Imagen == "" {
		return gko.ErrDatoInvalido().Msg("No hay imagen que eliminar")
	}
	err = os.Remove(filepath.Join(directorio, Tramo.Imagen))
	if err != nil {
		return gko.Err(err).Op("EliminarFotoTramo")
	}
	Tramo.Imagen = ""
	err = repo.UpdateTramo(*Tramo)
	if err != nil {
		return err
	}
	gko.LogInfof("Imagen eliminada %v", Tramo.Imagen)
	return nil
}

// ================================================================ //
// ================================================================ //

// Esta función es la chida.
func SetImagenProyecto(proyectoID string, format string, input io.Reader, dir string, repo Repo) error {
	pry, err := repo.GetProyecto(proyectoID)
	if err != nil {
		return err
	}
	if format != "jpeg" && format != "png" && format != "gif" {
		return gko.ErrNoSoportado().Msgf("Formato de imagen no soportado: %v", format)
	}
	Filename := fmt.Sprintf("p_%s.%s", proyectoID, format) // TODO: prefijo único para evitar imágenes en cache.
	Filepath := filepath.Join(dir, Filename)
	if pry.Imagen != "" { // Mover archivo anterior en lugar de borrarlo.
		err = os.Rename(filepath.Join(dir, pry.Imagen), filepath.Join(dir, "trash_"+pry.Imagen))
		if err != nil {
			gko.Err(err).Op("SetImagenProyecto.TrashOld").Log()
		}
	}
	err = GuardarImagen(input, format, Filepath, 3000)
	if err != nil {
		return err
	}
	pry.Imagen = Filename
	err = repo.UpdateProyecto(*pry)
	if err != nil {
		return err
	}
	gko.LogInfof("Imagen nueva %v", pry.Imagen)
	return nil
}
