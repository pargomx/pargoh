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
	"github.com/pargomx/gecko/gkoid"
)

// Escribe una imagen en un archivo con el formato especificado.
func GuardarImagen(input io.Reader, format string, outputPath string, maxPix int) error {
	op := gko.Op("GuardarImagen")
	img, _, err := image.Decode(input)
	if err != nil {
		return op.Err(err).ErrNoSoportado().Msg("Imposible decodificar imagen")
	}
	if img.Bounds().Max.X > maxPix {
		return op.ErrTooBig().Msgf("Suba máximo una imagen de %dpx, no %dpx", maxPix, img.Bounds().Max.X)
	}
	if img.Bounds().Max.Y > maxPix {
		return op.ErrTooBig().Msgf("Suba máximo una imagen de %dpx, no %dpx", maxPix, img.Bounds().Max.Y)
	}

	// Evitar sobreescribir un archivo existente.
	_, err = os.Stat(outputPath)
	if err == nil {
		return op.ErrYaExiste().Strf("ya existe una imagen en %s y no se va a sobreescribir", outputPath)
	} else if !os.IsNotExist(err) {
		return op.Err(err).Str("error verificando existencia del archivo")
	}

	// Crear nuevo archivo vacío
	outFile, err := os.OpenFile(outputPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0640)
	if err != nil {
		return op.Err(err).Str("crear archivo para imagen").Ctx("path", outputPath)
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
		return op.ErrNoSoportado().Msgf("Formato de imagen no soportado: %v", format)
	}
	if err != nil {
		return op.Err(err).Strf("codificar archivo %v", format)
	}
	return nil
}

// ================================================================ //

func SetFotoTramo(HistoriaID int, Posicion int, foto io.Reader, directorio string, MIME string, repo Repo) error {
	op := gko.Op("SetFotoTramo")
	Tramo, err := repo.GetTramo(HistoriaID, Posicion)
	if err != nil {
		return op.Err(err)
	}
	extension := ""
	switch MIME {
	case "image/jpeg", "image/jpg":
		extension = "jpeg"
	case "image/png":
		extension = "png"
	case "image/gif":
		extension = "gif"
	default:
		return op.ErrNoSoportado().Msgf("MIME no soportado: %v", MIME)
	}

	uniqueID, err := gkoid.New62(8)
	if err != nil {
		return op.Err(err)
	}
	Filename := fmt.Sprintf("t_%s.%s", uniqueID, extension)
	Filepath := filepath.Join(directorio, Filename)

	if Tramo.Imagen != "" { // Mover archivo anterior a _trash
		os.Rename(Filepath, filepath.Join(directorio, "trash_"+Filename))
	}
	err = GuardarImagen(foto, extension, Filepath, 3000)
	if err != nil {
		return op.Err(err)
	}
	Tramo.Imagen = Filename
	err = repo.UpdateTramo(*Tramo)
	if err != nil {
		return op.Err(err)
	}
	gko.LogInfof("Imagen nueva %v", Tramo.Imagen)
	return nil
}

func EliminarFotoTramo(HistoriaID int, Posicion int, directorio string, repo Repo) error {
	op := gko.Op("EliminarFotoTramo")
	Tramo, err := repo.GetTramo(HistoriaID, Posicion)
	if err != nil {
		return op.Err(err)
	}
	if Tramo.Imagen == "" {
		return op.ErrDatoInvalido().Msg("No hay imagen que eliminar")
	}
	err = os.Remove(filepath.Join(directorio, Tramo.Imagen))
	if err != nil {
		return op.Err(err)
	}
	Tramo.Imagen = ""
	err = repo.UpdateTramo(*Tramo)
	if err != nil {
		return op.Err(err)
	}
	gko.LogInfof("Imagen eliminada %v", Tramo.Imagen)
	return nil
}

// ================================================================ //
// ================================================================ //

// Esta función es la chida.
func SetImagenProyecto(proyectoID string, format string, input io.Reader, dir string, repo Repo) error {
	op := gko.Op("SetImagenProyecto")
	pry, err := repo.GetProyecto(proyectoID)
	if err != nil {
		return op.Err(err)
	}
	if format != "jpeg" && format != "png" && format != "gif" {
		return op.ErrNoSoportado().Msgf("Formato de imagen no soportado: %v", format)
	}
	uniqueID, err := gkoid.New62(8)
	if err != nil {
		return op.Err(err)
	}
	Filename := fmt.Sprintf("p_%s.%s", uniqueID, format)
	Filepath := filepath.Join(dir, Filename)
	if pry.Imagen != "" { // Mover archivo anterior en lugar de borrarlo.
		err = os.Rename(filepath.Join(dir, pry.Imagen), filepath.Join(dir, "trash_"+pry.Imagen))
		if err != nil {
			op.Err(err).Op("SetImagenProyecto.TrashOld").Log()
		}
	}
	err = GuardarImagen(input, format, Filepath, 3000)
	if err != nil {
		return op.Err(err)
	}
	pry.Imagen = Filename
	err = repo.UpdateProyecto(*pry)
	if err != nil {
		return op.Err(err)
	}
	gko.LogInfof("Imagen nueva %v", pry.Imagen)
	return nil
}
