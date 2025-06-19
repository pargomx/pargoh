package arbol

import (
	"os"
	"path/filepath"

	"github.com/pargomx/gecko/gko"
)

const EvBorrarImagen gko.EventKey = "imagen_eliminada"

type argsBorrarImagen struct {
	Filename string
}

func (s *AppTx) borrarImagen(args argsBorrarImagen) error {
	op := gko.Op("BorrarImagen")
	if args.Filename == "" {
		return op.E(gko.ErrDatoIndef).Str("filename indefinido")
	}

	err := os.Remove(filepath.Join(s.ImagesDir, args.Filename))
	if err != nil {
		return op.Err(err)
	}

	s.Results.Add(EvBorrarImagen.WithArgs(args))

	return nil
}
