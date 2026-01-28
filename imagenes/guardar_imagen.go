package imagenes

import (
	"bytes"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/disintegration/imaging"
	"github.com/pargomx/gecko/gko"
	"github.com/rwcarlsen/goexif/exif"
)

type ImagenEntrante struct {
	Reader  io.Reader // Archivo de entrada.
	Formato string    // Aceptados: "jpg", "jpeg", "png", "gif".
	ModTime time.Time // Posible fecha de captura original (modtime). Para JPEG se utiliza metadata.
}

// Escribe o sobreescribe una imagen en un archivo con el formato especificado.
// Retorna la ruta de la imagen guardada.
func (s *ImgService) GuardarImagen(imagenID string, input ImagenEntrante) (string, error) {
	op := gko.Op("GuardarImagen")
	err := ValidarImagenID(imagenID)
	if err != nil {
		return "", op.Err(err)
	}
	if input.Formato == "jpeg" {
		input.Formato = "jpg"
	}
	ima := Imagen{
		ImagenID:    imagenID,
		FechaUpload: time.Now(),
	}

	// Buffer para poder leer la imagen más de una vez.
	buf, err := io.ReadAll(input.Reader)
	if err != nil {
		return "", op.Err(err).Str("leer imagen")
	}

	// Fecha de captura original.
	if input.Formato == "jpg" {
		meta, err := exif.Decode(bytes.NewReader(buf))
		if err != nil {
			ima.FechaCaptura = input.ModTime // Fallback.
			gko.LogWarnf("Usando modtime para imagen jpg sin metadata: %v %v", err, imagenID)
			// return "",op.Str("leer metadata").Err(err).Ctx("img_bytes", len(buf))
		} else {
			ima.FechaCaptura, err = meta.DateTime()
			gko.LogWarnf("leer fecha captura: %v %v", err, imagenID)
			// if err != nil {
			// 	return "", op.Str("leer fecha captura").Err(err)
			// }
		}
	} else {
		ima.FechaCaptura = input.ModTime
	}
	if ima.FechaCaptura.IsZero() {
		ima.FechaCaptura = ima.FechaUpload
	}

	// Decodificar imagen para dimensiones y recodificar.
	imgReader := bytes.NewReader(buf)

	// img, _, err := image.Decode(imgReader)
	img, err := imaging.Decode(imgReader, imaging.AutoOrientation(true))
	if err != nil {
		return "", op.Key(gko.ErrNoSoportado).Err(err).Msg("Imposible decodificar imagen")
	}
	if img.Bounds().Max.X > s.maxPix {
		return "", op.Key(gko.ErrTooBig).Msgf("Suba máximo una imagen de %dpx, no %dpx", s.maxPix, img.Bounds().Max.X)
	}
	if img.Bounds().Max.Y > s.maxPix {
		return "", op.Key(gko.ErrTooBig).Msgf("Suba máximo una imagen de %dpx, no %dpx", s.maxPix, img.Bounds().Max.Y)
	}

	// Codificar en nuevo archivo para quitar metadata y a veces reducir tamaño.
	ima.Ruta = ima.ImagenID + "." + input.Formato
	outputPath := filepath.Join(s.directorio, ima.Ruta)
	outFile, err := os.OpenFile(outputPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0640)
	if err != nil {
		return "", op.Err(err).Str("abrir nuevo archivo").Ctx("path", outputPath)
	}
	defer outFile.Close()

	// Codificar y escribir.
	switch input.Formato {
	case "jpg":
		err = jpeg.Encode(outFile, img, nil)
	case "png":
		err = png.Encode(outFile, img)
	case "gif":
		err = gif.Encode(outFile, img, nil)
	default:
		os.Remove(outputPath)
		return "", op.Key(gko.ErrNoSoportado).Msgf("Formato de imagen no soportado: %v", input.Formato)
	}
	if err != nil {
		os.Remove(outputPath)
		return "", op.Err(err).Strf("codificar archivo %v", input.Formato)
	}
	outFile.Close()

	info, err := os.Stat(outputPath)
	if err != nil {
		os.Remove(outputPath)
		return "", op.Err(err).Str("obtener información de archivo recién creado")
	}
	ima.FileSize = info.Size()

	// El tiempo de modificación de la imagen es la fecha de captura.
	err = os.Chtimes(outputPath, time.Time{}, ima.FechaCaptura)
	if err != nil {
		gko.LogError(op.Err(err).Msg("Cambiar fecha de modificación de la imagen"))
	}

	// Hacer miniatura de img.
	err = s.hacerMiniatura(img, &ima)
	if err != nil {
		os.Remove(outputPath)
		return "", op.Err(err).Msg("No se puede hacer miniatura")
	}

	// Guardar en base de datos.
	err = s.repo.InsertImagen(ima)
	if err != nil {
		os.Remove(outputPath)
		return "", op.Err(err).Msg("No se puede guardar imagen en base de datos")
	}

	return ima.Ruta, nil
}

