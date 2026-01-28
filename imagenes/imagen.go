package imagenes

import (
	"time"
)

// Imagen corresponde a un elemento de la tabla 'imagenes'.
type Imagen struct {
	ImagenID     string    // `imagenes.imagen_id`  Alfanumérico de 16 caracteres.
	Ruta         string    // `imagenes.ruta`  Ruta al archivo de la imagen en resolución completa.
	FileSize     int64     // `imagenes.file_size`  Tamaño de la imagen en resolución completa en bytes.
	Miniatura    string    // `imagenes.miniatura`  Ruta al archivo miniatura de la imagen.
	FechaCaptura time.Time // `imagenes.fecha_captura`  Fecha en la que se capturó la imagen.
	FechaUpload  time.Time // `imagenes.fecha_upload`  Fecha en la que se subió el archivo.
}
