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

func ModificarProyecto(proyectoID string, nuevo ust.Proyecto, repo Repo) error {
	op := gko.Op("app.ModificarProyecto")
	pro, err := repo.GetProyecto(proyectoID)
	if err != nil {
		return op.Err(err)
	}

	nuevo.ProyectoID = nuevoProyectoID(nuevo.ProyectoID)
	nuevo.Titulo = gkt.SinEspaciosExtra(nuevo.Titulo)
	nuevo.Descripcion = gkt.SinEspaciosExtraConSaltos(nuevo.Descripcion)
	nuevo.Color = gkt.SinEspaciosNinguno(nuevo.Color)

	if nuevo.ProyectoID == "" {
		return op.E(gko.ErrDatoIndef).Msg("clave de proyecto indefinida")
	}
	if nuevo.Titulo == "" {
		return op.E(gko.ErrDatoIndef).Msg("titulo de proyecto indefinido")
	}

	if pro.ProyectoID != nuevo.ProyectoID {
		return op.E(gko.ErrNoSoportado).Msg("no se puede cambiar la clave del proyecto")
	}
	if pro.Titulo != nuevo.Titulo {
		pro.Titulo = nuevo.Titulo
	}
	if pro.Descripcion != nuevo.Descripcion {
		pro.Descripcion = nuevo.Descripcion
	}
	if pro.Color != nuevo.Color {
		pro.Color = nuevo.Color
	}
	if pro.Posicion != nuevo.Posicion {
		pro.Posicion = nuevo.Posicion
	}

	err = repo.UpdateProyecto(pro.ProyectoID, *pro)
	if err != nil {
		return op.Err(err)
	}
	return nil
}

func ParcharProyecto(proyectoID string, param string, newVal string, repo Repo) error {
	op := gko.Op("ParcharProyecto").Ctx("proyectoID", proyectoID)
	if proyectoID == "" {
		return op.Msg("el ID del proyecto debe estar definido")
	}
	Proyecto, err := repo.GetProyecto(proyectoID)
	if err != nil {
		return op.Err(err)
	}
	switch param {
	case "titulo":
		Proyecto.Titulo = gkt.SinEspaciosExtra(newVal)
	case "descripcion":
		Proyecto.Descripcion = gkt.SinEspaciosExtraConSaltos(newVal)
	case "color":
		Proyecto.Color = gkt.SinEspaciosNinguno(newVal)
	case "posicion":
		Proyecto.Posicion, err = gkt.ToInt(newVal)
		if err != nil {
			return op.Err(err)
		}
	default:
		return op.Msgf("Parámetro no soportado: %v", param)
	}
	err = repo.UpdateProyecto(Proyecto.ProyectoID, *Proyecto)
	if err != nil {
		return op.Err(err)
	}
	return nil
}

func EliminarProyecto(ProyectoID string, repo Repo) error {
	const op string = "app.QuitarProyecto"
	pers, err := repo.ListNodosPersonas(ProyectoID)
	if err != nil {
		return gko.Err(err).Op(op)
	}
	if len(pers) != 0 {
		return gko.ErrHayHuerfanos.Msg("Para eliminar este proyecto primero elimine todas sus historias y personajes")
	}
	err = repo.DeleteProyecto(ProyectoID)
	if err != nil {
		return gko.Err(err).Op(op)
	}
	return nil
}
