package dhistorias

import (
	"monorepo/ust"
	"regexp"
	"strings"

	"github.com/pargomx/gecko/gko"
	"github.com/pargomx/gecko/gkt"
)

func nuevoProyectoID(clave string) string {
	clave = strings.TrimSpace(regexp.MustCompile(`\s+`).ReplaceAllString(clave, ""))
	clave = strings.ReplaceAll(clave, "-", "_")
	return strings.ToLower(clave)
}

func NuevoProyecto(clave string, titulo string, desc string, repo Repo) error {
	const op string = "app.NuevoProyecto"
	pro := ust.Proyecto{
		ProyectoID:    nuevoProyectoID(clave),
		Titulo:        strings.TrimSpace(titulo),
		Descripcion:   strings.TrimSpace(desc),
		Color:         "lightblue",
		FechaRegistro: gkt.Now().Format(gkt.FormatoFechaHora),
		Posicion:      1,
	}
	if pro.ProyectoID == "" {
		return gko.ErrDatoIndef.Msg("clave de proyecto indefinida").Op(op)
	}
	if pro.Titulo == "" {
		return gko.ErrDatoIndef.Msg("titulo de proyecto indefinido").Op(op)
	}
	err := repo.InsertProyecto(pro)
	if err != nil {
		return gko.Err(err)
	}
	return nil
}
