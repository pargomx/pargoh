package dhistorias

import (
	"monorepo/ust"
	"regexp"
	"strings"

	"github.com/pargomx/gecko/gko"
)

func nuevoProyectoID(clave string) string {
	clave = strings.TrimSpace(regexp.MustCompile(`\s+`).ReplaceAllString(clave, ""))
	clave = strings.ReplaceAll(clave, "-", "_")
	return strings.ToLower(clave)
}

func NuevoProyecto(clave string, titulo string, desc string, repo Repo) error {
	const op string = "app.NuevoProyecto"
	pro := ust.Proyecto{
		ProyectoID:  nuevoProyectoID(clave),
		Titulo:      strings.TrimSpace(titulo),
		Descripcion: strings.TrimSpace(desc),
	}
	if pro.ProyectoID == "" {
		return gko.ErrDatoIndef().Op(op).Msg("clave de proyecto indefinida")
	}
	if pro.Titulo == "" {
		return gko.ErrDatoIndef().Op(op).Msg("titulo de proyecto indefinido")
	}
	err := repo.InsertProyecto(pro)
	if err != nil {
		return gko.Err(err)
	}
	return nil
}

func ModificarProyecto(proyectoID string, clave string, titulo string, desc string, repo Repo) error {
	const op string = "app.ModificarProyecto"
	pro, err := repo.GetProyecto(proyectoID)
	if err != nil {
		return gko.Err(err)
	}
	nuevo := ust.Proyecto{
		ProyectoID:  nuevoProyectoID(clave),
		Titulo:      strings.TrimSpace(titulo),
		Descripcion: strings.TrimSpace(desc),
	}
	if nuevo.ProyectoID == "" {
		return gko.ErrDatoIndef().Op(op).Msg("clave de proyecto indefinida")
	}
	if nuevo.Titulo == "" {
		return gko.ErrDatoIndef().Op(op).Msg("titulo de proyecto indefinido")
	}

	if pro.Titulo != nuevo.Titulo {
		pro.Titulo = nuevo.Titulo
	}
	if pro.Descripcion != nuevo.Descripcion {
		pro.Descripcion = nuevo.Descripcion
	}
	if pro.ProyectoID != nuevo.ProyectoID {
		return gko.ErrNoSoportado().Msg("no se puede cambiar la clave del proyecto")
	}

	err = repo.UpdateProyecto(*pro)
	if err != nil {
		return gko.Err(err)
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
		Proyecto.Titulo = newVal
	case "descripcion":
		Proyecto.Descripcion = newVal
	default:
		return op.Msgf("Parámetro no soportado: %v", param)
	}
	err = repo.UpdateProyecto(*Proyecto)
	if err != nil {
		return op.Err(err)
	}
	return nil
}

func QuitarProyecto(ProyectoID string, repo Repo) error {
	const op string = "app.QuitarProyecto"
	pers, err := repo.ListNodosPersonas(ProyectoID)
	if err != nil {
		return gko.Err(err).Op(op)
	}
	if len(pers) != 0 {
		return gko.ErrHayHuerfanos().Msg("Para eliminar este proyecto primero elimine todas sus historias y personajes")
	}
	err = repo.DeleteProyecto(ProyectoID)
	if err != nil {
		return gko.Err(err).Op(op)
	}
	return nil
}
