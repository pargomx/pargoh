package imagenes

import (
	"strings"

	"github.com/pargomx/gecko/gko"
)

type MiniaturaGetter struct{}

func (MiniaturaGetter) DerivarDe(path string) string {
	if path == "" {
		return ""
	}
	partes := strings.Split(path, ".")
	if len(partes) != 2 {
		gko.LogWarnf("img path no tiene UN punto: '%v'", path)
		return path
	}
	return partes[0] + "_256.jpg"
}
