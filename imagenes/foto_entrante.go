package imagenes

import (
	"mime/multipart"
	"strings"
	"time"

	"github.com/pargomx/gecko"
	"github.com/pargomx/gecko/gko"
	"github.com/pargomx/gecko/gkt"
)

func GetImagenEntrante(c *gecko.Context, name string) (*ImagenEntrante, error) {
	if name == "" {
		return nil, gko.ErrDatoIndef.Str("especifique el campo del form").Msg("No disponible")
	}
	form, err := c.MultipartForm()
	if err != nil {
		return nil, gko.Err(err).Msg("No se puede leer formulario con archivos adjuntos")
	}
	fileHeaders := form.File[name]
	if len(fileHeaders) == 0 {
		return nil, gko.ErrDatoIndef.Msg("Se debe subir al menos una imagen")
	}
	hdr := fileHeaders[0]
	file, err := hdr.Open()
	if err != nil {
		return nil, gko.Err(err).Msg("No se puede abrir archivo adjunto")
	}
	formato := hdr.Header.Get("Content-Type")
	switch formato {
	case "image/jpeg", "image/jpg":
		formato = "jpeg"
	case "image/png":
		formato = "png"
	case "image/gif":
		formato = "gif"
	default:
		return nil, gko.ErrNoSoportado.Msgf("MIME no soportado: %v", formato)
	}
	imagen := &ImagenEntrante{
		Formato: formato,
		Reader:  file,
	}
	imagen.ModTime = gkt.Now()
	return imagen, nil
}

func GetImagenesEntrante(c *gecko.Context, name string) ([]ImagenEntrante, error) {
	if name == "" {
		return nil, gko.ErrDatoIndef.Str("especifique el campo del form").Msg("No disponible")
	}
	form, err := c.MultipartForm()
	if err != nil {
		return nil, gko.Err(err).Msg("No se puede leer formulario con archivos adjuntos")
	}
	fileHeaders := form.File[name]
	if len(fileHeaders) == 0 {
		return nil, gko.ErrDatoIndef.Msg("Se debe subir al menos una imagen")
	}
	fotos := []ImagenEntrante{}
	files := []multipart.File{}
	fileModTimes := strings.Split(c.FormVal("modtimes"), ",")
	for i, hdr := range fileHeaders {
		file, err := hdr.Open()
		if err != nil {
			return nil, gko.Err(err).Msg("No se puede abrir archivo adjunto")
		}
		files = append(files, file)
		formato := hdr.Header.Get("Content-Type")
		switch formato {
		case "image/jpeg", "image/jpg":
			formato = "jpeg"
		case "image/png":
			formato = "png"
		case "image/gif":
			formato = "gif"
		default:
			return nil, gko.ErrNoSoportado.Msgf("MIME no soportado: %v", formato)
		}
		imagen := ImagenEntrante{
			Formato: formato,
			Reader:  files[i],
		}
		if i < len(fileModTimes) {
			if fileModTimes[i] == "" {
				return nil, gko.ErrDatoIndef.Msg("No se incluyó la fecha de modificación del archivo")
			}
			imagen.ModTime, err = time.Parse("2006-01-02T15:04:05", fileModTimes[i])
			if err != nil {
				return nil, gko.Err(err).Msg("No se puede leer fecha de modificación de archivo")
			}
		}
		fotos = append(fotos, imagen)
	}
	return fotos, nil
}
