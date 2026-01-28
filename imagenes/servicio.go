package imagenes

import (
	"io/fs"
	"os"

	"github.com/pargomx/gecko/gko"
)

type ImgService struct {
	directorio string // Directorio absoluto donde se guardan las imágenes.
	FS         fs.FS  // Donde se leen los archivos.
	maxPix     int    // Máximo tamaño de imagen permitido.
	repo       ImgRepo
}

type ImgRepo interface {
	InsertImagen(Imagen) error
	UpdateImagen(Imagen) error
	ListImagenes() ([]Imagen, error)
}

// NewImgService verifica que se pueda escribir en el directorio para guardar imágenes.
func NewImgService(directorio string, maxPix int, repo ImgRepo) (*ImgService, error) {
	if maxPix <= 0 {
		return nil, gko.ErrDatoInvalido.Msg("Tamaño máximo de imagen debe ser mayor a 0px")
	}
	if directorio == "" {
		return nil, gko.ErrDatoIndef.Msg("directorio para imágenes indefinido")
	}
	if directorio == "/" {
		return nil, gko.ErrDatoIndef.Msg("directorio para imágenes no puede ser '/'")
	}
	stat, err := os.Stat(directorio)
	if err != nil {
		if os.IsNotExist(err) {
			err := os.MkdirAll(directorio, 0750)
			if err != nil {
				return nil, gko.Err(err).Msg("No se puede crear directorio para imágenes")
			}
			gko.LogInfof("Directorio de imágenes creado: %v", directorio)
		} else {
			return nil, gko.Err(err).Msg("No se puede acceder a directorio de imágenes")
		}
	} else {
		if !stat.IsDir() {
			return nil, gko.ErrDatoInvalido.Msgf("No es un directorio válido para imágenes: %v", directorio)
		}
	}
	outFile, err := os.CreateTemp(directorio, "*.jpeg")
	if err != nil {
		return nil, gko.Err(err).Msg("No se puede guardar imágenes")
	}
	defer outFile.Close()
	_, err = outFile.WriteString("prueba")
	if err != nil {
		return nil, gko.Err(err).Msg("No se puede escribir imágenes")
	}
	err = os.Remove(outFile.Name())
	if err != nil {
		return nil, gko.Err(err).Msg("No se puede eliminar imágenes")
	}
	if repo == nil {
		return nil, gko.ErrDatoIndef.Msg("ImgService sin repositorio")
	}
	return &ImgService{
		directorio: directorio,
		FS:         os.DirFS(directorio),
		maxPix:     maxPix,
		repo:       repo,
	}, nil
}

type RepoMock struct {
}

func (r RepoMock) InsertImagen(img Imagen) error {
	gko.LogInfof("Nueva imagen %v %v", img.Ruta, img.Miniatura)
	return nil
}
func (r RepoMock) UpdateImagen(Imagen) error {
	return nil
}
func (r RepoMock) ListImagenes() ([]Imagen, error) {
	return []Imagen{}, nil
}
