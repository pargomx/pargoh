package imagenes

import (
	"github.com/pargomx/gecko/gko"
	"github.com/pargomx/gecko/gkoid"
)

const IMGAGEN_ID_LENGHT = 16

func NewImagenID() string {
	id, err := gkoid.New62(IMGAGEN_ID_LENGHT)
	if err != nil {
		gko.ErrInesperado.Op("NewImagenID").Err(err).Log()
		return ""
	}
	return id
}

// Debe ser alfanumérico con IMGAGEN_ID_LENGHT caracteres.
func ValidarImagenID(id string) error {
	if id == "" {
		return gko.ErrDatoIndef.Msg("ID de imagen indefinido")
	}
	if len(id) != IMGAGEN_ID_LENGHT {
		return gko.ErrDatoInvalido.Msg("ID de imagen debe tener 12 caracteres")
	}
	for _, r := range id {
		if !('0' <= r && r <= '9') && !('a' <= r && r <= 'z') && !('A' <= r && r <= 'Z') {
			return gko.ErrDatoInvalido.Msg("ID de imagen debe ser alfanumérico")
		}
	}
	return nil
}
