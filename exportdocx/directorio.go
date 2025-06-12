package exportdocx

import (
	"os"

	"github.com/pargomx/gecko/gko"
)

// Verificar que se pueda escribir en el directorio.
// Si no existe lo intenta crear.
func VerificarDirectorioExports(exportDir string) error {
	if exportDir == "" {
		return gko.ErrDatoIndef.Msg("Directorio para guardar exports indefinido")
	}
	inf, err := os.Stat(exportDir)
	if err != nil {
		if os.IsNotExist(err) {
			err := os.MkdirAll(exportDir, 0750)
			if err != nil {
				return gko.Err(err).Msg("No se puede crear directorio para exports")
			}
			gko.LogInfof("Directorio de exports creado: %v", exportDir)
		} else {
			return gko.Err(err).Msg("No se puede acceder a directorio de exports")
		}
	} else {
		if !inf.IsDir() {
			return gko.ErrDatoInvalido.Msgf("No es un directorio válido para exports: %v", exportDir)
		}
	}
	outFile, err := os.CreateTemp(exportDir, "*.doc")
	if err != nil {
		return gko.Err(err).Msg("No se puede guardar exports")
	}
	defer outFile.Close()
	_, err = outFile.WriteString("prueba")
	if err != nil {
		return gko.Err(err).Msg("No se puede escribir exports")
	}
	err = os.Remove(outFile.Name())
	if err != nil {
		return gko.Err(err).Msg("No se puede eliminar exports")
	}
	return nil
}
