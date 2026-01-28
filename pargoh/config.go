package main

import (
	"encoding/json"
	"os"

	"github.com/pargomx/gecko/gko"
)

// Si el archivo de configuración no existe se crea con la configuración default.
func getConfig(file string, cfg *configs) error {
	op := gko.Op("getConfig")
	if file == "" {
		return op.E(gko.ErrDatoIndef).Msg("archivo de configuración indefinido")
	}
	bytes, err := os.ReadFile(file)
	if len(bytes) == 0 {
		if err != nil && !os.IsNotExist(err) {
			return op.Err(err)
		}
		err = saveConfig(file, defaults)
		if err != nil {
			return op.Err(err)
		}
		*cfg = defaults
		return nil
	}
	err = json.Unmarshal(bytes, cfg)
	if err != nil {
		return op.Err(err)
	}
	return nil
}

// Crea el archivo de configuración en caso de que no exista.
func saveConfig(file string, cfg configs) error {
	op := gko.Op("saveConfig")
	bytes, err := json.MarshalIndent(&cfg, "", "\t")
	if err != nil {
		return op.Err(err)
	}

	// Para saber si ya existía el archivo.
	fileInfo, statErr := os.Stat(file)
	if statErr != nil && !os.IsNotExist(statErr) {
		return op.Err(statErr)
	}

	cfgFile, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return op.Err(err)
	}
	defer cfgFile.Close()
	_, err = cfgFile.Write(bytes)
	if err != nil {
		return op.Err(err)
	}

	if os.IsNotExist(statErr) || fileInfo.Size() == 0 {
		gko.LogInfo("Config file created: " + file)
	} else {
		gko.LogInfo("Config file updated: " + file)
	}
	return nil
}
