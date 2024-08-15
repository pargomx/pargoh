package dhistorias

import (
	"fmt"
	"image"
	"image/jpeg"
	_ "image/png"
	"io"
	"os"
	"path/filepath"

	"github.com/pargomx/gecko/gko"
)

// Se debe importar
//
//	_ "image/jpeg"
//	_ "image/png"

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

	img, _, err := image.Decode(foto)
	if err != nil {
		return gko.ErrNoSoportado().Msg("tipo de media no soportado").Err(err)
	}
	if img.Bounds().Max.X > 3000 {
		return gko.ErrTooBig().Msgf("Suba máximo una imagen de 3000px, no %vpx", img.Bounds().Max.X)
	}
	if img.Bounds().Max.Y > 3000 {
		return gko.ErrTooBig().Msgf("Suba máximo una imagen de 3000px, no %vpx", img.Bounds().Max.Y)
	}

	// Crear nuevo archivo vacío
	outFile, err := os.OpenFile(Filepath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0640)
	if err != nil {
		return gko.Err(err).Msg("no se puede crear archivo para foto").Ctx("directorio", directorio)
	}
	defer outFile.Close()

	// Encode into jpeg http://blog.golang.org/go-image-package
	err = jpeg.Encode(outFile, img, nil)
	if err != nil {
		return gko.Err(err).Msg("error al codificar archivo jpeg")
	}

	// Actualizar nombre de archivo en la base de datos.
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

func SetImagenProyecto(proyectoID string, foto io.Reader, directorio string, repo Repo) error {
	Proyecto, err := repo.GetProyecto(proyectoID)
	if err != nil {
		return err
	}
	Filename := fmt.Sprintf("p_%s.jpeg", proyectoID) // TODO: prefijo único para evitar imágenes en cache.
	Filepath := filepath.Join(directorio, Filename)

	if Proyecto.Imagen != "" { // Mover archivo anterior a _trash
		os.Rename(Filepath, filepath.Join(directorio, "trash_"+Filename))
	}

	img, _, err := image.Decode(foto)
	if err != nil {
		return gko.ErrNoSoportado().Msg("tipo de media no soportado").Err(err)
	}
	if img.Bounds().Max.X > 3000 {
		return gko.ErrTooBig().Msgf("Suba máximo una imagen de 3000px, no %vpx", img.Bounds().Max.X)
	}
	if img.Bounds().Max.Y > 3000 {
		return gko.ErrTooBig().Msgf("Suba máximo una imagen de 3000px, no %vpx", img.Bounds().Max.Y)
	}

	// Crear nuevo archivo vacío
	outFile, err := os.OpenFile(Filepath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0640)
	if err != nil {
		return gko.Err(err).Msg("no se puede crear archivo para foto").Ctx("directorio", directorio)
	}
	defer outFile.Close()

	// Encode into jpeg http://blog.golang.org/go-image-package
	err = jpeg.Encode(outFile, img, nil)
	if err != nil {
		return gko.Err(err).Msg("error al codificar archivo jpeg")
	}

	// Actualizar nombre de archivo en la base de datos.
	Proyecto.Imagen = Filename
	err = repo.UpdateProyecto(*Proyecto)
	if err != nil {
		return err
	}

	gko.LogInfof("Imagen nueva %v", Proyecto.Imagen)
	return nil
}