const IMG_MINIATURA_MAX = 256

func (s *ImgService) hacerMiniatura(img image.Image, ima *Imagen) error {
	op := gko.Op("HacerMiniatura")

	// Dimensiones de la imagen original.
	imgW, imgH := img.Bounds().Max.X, img.Bounds().Max.Y
	if imgW == 0 || imgH == 0 {
		return op.Str("imagen vacía")
	}

	// Redimensionar imagen
	newW, newH := IMG_MINIATURA_MAX, IMG_MINIATURA_MAX
	if imgW > imgH {
		newH = int(float64(newW) * float64(imgH) / float64(imgW))
	} else if imgW < imgH {
		newW = int(float64(newH) * float64(imgW) / float64(imgH))
	}
	mini := imaging.Resize(img, newW, newH, imaging.Lanczos)

	// Preparar archivo
	ima.Miniatura = ima.ImagenID + "_256.jpg"
	outputPath := filepath.Join(s.directorio, ima.Miniatura)
	outFile, err := os.OpenFile(outputPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0640)
	if err != nil {
		return op.Err(err).Str("abrir nuevo archivo").Ctx("path", outputPath)
	}
	defer outFile.Close()

	if ima.Miniatura == "" {
		return op.Str("ruta de miniatura vacía")
	}
	if ima.Miniatura == ima.Ruta {
		return op.Str("ruta de miniatura igual a imagen original")
	}

	// Codificar y escribir
	err = jpeg.Encode(outFile, mini, nil)
	if err != nil {
		os.Remove(outputPath)
		return op.Err(err).Str("codificar archivo")
	}
	return nil
}

func (s *ImgService) MakeMiniaturasFlatantes() error {
	op := gko.Op("MakeMiniaturasFlatantes")
	imgs, err := s.repo.ListImagenes()
	if err != nil {
		return op.Err(err).Msg("No se pueden obtener imágenes")
	}
	for _, ima := range imgs {
		if ima.Miniatura != "" {
			continue
		}
		imaFile, err := os.Open(filepath.Join(s.directorio, ima.Ruta))
		if err != nil {
			return op.Err(err).Str("abrir imagen").Ctx("imagen_id", ima.ImagenID)
		}
		img, _, err := image.Decode(imaFile)
		if err != nil {
			imaFile.Close()
			return op.Err(err).Str("decodificar imagen").Ctx("imagen_id", ima.ImagenID)
		}
		err = s.hacerMiniatura(img, &ima)
		if err != nil {
			imaFile.Close()
			return op.Err(err).Str("hacer miniatura").Ctx("imagen_id", ima.ImagenID)
		}
		imaFile.Close()
		err = s.repo.UpdateImagen(ima)
		if err != nil {
			return op.Err(err).Str("actualizar imagen").Ctx("imagen_id", ima.ImagenID)
		}
		gko.LogInfof("Miniatura creada: %v", ima.Miniatura)
	}
	return nil
}

func (s *ImgService) Remover(imgRuta string) error {
	op := gko.Op("imagenes.Remover")
	err := os.Remove(filepath.Join(s.directorio, imgRuta))
	if err != nil {
		return op.Err(err)
	}
	gko.LogInfof("Imagen eliminada %v", imgRuta)
	// TODO: Remover también la miniatura.
	return nil
}
