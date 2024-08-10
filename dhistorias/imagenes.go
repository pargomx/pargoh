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
		os.Rename(Filepath, filepath.Join(directorio, ".trash_"+Filename))
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

	gko.LogInfof("Nueva foto %v", Tramo.Imagen)
	return nil
}
